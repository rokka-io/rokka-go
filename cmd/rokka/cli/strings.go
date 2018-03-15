package cli

// ListContains check if s is in list l.
func ListContains(l []string, s string) bool {
	for _, v := range l {
		if s == v {
			return true
		}
	}
	return false
}
