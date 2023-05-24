package timemachine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/stealthrocket/timecraft/internal/object"
)

type ModuleInfo struct {
	ID        Hash
	Size      int64
	CreatedAt time.Time
}

type LogSegmentInfo struct {
	Number    int
	Size      int64
	CreatedAt time.Time
}

type Store struct {
	objects object.Store
}

func NewStore(objects object.Store) *Store {
	return &Store{objects: objects}
}

func (store *Store) CreateModule(ctx context.Context, id Hash, wasm []byte) error {
	return store.objects.CreateObject(ctx, store.moduleKey(id), bytes.NewReader(wasm))
}

func (store *Store) DeleteModule(ctx context.Context, id Hash) error {
	return store.objects.DeleteObject(ctx, store.moduleKey(id))
}

func (store *Store) ReadModule(ctx context.Context, id Hash) ([]byte, error) {
	key := store.moduleKey(id)
	s, err := store.objects.StatObject(ctx, key)
	if err != nil {
		return nil, err
	}
	r, err := store.objects.ReadObject(ctx, key)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	b := make([]byte, s.Size)
	n, err := io.ReadFull(r, b)
	return b[:n], err
}

func (store *Store) ListModules(ctx context.Context) object.Iter[ModuleInfo] {
	return object.ConvertIter(store.objects.ListObjects(ctx, "modules"), func(info object.Info) (ModuleInfo, error) {
		hash := path.Base(info.Name)
		algorithm, digest, ok := strings.Cut(hash, ":")
		if !ok {
			return ModuleInfo{}, fmt.Errorf("invalid module entry: %q", info.Name)
		}
		module := ModuleInfo{
			ID:        Hash{Algorithm: algorithm, Digest: digest},
			Size:      info.Size,
			CreatedAt: info.CreatedAt,
		}
		return module, nil
	})
}

func (store *Store) CreateLogSegment(ctx context.Context, processID Hash, segmentNumber int) (io.WriteCloser, error) {
	r, w := io.Pipe()
	name := store.logKey(processID, segmentNumber)
	go func() {
		r.CloseWithError(store.objects.CreateObject(ctx, name, r))
	}()
	return w, nil
}

func (store *Store) DeleteLogSegment(ctx context.Context, processID Hash, segmentNumber int) error {
	return store.objects.DeleteObject(ctx, store.logKey(processID, segmentNumber))
}

func (store *Store) ListLogSegments(ctx context.Context, processID Hash) object.Iter[LogSegmentInfo] {
	return object.ConvertIter(store.objects.ListObjects(ctx, "logs/"+processID.String()), func(info object.Info) (LogSegmentInfo, error) {
		number := path.Base(info.Name)
		n, err := strconv.ParseInt(number, 16, 32)
		if err != nil || n < 0 {
			return LogSegmentInfo{}, fmt.Errorf("invalid log segment entry: %q", info.Name)
		}
		segment := LogSegmentInfo{
			Number:    int(n),
			Size:      info.Size,
			CreatedAt: info.CreatedAt,
		}
		return segment, nil
	})
}

func (store *Store) ReadLogSegment(ctx context.Context, processID Hash, segmentNumber int) (io.ReadCloser, error) {
	return store.objects.ReadObject(ctx, store.logKey(processID, segmentNumber))
}

func (store *Store) logKey(processID Hash, segmentNumber int) string {
	return fmt.Sprintf("logs/%s/%08X", processID, segmentNumber)
}

func (store *Store) moduleKey(id Hash) string {
	return fmt.Sprintf("modules/%s", id)
}
