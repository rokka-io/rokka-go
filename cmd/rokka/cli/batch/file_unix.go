// +build !windows

package batch

// Fixpath returns an absolute path on the current OS which is mostly relevant for windows.
func Fixpath(name string) string {
	return name
}
