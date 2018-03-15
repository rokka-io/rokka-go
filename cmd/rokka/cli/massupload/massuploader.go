package massupload

type ImageDetails struct {
	Organization string
	Path         string
	UserMetadata map[string]interface{}
}

type UploadResult struct {
	Path      string
	RokkaHash string
	Error     error
}
