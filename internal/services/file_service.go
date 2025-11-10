package services

import (
    "fmt"
    "io"
    "log"
    "path/filepath"
    "strings"
    "time"

    "backend-fileprocessing/internal/models"
    "backend-fileprocessing/internal/processors"
)

// FileService serviço de processamento de arquivos usando APENAS Google Gemini
type FileService struct {
	geminiService *GeminiService
	processors    map[string]processors.FileProcessor
}

// NewFileService cria novo serviço de arquivos
func NewFileService() *FileService {
	// Inicializar serviço Gemini (OBRIGATÓRIO!)
	geminiService := NewGeminiService()
	if !geminiService.IsAvailable() {
		log.Printf("⚠️ ATENÇÃO: Gemini não disponível - GEMINI_API_KEY não configurada")
		log.Printf("⚠️ Configure GEMINI_API_KEY para processar arquivos")
	} else {
		log.Printf("✅ Gemini configurado e disponível")
	}

	// Mapear processadores por tipo de arquivo - TODOS usam Gemini!
	processorsMap := map[string]processors.FileProcessor{
		".pdf":  processors.NewPDFProcessor(geminiService),
		".png":  processors.NewImageProcessor(geminiService),
		".jpg":  processors.NewImageProcessor(geminiService),
		".jpeg": processors.NewImageProcessor(geminiService),
		".gif":  processors.NewImageProcessor(geminiService),
		".bmp":  processors.NewImageProcessor(geminiService),
		".webp": processors.NewImageProcessor(geminiService),
		".tiff": processors.NewImageProcessor(geminiService),
		".txt":  processors.NewTextProcessor(),
		".docx": processors.NewDocxProcessor(geminiService),
	}

	return &FileService{
		geminiService: geminiService,
		processors:    processorsMap,
	}
}

// ProcessFile processa arquivo e extrai texto
func (fs *FileService) ProcessFile(file io.Reader, filename string, size int64) (models.Response, error) {
	startTime := time.Now()
	fileType := strings.ToLower(filepath.Ext(filename))
	info := models.NewInfo(filename, fileType, size)

	log.Printf("📁 Processando arquivo: %s (%.2f MB)", filename, float64(size)/1024/1024)

	// Verificar se tipo é suportado
	processor, exists := fs.processors[fileType]
    if !exists {
        return models.NewErrorResponse(
            "UNSUPPORTED_FILE_TYPE",
            fmt.Sprintf("Tipo de arquivo não suportado: %s", fileType),
            "Tipos suportados: .pdf, .png, .jpg, .jpeg, .gif, .bmp, .webp, .tiff, .txt, .docx",
        ), nil
    }

	// Processar arquivo
	text, err := processor.Process(file, filename)
    if err != nil {
        return models.NewErrorResponse(
            "PROCESSING_ERROR",
            fmt.Sprintf("Erro ao processar arquivo: %v", err),
            "Verifique se o arquivo não está corrompido",
        ), nil
    }

	// Calcular tempo de processamento
	processingTime := time.Since(startTime)
	info.ProcessingTime = processingTime.String()

	log.Printf("✅ Arquivo processado com sucesso: %d caracteres em %v", len(text), processingTime)
    return models.NewSuccessResponse(text, info), nil
}

// GetSupportedTypes retorna tipos de arquivo suportados
func (fs *FileService) GetSupportedTypes() *models.SupportedTypes {
	return &models.SupportedTypes{
		Documents: []string{".pdf", ".txt", ".docx"},
		Images:    []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".tiff"},
		MaxSize:   "25MB",
		MaxSizeBytes: 25 * 1024 * 1024,
	}
}

// Close fecha recursos do serviço
func (fs *FileService) Close() {
	// Gemini não precisa de cleanup
}
