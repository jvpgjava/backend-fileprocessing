package processors

import "io"

// GeminiExtractor interface para extrair texto de arquivos usando Gemini
// Isso evita ciclo de importação
type GeminiExtractor interface {
	ExtractTextFromPDF(fileReader io.Reader, filename string) (string, error)
	ExtractTextFromFile(fileReader io.Reader, filename string) (string, error)
	IsAvailable() bool
}

