package root

import (
	"runtime/debug"
	"strings"
)

const modulePath = "github.com/arjunuxd/uxdguard"

var (
	linkerVersion     string
	linkerCommit      string
	linkerBuildDate   string
	versionOverridden bool
)

// BuildInfo describes the version metadata embedded in the executable.
type BuildInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

// SetVersion records a version supplied at build time through linker flags.
// Empty values are ignored so module metadata remains the fallback.
func SetVersion(v string) {
	v = normalizeVersion(v)
	if v != "" {
		linkerVersion = v
		versionOverridden = true
	}
}

// SetBuildMetadata records commit and build-date values supplied through linker
// flags. Empty values are ignored so VCS metadata remains the fallback.
func SetBuildMetadata(commit, buildDate string) {
	linkerCommit = strings.TrimSpace(commit)
	linkerBuildDate = strings.TrimSpace(buildDate)
}

// ResolveBuildInfo returns the effective version metadata. Linker-injected
// values take priority, then Go module/VCS build information, then truthful
// "dev"/"unknown" fallbacks.
func ResolveBuildInfo() BuildInfo {
	info := BuildInfo{
		Version:   "dev",
		Commit:    "unknown",
		BuildDate: "unknown",
	}

	if versionOverridden {
		info.Version = linkerVersion
	}

	if build, ok := debug.ReadBuildInfo(); ok {
		moduleVersion := strings.TrimSpace(build.Main.Version)
		if !versionOverridden &&
			build.Main.Path == modulePath &&
			moduleVersion != "" &&
			moduleVersion != "(devel)" {
			info.Version = normalizeVersion(moduleVersion)
		}

		for _, setting := range build.Settings {
			switch setting.Key {
			case "vcs.revision":
				if setting.Value != "" {
					info.Commit = setting.Value
				}
			case "vcs.time":
				if setting.Value != "" {
					info.BuildDate = setting.Value
				}
			}
		}
	}

	if linkerCommit != "" {
		info.Commit = linkerCommit
	}
	if linkerBuildDate != "" {
		info.BuildDate = linkerBuildDate
	}

	return info
}

// CurrentVersion returns the resolved version string without a leading "v".
func CurrentVersion() string {
	return ResolveBuildInfo().Version
}

func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	return strings.TrimPrefix(v, "v")
}
