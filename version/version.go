package version

import (
	"fmt"
	"runtime/debug"
)

type BuildInfo struct {
	Date        string
	FullGitSHA  string
	GoVersion   string
	ShortGitSHA string
	Version     string
}

func (bi BuildInfo) String() string {
	return fmt.Sprintf(`Version: %s, %s
Build Date: %s
Go Version: %s`, bi.Version, bi.FullGitSHA, bi.Date, bi.GoVersion)
}

// GetBuildInfo reads the metadata the go toolchain embeds in every binary:
// the module version for go install pkg@version, and the VCS revision and
// commit time for builds inside a git checkout.
func GetBuildInfo() BuildInfo {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return BuildInfo{}
	}

	bi := BuildInfo{
		GoVersion: info.GoVersion,
		Version:   info.Main.Version,
	}
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			bi.FullGitSHA = s.Value
		case "vcs.time":
			bi.Date = s.Value
		}
	}
	if len(bi.FullGitSHA) >= 7 {
		bi.ShortGitSHA = bi.FullGitSHA[:7]
	}
	return bi
}
