package processors

import (
    "archive/zip"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "html"

    "github.com/otiai10/gosseract/v2"
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
    // DOCX é um ZIP; vamos ler word/document.xml e extrair o texto simples
    zr, err := zip.OpenReader(filePath)
    if err != nil {
        return "", fmt.Errorf("erro ao abrir DOCX (zip): %v", err)
    }
    defer zr.Close()

    var xmlData string
    for _, f := range zr.File {
        // Documento principal
        if f.Name == "word/document.xml" || f.Name == filepath.ToSlash("word/document.xml") {
            rc, err := f.Open()
            if err != nil {
                return "", fmt.Errorf("erro ao abrir document.xml: %v", err)
            }
            b, err := io.ReadAll(rc)
            rc.Close()
            if err != nil {
                return "", fmt.Errorf("erro ao ler document.xml: %v", err)
            }
            xmlData = string(b)
            break
        }
    }

    if xmlData == "" {
        return "", fmt.Errorf("document.xml não encontrado no DOCX")
    }

    // Remover tags XML simples e normalizar espaços
    // 1) substituir quebras de parágrafo por nova linha
    xmlData = strings.ReplaceAll(xmlData, "</w:p>", "\n")
    // 2) remover todas as tags
    re := regexp.MustCompile(`<[^>]+>`) // tags XML
    plain := re.ReplaceAllString(xmlData, "")
    // 3) unescape entidades
    plain = html.UnescapeString(plain)
    // 4) normalizar espaços
    plain = strings.TrimSpace(plain)

    if len(strings.TrimSpace(plain)) == 0 {
        return "", fmt.Errorf("arquivo DOCX não contém texto")
    }

    log.Printf("📖 DOCX lido com sucesso: %d caracteres", len(plain))
    return plain, nil
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
