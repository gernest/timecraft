package sandbox

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

type dirFS string

func (root dirFS) Open(name string, flags int, mode fs.FileMode) (File, error) {
	dirfd, err := openat(unix.AT_FDCWD, string(root), O_DIRECTORY, 0)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: string(root), Err: err}
	}
	if name = cleanPath(name); name == "/" || name == "." { // root?
		return dirFile{os.NewFile(uintptr(dirfd), "/")}, nil
	}
	defer closeTraceError(dirfd)
	relPath := "/" + trimLeadingSlash(name)
	fd, err := openat(dirfd, name, flags, uint32(mode.Perm()))
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: relPath, Err: err}
	}
	return dirFile{os.NewFile(uintptr(fd), relPath)}, nil
}

type dirFile struct{ *os.File }

func (f dirFile) Open(name string, flags int, mode fs.FileMode) (File, error) {
	name = cleanPath(name)
	relPath := f.join(name)
	fd, err := openat(f.fd(), name, flags, uint32(mode.Perm()))
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: relPath, Err: err}
	}
	return dirFile{os.NewFile(uintptr(fd), relPath)}, nil
}

func (f dirFile) Datasync() error {
	if err := fdatasync(f.fd()); err != nil {
		return &fs.PathError{Op: "datasync", Path: f.Name(), Err: err}
	}
	return nil
}

func (f dirFile) Flags() (int, error) {
	flags, err := unix.FcntlInt(f.Fd(), unix.F_GETFL, 0)
	if err != nil {
		return 0, &fs.PathError{Op: "fcntl", Path: f.Name(), Err: err}
	}
	return flags, nil
}

func (f dirFile) SetFlags(flags int) error {
	_, err := unix.FcntlInt(f.Fd(), unix.F_SETFL, flags)
	if err != nil {
		return &fs.PathError{Op: "fcntl", Path: f.Name(), Err: err}
	}
	return nil
}

func (f dirFile) Stat(name string, flags int) (FileInfo, error) {
	var stat unix.Stat_t
	var err error

	fd := f.fd()
	if name == "" {
		err = fstat(fd, &stat)
	} else {
		err = fstatat(fd, name, &stat, flags)
	}
	if err != nil {
		return FileInfo{}, &fs.PathError{Op: "stat", Path: f.join(name), Err: err}
	}

	info := FileInfo{
		Dev:   uint64(stat.Dev),
		Ino:   uint64(stat.Ino),
		Nlink: uint64(stat.Nlink),
		Mode:  fs.FileMode(stat.Mode & 0777), // perm
		Size:  int64(stat.Size),
		Atime: time.Unix(stat.Atim.Unix()),
		Mtime: time.Unix(stat.Mtim.Unix()),
		Ctime: time.Unix(stat.Ctim.Unix()),
	}

	switch stat.Mode & unix.S_IFMT {
	case unix.S_IFREG:
	case unix.S_IFBLK:
		info.Mode |= fs.ModeDevice
	case unix.S_IFCHR:
		info.Mode |= fs.ModeDevice | fs.ModeCharDevice
	case unix.S_IFDIR:
		info.Mode |= fs.ModeDir
	case unix.S_IFIFO:
		info.Mode |= fs.ModeNamedPipe
	case unix.S_IFLNK:
		info.Mode |= fs.ModeSymlink
	case unix.S_IFSOCK:
		info.Mode |= fs.ModeSocket
	default:
		info.Mode |= fs.ModeIrregular
	}
	return info, nil
}

func (f dirFile) Readlink(name string, buf []byte) (n int, err error) {
	fd := f.fd()
	if name == "" {
		n, err = freadlink(fd, buf)
	} else {
		n, err = readlinkat(fd, name, buf)
	}
	if err != nil {
		err = &fs.PathError{Op: "readlink", Path: f.join(name), Err: err}
	}
	return n, err
}

func (f dirFile) Chtimes(name string, atime, mtime time.Time, flags int) error {
	var ts [2]unix.Timespec
	var err error
	ts[0], _ = unix.TimeToTimespec(atime)
	ts[1], _ = unix.TimeToTimespec(mtime)
	if name == "" {
		err = futimens(f.fd(), &ts)
	} else {
		err = utimensat(f.fd(), name, &ts, flags)
	}
	if err != nil {
		return &fs.PathError{Op: "chtimes", Path: f.join(name), Err: err}
	}
	return err
}

func (f dirFile) Mkdir(name string, mode fs.FileMode) error {
	if err := mkdirat(f.fd(), name, uint32(mode.Perm())); err != nil {
		return &fs.PathError{Op: "mkdir", Path: f.join(name), Err: err}
	}
	return nil
}

func (f dirFile) Rmdir(name string) error {
	if err := unlinkat(f.fd(), name, unix.AT_REMOVEDIR); err != nil {
		return &fs.PathError{Op: "rmdir", Path: f.join(name), Err: err}
	}
	return nil
}

func (f dirFile) Rename(oldName string, newDir File, newName string) error {
	fd1 := f.fd()
	fd2 := int(newDir.Fd())
	if err := renameat(fd1, oldName, fd2, newName); err != nil {
		path1 := f.join(oldName)
		path2 := joinPath(newDir.Name(), newName)
		return &os.LinkError{Op: "rename", Old: path1, New: path2, Err: err}
	}
	return nil
}

func (f dirFile) Link(oldName string, newDir File, newName string, flags int) error {
	linkFlags := 0
	if (flags & AT_SYMLINK_NOFOLLOW) == 0 {
		linkFlags |= unix.AT_SYMLINK_FOLLOW
	}
	fd1 := f.fd()
	fd2 := int(newDir.Fd())
	if err := linkat(fd1, oldName, fd2, newName, linkFlags); err != nil {
		path1 := f.join(oldName)
		path2 := joinPath(newDir.Name(), newName)
		return &os.LinkError{Op: "link", Old: path1, New: path2, Err: err}
	}
	return nil
}

func (f dirFile) Symlink(oldName string, newName string) error {
	if err := symlinkat(oldName, f.fd(), newName); err != nil {
		return &fs.PathError{Op: "symlink", Path: f.join(newName), Err: err}
	}
	return nil
}

func (f dirFile) Unlink(name string) error {
	if err := unlinkat(f.fd(), name, 0); err != nil {
		return &fs.PathError{Op: "unlink", Path: f.join(name), Err: err}
	}
	return nil
}

func (f dirFile) fd() int {
	return int(f.Fd())
}

func (f dirFile) join(name string) string {
	return joinPath(f.Name(), name)
}

type dirFileInfo struct {
	stat unix.Stat_t
	name string
}

func (info *dirFileInfo) Name() string {
	return path.Base(info.name)
}

func (info *dirFileInfo) Size() int64 {
	return info.stat.Size
}

func (info *dirFileInfo) Mode() (mode fs.FileMode) {
	mode = fs.FileMode(info.stat.Mode & 0777) // perm
	switch info.stat.Mode & unix.S_IFMT {
	case unix.S_IFREG:
	case unix.S_IFBLK:
		mode |= fs.ModeDevice
	case unix.S_IFCHR:
		mode |= fs.ModeDevice | fs.ModeCharDevice
	case unix.S_IFDIR:
		mode |= fs.ModeDir
	case unix.S_IFIFO:
		mode |= fs.ModeNamedPipe
	case unix.S_IFLNK:
		mode |= fs.ModeSymlink
	case unix.S_IFSOCK:
		mode |= fs.ModeSocket
	default:
		mode |= fs.ModeIrregular
	}
	return mode
}

func (info *dirFileInfo) ModTime() time.Time {
	return time.Unix(info.stat.Mtim.Unix())
}

func (info *dirFileInfo) IsDir() bool {
	return (info.stat.Mode & unix.S_IFDIR) != 0
}

func (info *dirFileInfo) Sys() any {
	return &info.stat
}

func (info *dirFileInfo) String() string {
	group, name := "none", "nobody"

	g, err := user.LookupGroupId(strconv.Itoa(int(info.stat.Gid)))
	if err == nil {
		group = g.Name
	}

	u, err := user.LookupId(strconv.Itoa(int(info.stat.Uid)))
	if err == nil {
		name = u.Username
	}

	return fmt.Sprintf("%s %d %s %s %d %s %s",
		info.Mode(),
		info.stat.Nlink,
		name,
		group,
		info.Size(),
		info.ModTime().Format(time.Stamp),
		info.Name(),
	)
}
