package buildinfo

import "fmt"

func VersionString(full bool) string {
	version := Version
	if version == "" {
		version = "unknown"
	}

	if full {
		buildtime := BuildTime
		if buildtime == "" {
			buildtime = "unknown"
		}
		origin := Origin
		if origin == "" {
			origin = "unknown"
		}
		return fmt.Sprintf("%s_%s_%s", version, buildtime, origin)
	}

	return version
}
