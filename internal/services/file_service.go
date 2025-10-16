package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend-fileprocessing/internal/models"
	"backend-fileprocessing/internal/processors"

	"github.com/otiai10/gosseract/v2"
)

// FileService serviço de processamento de arquivos
type FileService struct {
	ocrClient *gosseract.Client
	processors map[string]processors.FileProcessor
}

// NewFileService cria novo serviço de arquivos
func NewFileService() *FileService {
	// Inicializar cliente OCR
	ocrClient := gosseract.NewClient()
	ocrClient.SetLanguage("por", "eng")

	// Mapear processadores por tipo de arquivo
	processorsMap := map[string]processors.FileProcessor{
		".pdf":  processors.NewPDFProcessor(ocrClient),
		".png":  processors.NewImageProcessor(ocrClient),
		".jpg":  processors.NewImageProcessor(ocrClient),
		".jpeg": processors.NewImageProcessor(ocrClient),
		".gif":  processors.NewImageProcessor(ocrClient),
		".bmp":  processors.NewImageProcessor(ocrClient),
		".webp": processors.NewImageProcessor(ocrClient),
		".tiff": processors.NewImageProcessor(ocrClient),
		".txt":  processors.NewTextProcessor(),
		".docx": processors.NewDocxProcessor(ocrClient),
	}

	return &FileService{
		ocrClient:  ocrClient,
		processors: processorsMap,
	}
}

// ProcessFile processa arquivo e extrai texto
func (fs *FileService) ProcessFile(file io.Reader, filename string, size int64) (*models.Response, error) {
	startTime := time.Now()
	fileType := strings.ToLower(filepath.Ext(filename))
	info := models.NewInfo(filename, fileType, size)

	log.Printf("📁 Processando arquivo: %s (%.2f MB)", filename, float64(size)/1024/1024)

	// Verificar se tipo é suportado
	processor, exists := fs.processors[fileType]
	if !exists {
		return &models.NewErrorResponse(
			"UNSUPPORTED_FILE_TYPE",
			fmt.Sprintf("Tipo de arquivo não suportado: %s", fileType),
			"Tipos suportados: .pdf, .png, .jpg, .jpeg, .gif, .bmp, .webp, .tiff, .txt, .docx",
		), nil
	}

	// Processar arquivo
	text, err := processor.Process(file, filename)
	if err != nil {
		return &models.NewErrorResponse(
			"PROCESSING_ERROR",
			fmt.Sprintf("Erro ao processar arquivo: %v", err),
			"Verifique se o arquivo não está corrompido",
		), nil
	}

	// Calcular tempo de processamento
	processingTime := time.Since(startTime)
	info.ProcessingTime = processingTime.String()

	log.Printf("✅ Arquivo processado com sucesso: %d caracteres em %v", len(text), processingTime)
	return &models.NewSuccessResponse(text, info), nil
}

// GetSupportedTypes retorna tipos de arquivo suportados
func (fs *FileService) GetSupportedTypes() *models.SupportedTypes {
	return &models.SupportedTypes{
		Documents: []string{".pdf", ".txt", ".docx"},
		Images:    []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".tiff"},
		MaxSize:   "5MB",
		MaxSizeBytes: 5 * 1024 * 1024,
	}
}

// Close fecha recursos do serviço
func (fs *FileService) Close() {
	if fs.ocrClient != nil {
		fs.ocrClient.Close()
	}
}
