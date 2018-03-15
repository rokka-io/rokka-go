// +build !windows

package massupload

// fixpath returns an absolute path on the current OS, so we can open long
// file names. See Restic file_unix.go
func fixpath(name string) string {
	return name
}
