package processors

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"github.com/nguyenthenguyen/docx"
)

// DocxProcessor processador de arquivos DOCX
type DocxProcessor struct {
	ocrClient *gosseract.Client
}

// NewDocxProcessor cria novo processador de DOCX
func NewDocxProcessor(ocrClient *gosseract.Client) *DocxProcessor {
	return &DocxProcessor{
		ocrClient: ocrClient,
	}
}

// Process processa arquivo DOCX
func (p *DocxProcessor) Process(file io.Reader, filename string) (string, error) {
	log.Printf("📄 Processando DOCX: %s", filename)

	// Criar arquivo temporário
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

	// Tentar extrair texto diretamente do DOCX
	text, err := p.extractTextFromDocx(tempFile.Name())
	if err != nil {
		log.Printf("⚠️ Extração direta falhou: %v", err)
		// Se falhar, tentar OCR
		log.Printf("🔄 Tentando OCR...")
		text, err = p.extractTextWithOCR(tempFile.Name())
		if err != nil {
			return "", fmt.Errorf("erro ao processar DOCX: %v", err)
		}
	}

	if len(strings.TrimSpace(text)) < 10 {
		return "", fmt.Errorf("DOCX não contém texto suficiente para análise")
	}

	log.Printf("✅ DOCX processado com sucesso: %d caracteres extraídos", len(text))
	return strings.TrimSpace(text), nil
}

// extractTextFromDocx extrai texto diretamente do DOCX
func (p *DocxProcessor) extractTextFromDocx(filePath string) (string, error) {
	// Abrir arquivo DOCX
	doc, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo DOCX: %v", err)
	}
	defer doc.Close()

	// Extrair texto
	text := doc.Editable().GetText()
	if len(strings.TrimSpace(text)) == 0 {
		return "", fmt.Errorf("arquivo DOCX não contém texto")
	}

	log.Printf("📖 DOCX lido com sucesso: %d caracteres", len(text))
	return text, nil
}

// extractTextWithOCR extrai texto usando OCR
func (p *DocxProcessor) extractTextWithOCR(filePath string) (string, error) {
	log.Printf("🔍 Iniciando OCR para DOCX: %s", filePath)
	
	p.ocrClient.SetImage(filePath)
	text, err := p.ocrClient.Text()
	if err != nil {
		return "", fmt.Errorf("erro no OCR: %v", err)
	}
	
	log.Printf("✅ OCR concluído: %d caracteres extraídos", len(text))
	return text, nil
}
