package chezmoi

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetUmask returns the umask.
func GetUmask() os.FileMode {
	return os.FileMode(0)
}

// TrimDirPrefix returns path with the directory prefix dir stripped. path must
// be an absolute path with forward slashes.
func TrimDirPrefix(path, dir string) (string, error) {
	prefix := strings.ToLower(dir + "/")
	if !strings.HasPrefix(strings.ToLower(path), prefix) {
		return "", fmt.Errorf("%q does not have dir prefix %q", path, dir)
	}
	return path[len(prefix):], nil
}

// SetUmask sets the umask.
func SetUmask(umask os.FileMode) {}

// isExecutable returns false on Windows.
func isExecutable(info os.FileInfo) bool {
	return false
}

// isPrivate returns false on Windows.
func isPrivate(info os.FileInfo) bool {
	return false
}

func isSlash(c uint8) bool {
	return c == '\\' || c == '/'
}

func normalizePath(p, homeDir string) (string, error) {
	switch {
	case p == "~":
		return homeDir, nil
	case len(p) >= 2 && p[0] == '~' && isSlash(p[1]):
		return filepath.ToSlash(filepath.Join(homeDir, p[2:])), nil
	default:
		var err error
		p, err = filepath.Abs(p)
		if err != nil {
			return "", err
		}
		if n := volumeNameLen(p); n > 0 {
			p = strings.ToUpper(p[:n]) + p[n:]
		}
		return filepath.ToSlash(p), nil
	}
}

// umaskPermEqual returns true on Windows.
func umaskPermEqual(perm1 os.FileMode, perm2 os.FileMode, umask os.FileMode) bool {
	return true
}

// volumeNameLen returns length of the leading volume name on Windows. It
// returns 0 elsewhere.
func volumeNameLen(path string) int {
	if len(path) < 2 {
		return 0
	}
	// with drive letter
	c := path[0]
	if path[1] == ':' && ('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z') {
		return 2
	}
	// is it UNC? https://msdn.microsoft.com/en-us/library/windows/desktop/aa365247(v=vs.85).aspx
	if l := len(path); l >= 5 && isSlash(path[0]) && isSlash(path[1]) &&
		!isSlash(path[2]) && path[2] != '.' {
		// first, leading `\\` and next shouldn't be `\`. its server name.
		for n := 3; n < l-1; n++ {
			// second, next '\' shouldn't be repeated.
			if isSlash(path[n]) {
				n++
				// third, following something characters. its share name.
				if !isSlash(path[n]) {
					if path[n] == '.' {
						break
					}
					for ; n < l; n++ {
						if isSlash(path[n]) {
							break
						}
					}
					return n
				}
				break
			}
		}
	}
	return 0
}
