package definitions

import "io"

type File struct {
	File     io.Reader
	Filename string
	Size     int64
}
