package handlers

import (
	"net/http"
	"strconv"

	"backend-fileprocessing/internal/models"
	"backend-fileprocessing/internal/services"

	"github.com/gin-gonic/gin"
)

// FileHandler handler para processamento de arquivos
type FileHandler struct {
	fileService *services.FileService
}

// NewFileHandler cria novo handler de arquivos
func NewFileHandler(fileService *services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// ProcessFile processa arquivo enviado
// @Summary Processar arquivo
// @Description Processa arquivo e extrai texto usando OCR ou extração nativa
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Arquivo para processar (PDF, imagem, TXT, DOCX)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/files/process [post]
func (h *FileHandler) ProcessFile(c *gin.Context) {
	// Verificar se há arquivo
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"NO_FILE",
			"Nenhum arquivo foi enviado",
			"Envie um arquivo usando o campo 'file'",
		))
		return
	}
	defer file.Close()

	// Validar tamanho do arquivo (5MB máximo)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if header.Size > maxSize {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse(
			"FILE_TOO_LARGE",
			"Arquivo muito grande",
			"Tamanho máximo permitido: 5MB",
		))
		return
	}

	// Processar arquivo
	response, err := h.fileService.ProcessFile(file, header.Filename, header.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse(
			"PROCESSING_ERROR",
			"Erro interno do servidor",
			err.Error(),
		))
		return
	}

	// Retornar resposta
	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, response)
	}
}

// GetSupportedTypes retorna tipos de arquivo suportados
// @Summary Tipos de arquivo suportados
// @Description Retorna tipos de arquivo suportados pelo serviço
// @Tags files
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/files/supported-types [get]
func (h *FileHandler) GetSupportedTypes(c *gin.Context) {
	supportedTypes := h.fileService.GetSupportedTypes()
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    supportedTypes,
	})
}
