package chezmoi

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/twpayne/chezmoi/next/internal/chezmoitest"
)

func TestOSPathTildeAbsSlash(t *testing.T) {
	var homeDirNormalized string
	switch runtime.GOOS {
	case "windows":
		homeDirNormalized = "C:/home/user"
	default:
		homeDirNormalized = "/home/user"
	}

	wd, err := os.Getwd()
	require.NoError(t, err)
	wdNormalized, err := normalizePath(wd, homeDirNormalized)
	require.NoError(t, err)

	for _, tc := range []struct {
		name     string
		s        string
		expected string
	}{
		{
			name:     "empty",
			expected: wdNormalized,
		},
		{
			name:     "file",
			s:        "file",
			expected: path.Join(wdNormalized, "file"),
		},
		{
			name:     "tilde",
			s:        "~",
			expected: homeDirNormalized,
		},
		{
			name:     "tilde_home_file",
			s:        "~/file",
			expected: homeDirNormalized + "/file",
		},
		{
			name:     "tilde_home_file_windows",
			s:        `~\file`,
			expected: homeDirNormalized + "/file",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			chezmoitest.SkipUnlessGOOS(t, tc.name)

			actual, err := NewOSPath(tc.s).TildeAbsSlash(homeDirNormalized)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestOSPathFormat(t *testing.T) {
	type s struct {
		Dir *OSPath
	}

	for name, format := range Formats {
		t.Run(name, func(t *testing.T) {
			var dirStr string
			switch runtime.GOOS {
			case "windows":
				dirStr = `C:\home\user`
			default:
				dirStr = "/home/user"
			}
			expectedS := &s{
				Dir: NewOSPath(dirStr),
			}
			data, err := format.Marshal(expectedS)
			assert.NoError(t, err)
			actualS := &s{}
			assert.NoError(t, format.Decode(data, actualS))
			assert.Equal(t, expectedS, actualS)
		})
	}
}
