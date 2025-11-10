package processors

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// DocxProcessor processador de arquivos DOCX usando Google Gemini
type DocxProcessor struct {
	geminiExtractor GeminiExtractor
}

// NewDocxProcessor cria novo processador de DOCX
func NewDocxProcessor(geminiExtractor GeminiExtractor) *DocxProcessor {
	return &DocxProcessor{
		geminiExtractor: geminiExtractor,
	}
}

// Process processa arquivo DOCX usando Google Gemini
func (p *DocxProcessor) Process(file io.Reader, filename string) (string, error) {
	log.Printf("📄 Processando DOCX: %s", filename)

	// Verificar se Gemini está disponível
	if p.geminiExtractor == nil || !p.geminiExtractor.IsAvailable() {
		return "", fmt.Errorf("Gemini não está disponível - GEMINI_API_KEY não configurada")
	}

	// Criar arquivo temporário para poder reler
	tempFile, err := os.CreateTemp("", "temp_*.docx")
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
	log.Printf("🤖 Processando DOCX com Google Gemini...")

	// Ler arquivo novamente para passar para Gemini
	fileReader, err := os.Open(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("erro ao reabrir arquivo para Gemini: %v", err)
	}
	defer fileReader.Close()

	text, err := p.geminiExtractor.ExtractTextFromFile(fileReader, filename)
	if err != nil {
		return "", fmt.Errorf("erro ao processar DOCX com Gemini: %v", err)
	}

	if len(strings.TrimSpace(text)) < 10 {
		return "", fmt.Errorf("Gemini extraiu pouco texto (menos de 10 caracteres)")
	}

	log.Printf("✅ Gemini extraiu texto do DOCX: %d caracteres", len(text))
	return strings.TrimSpace(text), nil
}
