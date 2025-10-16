package models

import "time"

// Response estrutura para resposta da API
type Response struct {
	Success bool   `json:"success"`
	Data    *Data  `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

// Data dados da resposta
type Data struct {
	Text string `json:"text"`
	Info Info   `json:"info"`
}

// Error estrutura de erro
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Info informações do arquivo processado
type Info struct {
	FileName    string `json:"fileName"`
	FileType    string `json:"fileType"`
	FileSize    int64  `json:"fileSize"`
	ProcessedAt string `json:"processedAt"`
	ProcessingTime string `json:"processingTime,omitempty"`
}

// SupportedTypes tipos de arquivo suportados
type SupportedTypes struct {
	Documents []string `json:"documents"`
	Images    []string `json:"images"`
	MaxSize   string   `json:"maxSize"`
	MaxSizeBytes int64  `json:"maxSizeBytes"`
}

// HealthResponse resposta do health check
type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}

// StatusResponse resposta de status
type StatusResponse struct {
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Environment string  `json:"environment"`
	Features  []string  `json:"features"`
}

// NewSuccessResponse cria nova resposta de sucesso
func NewSuccessResponse(text string, info Info) Response {
	return Response{
		Success: true,
		Data: &Data{
			Text: text,
			Info: info,
		},
	}
}

// NewErrorResponse cria nova resposta de erro
func NewErrorResponse(code, message, details string) Response {
	return Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// NewInfo cria nova estrutura Info
func NewInfo(fileName, fileType string, fileSize int64) Info {
	return Info{
		FileName:    fileName,
		FileType:    fileType,
		FileSize:    fileSize,
		ProcessedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
}
