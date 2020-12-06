package chezmoi

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// An OSPath is a native OS path.
type OSPath struct {
	s string
}

// NewOSPath returns a new OSPath.
func NewOSPath(s string) *OSPath {
	return &OSPath{
		s: filepath.FromSlash(s),
	}
}

// TildeAbsSlash returns p converted to an absolute path with a leading tilde
// expanded to homeDirStr and backslashes replaced by forward slashes.
func (p *OSPath) TildeAbsSlash(homeDirStr string) (string, error) {
	s := p.s
	switch {
	case s == "~":
		absHomeDirStr, err := filepath.Abs(homeDirStr)
		if err != nil {
			return "", err
		}
		return filepath.ToSlash(absHomeDirStr), nil
	case strings.HasPrefix(s, "~/"):
		fallthrough
	case runtime.GOOS == "windows" && strings.HasPrefix(s, "~"+string(rune(os.PathSeparator))):
		s = homeDirStr + string(rune(os.PathSeparator)) + s[2:]
	}
	abs, err := filepath.Abs(s)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(abs), nil
}

// Dir returns p's directory.
func (p *OSPath) Dir() *OSPath {
	return &OSPath{
		s: filepath.Dir(p.s),
	}
}

// Empty returns if p is empty.
func (p *OSPath) Empty() bool {
	return p.s != ""
}

// Join joins elems on to p.
func (p *OSPath) Join(elems ...string) *OSPath {
	return &OSPath{
		s: filepath.Join(append([]string{p.s}, elems...)...),
	}
}

// MarshalText implements encoding.TextMarshaler.MarshalText.
func (p *OSPath) MarshalText() ([]byte, error) {
	return []byte(p.s), nil
}

func (p *OSPath) String() string {
	return p.s
}

// UnmarshalText implements encoding.TextUnmarshaler.UnmarshalText.
func (p *OSPath) UnmarshalText(data []byte) error {
	p.s = string(data)
	return nil
}
