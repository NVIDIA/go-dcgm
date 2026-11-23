package dcgm

import (
	"strconv"
	"strings"
)

// runtimeVersion reads the major, minor, and patch numbers from DCGM build info.
// It reports an unknown version if the version key is missing, repeated, or invalid.
func runtimeVersion(info VersionInfo) ([3]uint64, bool) {
	var version [3]uint64
	text := ""
	found := false
	for _, pair := range strings.Split(info.RawBuildInfoString, ";") {
		key, value, ok := strings.Cut(pair, ":")
		if !ok || strings.TrimSpace(key) != "version" {
			continue
		}
		if found {
			return version, false
		}
		found = true
		text = strings.TrimSpace(value)
	}
	parts := strings.Split(text, ".")
	if len(parts) != len(version) {
		return version, false
	}
	for i, part := range parts {
		value, err := strconv.ParseUint(part, 10, 32)
		if err != nil || (len(part) > 1 && part[0] == '0') {
			return [3]uint64{}, false
		}
		version[i] = value
	}
	return version, true
}
