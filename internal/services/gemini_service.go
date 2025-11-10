package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// GeminiService serviço para comunicação com Google Gemini API
type GeminiService struct {
	apiKey       string
	apiURL       string
	modelOptions []string
}

// GeminiRequest estrutura da requisição para Gemini
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent conteúdo para enviar ao Gemini
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart parte do conteúdo (pode ser texto ou dados inline)
type GeminiPart struct {
	InlineData *GeminiInlineData `json:"inline_data,omitempty"`
	Text       string            `json:"text,omitempty"`
}

// GeminiInlineData dados inline (PDF, imagem, etc)
type GeminiInlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"` // Base64
}

// GeminiResponse resposta do Gemini
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate candidato de resposta
type GeminiCandidate struct {
	Content GeminiResponseContent `json:"content"`
}

// GeminiResponseContent conteúdo da resposta
type GeminiResponseContent struct {
	Parts []GeminiResponsePart `json:"parts"`
}

// GeminiResponsePart parte da resposta
type GeminiResponsePart struct {
	Text string `json:"text"`
}

// NewGeminiService cria novo serviço Gemini
func NewGeminiService() *GeminiService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Printf("⚠️ GEMINI_API_KEY não configurada - funcionalidade Gemini desabilitada")
	}

	// Usar gemini-1.5-flash (sem -latest) - modelo correto para processar arquivos
	// Formato: v1beta/models/MODEL_NAME:generateContent
	apiURL := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"
	
	if apiKey != "" {
		log.Printf("✅ Gemini configurado - API URL: %s", apiURL)
		log.Printf("✅ API Key configurada (primeiros 10 chars): %s...", apiKey[:min(10, len(apiKey))])
	} else {
		log.Printf("⚠️ GEMINI_API_KEY não configurada")
	}
	
	return &GeminiService{
		apiKey: apiKey,
		apiURL: apiURL,
		modelOptions: []string{
			"gemini-1.5-flash-latest",
			"gemini-1.5-pro-latest",
			"gemini-1.5-flash",
			"gemini-1.5-pro",
		},
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// IsAvailable verifica se o serviço está disponível
func (s *GeminiService) IsAvailable() bool {
	return s.apiKey != ""
}

// ExtractTextFromPDF extrai texto de PDF usando Gemini
func (s *GeminiService) ExtractTextFromPDF(fileReader io.Reader, filename string) (string, error) {
	if !s.IsAvailable() {
		return "", fmt.Errorf("Gemini não está disponível - GEMINI_API_KEY não configurada")
	}

	log.Printf("🤖 Enviando PDF para Google Gemini (gratuito)...")

	// Ler arquivo completo em buffer
	fileBuffer := new(bytes.Buffer)
	_, err := io.Copy(fileBuffer, fileReader)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo: %v", err)
	}

	// Converter para base64
	base64Content := base64.StdEncoding.EncodeToString(fileBuffer.Bytes())

	// Criar requisição para Gemini
	requestBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						InlineData: &GeminiInlineData{
							MimeType: "application/pdf",
							Data:     base64Content,
						},
					},
					{
						Text: fmt.Sprintf(`Extraia TODO o texto deste PDF (%s) e retorne APENAS o texto extraído, sem comentários ou explicações adicionais. 
Se o PDF contiver imagens escaneadas, descreva o conteúdo das imagens também.

Retorne apenas o texto puro extraído do documento.`, filename),
					},
				},
			},
		},
	}

	// Converter para JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao criar JSON: %v", err)
	}

	// Usar a mesma lógica de tentar múltiplos modelos
	return s.tryRequestWithModels(jsonData, "PDF")
}

// ExtractTextFromFile extrai texto de qualquer arquivo usando Gemini (PDF, imagens, DOCX, etc)
func (s *GeminiService) ExtractTextFromFile(fileReader io.Reader, filename string) (string, error) {
	if !s.IsAvailable() {
		return "", fmt.Errorf("Gemini não está disponível - GEMINI_API_KEY não configurada")
	}

	// Detectar tipo MIME baseado na extensão
	mimeType := getMimeType(filename)
	
	log.Printf("🤖 Enviando arquivo %s (%s) para Google Gemini...", filename, mimeType)

	// Ler arquivo completo em buffer
	fileBuffer := new(bytes.Buffer)
	_, err := io.Copy(fileBuffer, fileReader)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo: %v", err)
	}

	fileSize := fileBuffer.Len()
	log.Printf("📊 Tamanho do arquivo: %d bytes (%.2f MB)", fileSize, float64(fileSize)/1024/1024)

	// Converter para base64
	base64Content := base64.StdEncoding.EncodeToString(fileBuffer.Bytes())
	base64Size := len(base64Content)
	log.Printf("📦 Arquivo convertido para base64: %d caracteres (%.2f MB)", base64Size, float64(base64Size)/1024/1024)

	// Verificar tamanho (Gemini tem limite de ~20MB em base64)
	maxSize := 20 * 1024 * 1024 // 20MB
	if base64Size > maxSize {
		return "", fmt.Errorf("arquivo muito grande para Gemini: %.2f MB (limite: 20MB)", float64(base64Size)/1024/1024)
	}

	// Criar requisição para Gemini
	requestBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						InlineData: &GeminiInlineData{
							MimeType: mimeType,
							Data:     base64Content,
						},
					},
					{
						Text: fmt.Sprintf(`Extraia TODO o texto deste arquivo (%s) e retorne APENAS o texto extraído, sem comentários ou explicações adicionais.
Se o arquivo contiver imagens, descreva o conteúdo das imagens também.
Se for um documento (PDF, DOCX), extraia todo o texto presente.

Retorne apenas o texto puro extraído do documento.`, filename),
					},
				},
			},
		},
	}

	// Converter para JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao criar JSON: %v", err)
	}

	log.Printf("📤 Enviando requisição para Gemini API (tamanho JSON: %d bytes)...", len(jsonData))

	// Tentar diferentes modelos até encontrar um disponível
	return s.tryRequestWithModels(jsonData, "arquivo")
}

// tryRequestWithModels tenta diferentes modelos até encontrar um disponível
func (s *GeminiService) tryRequestWithModels(jsonData []byte, fileType string) (string, error) {
	// Primeiro, tentar listar modelos disponíveis
	availableModels, err := s.listAvailableModels()
	if err == nil && len(availableModels) > 0 {
		log.Printf("✅ Modelos disponíveis encontrados: %v", availableModels)
		modelsToTry := availableModels
		// Tentar com os modelos disponíveis
		return s.tryModels(jsonData, fileType, modelsToTry)
	}
	
	log.Printf("⚠️ Não foi possível listar modelos, tentando lista padrão...")
	// Lista padrão de modelos para tentar (em ordem de preferência)
	modelsToTry := []string{
		"gemini-1.5-flash-latest",  // Versão latest
		"gemini-1.5-pro-latest",    // Versão latest
		"gemini-1.5-flash",         // Versão estável
		"gemini-1.5-pro",           // Versão estável
		"gemini-pro",               // Modelo básico
		"gemini-1.0-pro",           // Modelo mais antigo
		"gemini-1.5-flash-002",     // Versão específica mais recente
		"gemini-1.5-pro-002",       // Versão específica mais recente
	}
	
	return s.tryModels(jsonData, fileType, modelsToTry)
}

// tryModels tenta uma lista específica de modelos
func (s *GeminiService) tryModels(jsonData []byte, fileType string, modelsToTry []string) (string, error) {
	
	var lastErr error
	
	// Tentar com diferentes versões da API
	apiVersions := []string{"v1beta", "v1"}
	
	for _, apiVersion := range apiVersions {
		for _, model := range modelsToTry {
			modelURL := fmt.Sprintf("https://generativelanguage.googleapis.com/%s/models/%s:generateContent", apiVersion, model)
			log.Printf("🔄 Tentando modelo: %s na API %s (para %s)", model, apiVersion, fileType)
			
			// Fazer requisição HTTP
			req, err := http.NewRequest("POST", fmt.Sprintf("%s?key=%s", modelURL, s.apiKey), bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("❌ Erro ao criar requisição HTTP: %v", err)
				lastErr = fmt.Errorf("erro ao criar requisição: %v", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")
			
			// Adicionar contexto com timeout explícito
			ctx := req.Context()
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*60*time.Second)
			defer cancel()
			req = req.WithContext(ctxWithTimeout)

			// Cliente HTTP com timeout de 5 minutos (para processar arquivos grandes)
			client := &http.Client{
				Timeout: 5 * 60 * time.Second, // 5 minutos
			}
			
			log.Printf("📡 Fazendo requisição HTTP para Gemini (timeout: 5 minutos)...")
			log.Printf("📡 URL: %s (modelo: %s, API: %s)", modelURL, model, apiVersion)
			requestStartTime := time.Now()
			resp, err := client.Do(req)
			requestDuration := time.Since(requestStartTime)
			
			if err != nil {
				log.Printf("❌ Erro HTTP ao fazer requisição para Gemini: %v (após %v)", err, requestDuration)
				// Verificar se foi timeout
				if requestDuration >= 5*60*time.Second {
					log.Printf("⚠️ Timeout! Requisição demorou mais de 5 minutos")
				}
				lastErr = fmt.Errorf("erro ao fazer requisição: %v", err)
				continue
			}
			defer resp.Body.Close()

			log.Printf("📥 Resposta do Gemini recebida (status: %d) para modelo %s na API %s (tempo: %v)", resp.StatusCode, model, apiVersion, requestDuration)

			if resp.StatusCode == http.StatusOK {
				// Sucesso! Usar este modelo
				log.Printf("✅ Modelo %s funcionou na API %s!", model, apiVersion)
				// Parsear resposta normalmente abaixo
				return s.parseGeminiResponse(resp, model)
			}
			
			// Se não foi 200, ler o erro mas continuar tentando
			bodyBytes, _ := io.ReadAll(resp.Body)
			errorMsg := string(bodyBytes)
			
			if resp.StatusCode == 404 {
				// Modelo não encontrado nesta versão, continuar tentando
				log.Printf("⚠️ Modelo %s não encontrado na API %s, continuando...", model, apiVersion)
				lastErr = fmt.Errorf("modelo %s não encontrado na API %s", model, apiVersion)
				continue
			}
			
			if resp.StatusCode == 429 {
				// Quota excedida - tentar próximo modelo (pode ser que outro modelo tenha quota disponível)
				log.Printf("⚠️ Cota excedida para modelo %s na API %s, tentando próximo modelo...", model, apiVersion)
				// Parsear mensagem de retry
				var errorResp struct {
					Error struct {
						Message string `json:"message"`
						Details []struct {
							RetryInfo struct {
								RetryDelay string `json:"retryDelay"`
							} `json:"retryInfo"`
						} `json:"details"`
					} `json:"error"`
				}
				if err := json.Unmarshal(bodyBytes, &errorResp); err == nil {
					retryDelay := "alguns segundos"
					if len(errorResp.Error.Details) > 0 && errorResp.Error.Details[0].RetryInfo.RetryDelay != "" {
						retryDelay = errorResp.Error.Details[0].RetryInfo.RetryDelay
					}
					lastErr = fmt.Errorf("cota excedida para modelo %s. Tente novamente em %s", model, retryDelay)
				} else {
					lastErr = fmt.Errorf("cota excedida para modelo %s", model)
				}
				continue // Tentar próximo modelo
			}
			
			// Outro erro (400, 403, etc) - parar e retornar
			log.Printf("❌ Erro da API Gemini (status %d) com modelo %s na API %s: %s", resp.StatusCode, model, apiVersion, errorMsg)
			
			var errorResp struct {
				Error struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Status  string `json:"status"`
				} `json:"error"`
			}
			
			if err := json.Unmarshal(bodyBytes, &errorResp); err == nil && errorResp.Error.Message != "" {
				return "", fmt.Errorf("erro da API Gemini: %s (status: %s, code: %d)", errorResp.Error.Message, errorResp.Error.Status, errorResp.Error.Code)
			}
			
			return "", fmt.Errorf("erro da API Gemini (status %d): %s", resp.StatusCode, errorMsg)
		}
	}
	
	// Se chegou aqui, nenhum modelo funcionou
	if strings.Contains(lastErr.Error(), "cota excedida") || strings.Contains(lastErr.Error(), "quota") {
		return "", fmt.Errorf("cota gratuita do Gemini foi excedida. Por favor: 1) Aguarde alguns minutos e tente novamente, 2) Verifique sua cota em https://ai.dev/usage?tab=rate-limit, 3) Considere usar uma API Key diferente ou upgrade do plano. Último erro: %v", lastErr)
	}
	return "", fmt.Errorf("nenhum modelo Gemini disponível. Último erro: %v", lastErr)
}

// listAvailableModels lista os modelos disponíveis na API
func (s *GeminiService) listAvailableModels() ([]string, error) {
	if !s.IsAvailable() {
		return nil, fmt.Errorf("Gemini não está disponível")
	}
	
	// Endpoint para listar modelos
	listURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models?key=%s", s.apiKey)
	
	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("⚠️ Erro ao listar modelos (status %d): %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("erro ao listar modelos: status %d", resp.StatusCode)
	}
	
	var modelsResponse struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&modelsResponse); err != nil {
		return nil, err
	}
	
	var modelNames []string
	// Modelos prioritários (gratuitos e simples)
	priorityModels := []string{
		"gemini-flash-latest",      // Modelo gratuito mais rápido
		"gemini-pro-latest",        // Modelo gratuito básico
		"gemini-2.0-flash",         // Flash 2.0 (gratuito)
		"gemini-2.0-flash-lite",    // Flash Lite (mais barato)
	}
	
	// Primeiro, adicionar modelos prioritários se estiverem na lista
	foundModels := make(map[string]bool)
	for _, model := range modelsResponse.Models {
		name := strings.TrimPrefix(model.Name, "models/")
		foundModels[name] = true
	}
	
	// Adicionar modelos prioritários primeiro
	for _, priority := range priorityModels {
		if foundModels[priority] {
			modelNames = append(modelNames, priority)
		}
	}
	
	// Depois, adicionar outros modelos gemini que não sejam de embedding/imagem
	for _, model := range modelsResponse.Models {
		name := strings.TrimPrefix(model.Name, "models/")
		// Filtrar: apenas modelos gemini que não sejam embedding, imagem, ou outros tipos especiais
		if strings.HasPrefix(name, "gemini-") && 
		   !strings.Contains(name, "embedding") &&
		   !strings.Contains(name, "image") &&
		   !strings.Contains(name, "imagen") &&
		   !strings.Contains(name, "text-embedding") &&
		   !strings.Contains(name, "aqa") &&
		   !strings.Contains(name, "robotics") &&
		   !strings.Contains(name, "computer-use") &&
		   !foundModels[name] { // Evitar duplicatas
			// Adicionar no final (depois dos prioritários)
			modelNames = append(modelNames, name)
		}
	}
	
	return modelNames, nil
}

// parseGeminiResponse parseia a resposta do Gemini
func (s *GeminiService) parseGeminiResponse(resp *http.Response, modelName string) (string, error) {

	// Parsear resposta
	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", fmt.Errorf("erro ao parsear resposta: %v", err)
	}

	// Extrair texto da resposta
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("resposta do Gemini não contém texto")
	}

	extractedText := geminiResp.Candidates[0].Content.Parts[0].Text
	extractedText = strings.TrimSpace(extractedText)

	if len(extractedText) < 10 {
		return "", fmt.Errorf("Gemini extraiu pouco texto (menos de 10 caracteres)")
	}

	log.Printf("✅ Gemini extraiu texto: %d caracteres", len(extractedText))
	return extractedText, nil
}

// getMimeType retorna MIME type baseado na extensão do arquivo
func getMimeType(filename string) string {
	ext := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(ext, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(ext, ".png"):
		return "image/png"
	case strings.HasSuffix(ext, ".jpg") || strings.HasSuffix(ext, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(ext, ".gif"):
		return "image/gif"
	case strings.HasSuffix(ext, ".bmp"):
		return "image/bmp"
	case strings.HasSuffix(ext, ".webp"):
		return "image/webp"
	case strings.HasSuffix(ext, ".tiff") || strings.HasSuffix(ext, ".tif"):
		return "image/tiff"
	case strings.HasSuffix(ext, ".docx"):
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}
