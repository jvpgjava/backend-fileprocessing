package processors

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// ImageProcessor processador de arquivos de imagem usando Google Gemini
type ImageProcessor struct {
	geminiExtractor GeminiExtractor
}

// NewImageProcessor cria novo processador de imagem
func NewImageProcessor(geminiExtractor GeminiExtractor) *ImageProcessor {
	return &ImageProcessor{
		geminiExtractor: geminiExtractor,
	}
}

// Process processa arquivo de imagem usando Google Gemini
func (p *ImageProcessor) Process(file io.Reader, filename string) (string, error) {
	log.Printf("🖼️ Processando imagem: %s", filename)

	// Verificar se Gemini está disponível
	if p.geminiExtractor == nil || !p.geminiExtractor.IsAvailable() {
		return "", fmt.Errorf("Gemini não está disponível - GEMINI_API_KEY não configurada")
	}

	// Criar arquivo temporário para poder reler
	tempFile, err := os.CreateTemp("", "temp_*")
	if err != nil {
		return "", fmt.Errorf("erro ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copiar conteúdo do arquivo
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", fmt.Errorf("erro ao copiar arquivo: %v", err)
	}

	// Processar com Gemini
	log.Printf("🤖 Processando imagem com Google Gemini...")

	// Ler arquivo novamente para passar para Gemini
	fileReader, err := os.Open(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("erro ao reabrir arquivo para Gemini: %v", err)
	}
	defer fileReader.Close()

	text, err := p.geminiExtractor.ExtractTextFromFile(fileReader, filename)
	if err != nil {
		return "", fmt.Errorf("erro ao processar imagem com Gemini: %v", err)
	}

	if len(strings.TrimSpace(text)) < 10 {
		return "", fmt.Errorf("Gemini extraiu pouco texto (menos de 10 caracteres)")
	}

	log.Printf("✅ Gemini extraiu texto da imagem: %d caracteres", len(text))
	return strings.TrimSpace(text), nil
}
