package processors

import (
	"io"
	"log"
)

// TextProcessor processador de arquivos de texto
type TextProcessor struct{}

// NewTextProcessor cria novo processador de texto
func NewTextProcessor() *TextProcessor {
	return &TextProcessor{}
}

// Process processa arquivo de texto
func (p *TextProcessor) Process(file io.Reader, filename string) (string, error) {
	log.Printf("📝 Processando texto: %s", filename)

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	text := string(content)
	log.Printf("✅ Texto processado com sucesso: %d caracteres", len(text))
	return text, nil
}
