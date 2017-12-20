package rokka

// StrPtr returns a pointer to the passed string
func StrPtr(v string) *string { return &v }

// IntPtr returns a pointer to the passed int
func IntPtr(v int) *int { return &v }

// Float64Ptr returns a pointer to the passed float64
func Float64Ptr(v float64) *float64 { return &v }

// BoolPtr returns a pointer to the passed bool
func BoolPtr(v bool) *bool { return &v }
