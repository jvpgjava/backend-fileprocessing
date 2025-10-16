package processors

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

// ImageProcessor processador de arquivos de imagem
type ImageProcessor struct {
	ocrClient *gosseract.Client
}

// NewImageProcessor cria novo processador de imagem
func NewImageProcessor(ocrClient *gosseract.Client) *ImageProcessor {
	return &ImageProcessor{
		ocrClient: ocrClient,
	}
}

// Process processa arquivo de imagem
func (p *ImageProcessor) Process(file io.Reader, filename string) (string, error) {
	log.Printf("🖼️ Processando imagem: %s", filename)

	// Obter extensão do arquivo
	ext := strings.ToLower(filepath.Ext(filename))

	// Criar arquivo temporário
	tempFile, err := os.CreateTemp("", "temp_*"+ext)
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

	// Extrair texto com OCR
	text, err := p.extractTextWithOCR(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("erro ao processar imagem: %v", err)
	}

	if len(strings.TrimSpace(text)) == 0 {
		imageType := strings.ToUpper(ext[1:])
		text = fmt.Sprintf("[IMAGEM %s] - Nenhum texto detectado na imagem", imageType)
	}

	log.Printf("✅ Imagem processada com sucesso: %d caracteres extraídos", len(text))
	return strings.TrimSpace(text), nil
}

// extractTextWithOCR extrai texto usando OCR
func (p *ImageProcessor) extractTextWithOCR(filePath string) (string, error) {
	log.Printf("🔍 Iniciando OCR para imagem: %s", filePath)
	
	p.ocrClient.SetImage(filePath)
	text, err := p.ocrClient.Text()
	if err != nil {
		return "", fmt.Errorf("erro no OCR: %v", err)
	}
	
	log.Printf("✅ OCR concluído: %d caracteres extraídos", len(text))
	return text, nil
}
