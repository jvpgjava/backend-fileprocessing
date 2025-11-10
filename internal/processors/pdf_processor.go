package processors

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// PDFProcessor processador de arquivos PDF usando APENAS Google Gemini
type PDFProcessor struct {
	geminiExtractor GeminiExtractor
}

// NewPDFProcessor cria novo processador de PDF
func NewPDFProcessor(geminiExtractor GeminiExtractor) *PDFProcessor {
	return &PDFProcessor{
		geminiExtractor: geminiExtractor,
	}
}

// Process processa arquivo PDF usando APENAS Google Gemini
func (p *PDFProcessor) Process(file io.Reader, filename string) (string, error) {
	log.Printf("📄 Processando PDF: %s", filename)

	// Verificar se Gemini está disponível
	if p.geminiExtractor == nil || !p.geminiExtractor.IsAvailable() {
		return "", fmt.Errorf("Gemini não está disponível - GEMINI_API_KEY não configurada. Configure a variável de ambiente GEMINI_API_KEY")
	}

	// Criar arquivo temporário para poder reler
	tempFile, err := os.CreateTemp("", "temp_*.pdf")
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

	// Processar com Gemini (APENAS!)
	log.Printf("🤖 Processando PDF com Google Gemini (gratuito)...")

	// Ler arquivo novamente para passar para Gemini
	fileReader, err := os.Open(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("erro ao reabrir arquivo para Gemini: %v", err)
	}
	defer fileReader.Close()

	geminiText, err := p.geminiExtractor.ExtractTextFromFile(fileReader, filename)
	if err != nil {
		return "", fmt.Errorf("erro ao processar PDF com Gemini: %v", err)
	}

	if len(strings.TrimSpace(geminiText)) < 10 {
		return "", fmt.Errorf("Gemini extraiu pouco texto (menos de 10 caracteres)")
	}

	log.Printf("✅ Gemini extraiu texto com sucesso: %d caracteres", len(geminiText))
	return strings.TrimSpace(geminiText), nil
}

