package bees

import "io"

type BinaryValue struct {
	MimeType string
	Data     io.ReadCloser
}
