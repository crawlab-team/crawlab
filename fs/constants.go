package fs

import "os"

const (
	FilerResponseNotFoundErrorMessage = "response status code: 404"
	FilerStatusNotFoundErrorMessage   = "Status:404 Not Found"
)

const (
	DefaultDirMode  = os.FileMode(0766)
	DefaultFileMode = os.FileMode(0666)
)

const (
	MethodUpdateFile = "update-file"
	MethodUploadFile = "upload-file"
	MethodUploadDir  = "upload-dir"
)
