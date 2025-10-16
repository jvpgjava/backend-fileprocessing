package processors

import "io"

// FileProcessor interface para processadores de arquivo
type FileProcessor interface {
	Process(file io.Reader, filename string) (string, error)
}
