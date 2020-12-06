// +build !linux

package chezmoi

import "github.com/twpayne/go-vfs"

func KernelInfo(fs vfs.FS) (map[string]string, error) {
	return nil, nil
}

func OSRelease(fs vfs.FS) (map[string]string, error) {
	return nil, nil
}
