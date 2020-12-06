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
	wd, err := os.Getwd()
	require.NoError(t, err)
	for _, tc := range []struct {
		name       string
		s          string
		homeDirStr string
		expected   string
	}{
		{
			name:     "empty",
			expected: wd,
		},
		{
			name:     "file",
			s:        "file",
			expected: path.Join(wd, "file"),
		},
		{
			name:       "tilde",
			s:          "~",
			homeDirStr: "/home/user",
			expected:   "/home/user",
		},
		{
			name:       "tilde_home",
			s:          "~",
			homeDirStr: "/home/user",
			expected:   "/home/user",
		},
		{
			name:       "tilde_home_file",
			s:          "~/file",
			homeDirStr: "/home/user",
			expected:   "/home/user/file",
		},
		{
			name:       "tilde_home_file_windows",
			s:          `~\file`,
			homeDirStr: `C:\home\user`,
			expected:   `C:\home\user\file`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			chezmoitest.SkipUnlessGOOS(t, tc.name)

			actual, err := NewOSPath(tc.s).TildeAbsSlash(tc.homeDirStr)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

//nolint:paralleltest,tparallel
func TestOSPathFormat(t *testing.T) {
	t.Parallel()

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
