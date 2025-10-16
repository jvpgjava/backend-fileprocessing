package processors

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// PDFProcessor processador de arquivos PDF
type PDFProcessor struct {
	ocrClient *gosseract.Client
}

// NewPDFProcessor cria novo processador de PDF
func NewPDFProcessor(ocrClient *gosseract.Client) *PDFProcessor {
	return &PDFProcessor{
		ocrClient: ocrClient,
	}
}

// Process processa arquivo PDF
func (p *PDFProcessor) Process(file io.Reader, filename string) (string, error) {
	log.Printf("📄 Processando PDF: %s", filename)

	// Criar arquivo temporário
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

	// Tentar extrair texto diretamente do PDF
	text, err := p.extractTextFromPDF(tempFile.Name())
	if err != nil {
		log.Printf("⚠️ Extração direta falhou: %v", err)
		// Se falhar, tentar OCR
		log.Printf("🔄 Tentando OCR...")
		text, err = p.extractTextWithOCR(tempFile.Name())
		if err != nil {
			return "", fmt.Errorf("erro ao processar PDF: %v", err)
		}
	}

	if len(strings.TrimSpace(text)) < 10 {
		return "", fmt.Errorf("PDF não contém texto suficiente para análise")
	}

	log.Printf("✅ PDF processado com sucesso: %d caracteres extraídos", len(text))
	return strings.TrimSpace(text), nil
}

// extractTextFromPDF extrai texto diretamente do PDF usando unipdf
func (p *PDFProcessor) extractTextFromPDF(filePath string) (string, error) {
	// Abrir arquivo PDF
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Criar leitor PDF
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return "", err
	}

	var fullText strings.Builder
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	log.Printf("📖 PDF tem %d páginas", numPages)

	// Extrair texto de cada página
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			log.Printf("⚠️ Erro ao ler página %d: %v", i, err)
			continue
		}

		ex, err := extractor.New(page)
		if err != nil {
			log.Printf("⚠️ Erro ao criar extrator para página %d: %v", i, err)
			continue
		}

		text, err := ex.ExtractText()
		if err != nil {
			log.Printf("⚠️ Erro ao extrair texto da página %d: %v", i, err)
			continue
		}

		if len(strings.TrimSpace(text)) > 0 {
			fullText.WriteString(text)
			fullText.WriteString("\n")
		}
	}

	return fullText.String(), nil
}

// extractTextWithOCR extrai texto usando OCR
func (p *PDFProcessor) extractTextWithOCR(filePath string) (string, error) {
	log.Printf("🔍 Iniciando OCR para PDF: %s", filePath)
	
	p.ocrClient.SetImage(filePath)
	text, err := p.ocrClient.Text()
	if err != nil {
		return "", fmt.Errorf("erro no OCR: %v", err)
	}
	
	log.Printf("✅ OCR concluído: %d caracteres extraídos", len(text))
	return text, nil
}
