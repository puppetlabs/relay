package version

import "runtime/debug"

var Version string

func GetVersion() string {
	if Version == "" {
		i, ok := debug.ReadBuildInfo()
		if !ok {
			return ""
		}
		Version = i.Main.Version
	}
	return Version
}
