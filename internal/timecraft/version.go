package timecraft

import "runtime/debug"

// Version returns the timecraft version.
func Version() string {
	version := "devel"
	if info, ok := debug.ReadBuildInfo(); ok {
		switch info.Main.Version {
		case "":
		case "(devel)":
		default:
			version = info.Main.Version
		}
	}
	return version
}
