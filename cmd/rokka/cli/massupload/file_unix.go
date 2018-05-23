// +build !windows

package massupload

// fixpath returns an absolute path on the current OS, so we can open long
// file names. See Restic file_unix.go
func Fixpath(name string) string {
	return name
}
