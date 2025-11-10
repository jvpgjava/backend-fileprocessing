# Backend File Processing

Serviço profissional de processamento de arquivos em Go para extração de texto de PDFs, imagens e outros tipos de arquivo.

## 🏗️ Arquitetura

```
backend-fileprocessing/
├── cmd/
│   └── server/
│       └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── config.go            # Configurações da aplicação
│   ├── handlers/
│   │   ├── file_handler.go      # Handlers para processamento de arquivos
│   │   └── health_handler.go    # Handlers para health checks
│   ├── middleware/
│   │   ├── cors.go              # Middleware CORS
│   │   ├── logger.go            # Middleware de logging
│   │   └── recovery.go          # Middleware de recovery
│   ├── models/
│   │   └── models.go            # Estruturas de dados
│   ├── processors/
│   │   ├── interface.go         # Interface para processadores
│   │   ├── pdf_processor.go     # Processador de PDFs
│   │   ├── image_processor.go   # Processador de imagens
│   │   ├── text_processor.go    # Processador de textos
│   │   └── docx_processor.go    # Processador de DOCX
│   └── services/
│       └── file_service.go      # Serviço de processamento
├── go.mod                       # Dependências Go
├── go.sum                       # Checksums das dependências
├── Makefile                     # Comandos de desenvolvimento
├── Dockerfile                   # Container Docker
├── vercel.json                  # Configuração Vercel
└── README.md                    # Este arquivo
```

## 🚀 Funcionalidades

- **PDF**: Extração de texto nativa (UniPDF) + Google Gemini (GRATUITO!) para PDFs escaneados
- **Imagens**: OCR para PNG, JPG, JPEG, GIF, BMP, WEBP, TIFF
- **Texto**: Leitura direta de arquivos TXT
- **DOCX**: Extração de texto nativa + OCR como fallback
- **API REST**: Interface profissional com versionamento
- **Middleware**: CORS, Logging, Recovery
- **Deploy**: Suporte para Vercel, Railway, Render

## 📋 Requisitos

- Go 1.21+
- Tesseract OCR instalado no sistema
- Bibliotecas de idioma (português e inglês)

### Instalação do Tesseract

#### Windows
```bash
# Via Chocolatey
choco install tesseract

# Ou baixar do site oficial
# https://github.com/UB-Mannheim/tesseract/wiki
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install tesseract-ocr tesseract-ocr-por tesseract-ocr-eng
```

#### macOS
```bash
brew install tesseract tesseract-lang
```

## 🛠️ Desenvolvimento

### Instalação e Execução

1. **Clone o repositório**
```bash
git clone <seu-repositorio>
cd backend-fileprocessing
```

2. **Instale as dependências**
```bash
make deps
```

3. **Execute o servidor**
```bash
make run
```

### Comandos Disponíveis

```bash
make build         # Build da aplicação
make run           # Executar aplicação
make dev           # Executar com hot reload (requer air)
make deps          # Instalar dependências
make test          # Executar testes
make test-coverage # Executar testes com coverage
make lint          # Executar linter
make format        # Formatar código
make clean         # Limpar arquivos de build
make install-tools # Instalar ferramentas de desenvolvimento
make help          # Mostrar ajuda
```

### Instalação de Ferramentas de Desenvolvimento

```bash
make install-tools
```

## 📡 API Endpoints

### Base URL
```
http://localhost:9091/api/v1
```

### Documentação Swagger
```
http://localhost:9091/swagger/index.html
```

### Health Check
```http
GET /health
```

**Resposta:**
```json
{
  "success": true,
  "data": {
    "status": "ok",
    "service": "backend-fileprocessing",
    "version": "1.0.0",
    "timestamp": "2025-10-16T09:30:00Z",
    "uptime": "2h30m15s"
  }
}
```

### Status Detalhado
```http
GET /status
```

**Resposta:**
```json
{
  "success": true,
  "data": {
    "service": "backend-fileprocessing",
    "version": "1.0.0",
    "status": "running",
    "timestamp": "2025-10-16T09:30:00Z",
    "environment": "debug",
    "features": [
      "PDF Processing",
      "Image OCR",
      "Text Extraction",
      "DOCX Support",
      "REST API"
    ]
  }
}
```

### Processar Arquivo
```http
POST /files/process
Content-Type: multipart/form-data
```

**Parâmetros:**
- `file`: Arquivo para processar (máximo 5MB)

**Resposta de Sucesso:**
```json
{
  "success": true,
  "data": {
    "text": "Texto extraído do arquivo...",
    "info": {
      "fileName": "documento.pdf",
      "fileType": ".pdf",
      "fileSize": 1024000,
      "processedAt": "2025-10-16 09:30:00",
      "processingTime": "1.234s"
    }
  }
}
```

**Resposta de Erro:**
```json
{
  "success": false,
  "error": {
    "code": "UNSUPPORTED_FILE_TYPE",
    "message": "Tipo de arquivo não suportado: .xyz",
    "details": "Tipos suportados: .pdf, .png, .jpg, .jpeg, .gif, .bmp, .webp, .tiff, .txt, .docx"
  }
}
```

### Tipos de Arquivo Suportados
```http
GET /files/supported-types
```

**Resposta:**
```json
{
  "success": true,
  "data": {
    "documents": [".pdf", ".txt", ".docx"],
    "images": [".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp", ".tiff"],
    "maxSize": "5MB",
    "maxSizeBytes": 5242880
  }
}
```

## 🔧 Configuração

### Variáveis de Ambiente

- `PORT`: Porta do servidor (padrão: 9091)
- `GIN_MODE`: Modo do Gin (release, debug, test)
- `LOG_LEVEL`: Nível de log (debug, info, warn, error)
- `GEMINI_API_KEY`: **Google Gemini API Key (GRATUITO!)** - Para processar PDFs diretamente

### Configurar Google Gemini (Recomendado!)

O Gemini permite processar PDFs escaneados diretamente, sem precisar converter para imagens primeiro. É **GRATUITO** e funciona muito melhor que OCR tradicional!

1. **Obter API Key:**
   - Acesse: https://makersuite.google.com/app/apikey
   - Faça login com sua conta Google
   - Crie uma nova API Key (gratuito!)

2. **Configurar (Escolha uma opção):**

   **Opção 1: Criar arquivo `.env` (Recomendado para desenvolvimento local)**
   ```bash
   # No diretório backend-fileprocessing, crie um arquivo .env:
   PORT=9091
   GIN_MODE=debug
   LOG_LEVEL=info
   GEMINI_API_KEY=AIzaSyBrMqmufkvDulFkCLu9XzYCOqzmPEz7tFk
   ```
   O Go agora carrega o `.env` automaticamente! ✅

   **Opção 2: Variáveis de ambiente do sistema**
   ```bash
   # Windows PowerShell
   $env:GEMINI_API_KEY="AIzaSyBrMqmufkvDulFkCLu9XzYCOqzmPEz7tFk"
   
   # Linux/macOS
   export GEMINI_API_KEY="AIzaSyBrMqmufkvDulFkCLu9XzYCOqzmPEz7tFk"
   ```

3. **Fluxo de Processamento:**
   ```
   PDF → UniPDF (texto nativo) → Se falhar → Gemini (GRATUITO!) ✅
   ```

### Exemplo de Uso com cURL

```bash
# Health check
curl http://localhost:9091/api/v1/health

# Status detalhado
curl http://localhost:9091/api/v1/status

# Processar arquivo
curl -X POST -F "file=@documento.pdf" http://localhost:9091/api/v1/files/process

# Listar tipos suportados
curl http://localhost:9091/api/v1/files/supported-types

# Acessar documentação Swagger
open http://localhost:9091/swagger/index.html
```

## 🧪 Testes Locais

### Executar Testes

```bash
# Executar todos os testes
make test

# Executar testes com coverage
make test-coverage

# Executar testes específicos
go test -v ./internal/processors/
```

## 🚀 Deploy

### Vercel

1. **Deploy automático**
```bash
vercel --prod
```

2. **Configuração manual**
- Conectar repositório GitHub
- Configurar build command: `go build -o main cmd/server/main.go`
- Configurar output directory: `.`

### Railway

1. **Deploy via GitHub**
- Conectar repositório
- Configurar variáveis de ambiente
- Deploy automático

### Render

1. **Deploy via GitHub**
- Conectar repositório
- Configurar build command: `go build -o main cmd/server/main.go`
- Configurar start command: `./main`

## 🔗 Integração com NewsTrust

Para integrar com o backend do NewsTrust, modifique o `FileProcessor.js`:

```javascript
async processPdfBuffer(buffer) {
    try {
        const formData = new FormData();
        formData.append('file', new Blob([buffer]), 'document.pdf');
        
        const response = await fetch('https://seu-servico-golang.vercel.app/api/v1/files/process', {
            method: 'POST',
            body: formData
        });
        
        const result = await response.json();
        
        if (result.success) {
            return result.data.text;
        } else {
            throw new Error(result.error.message);
        }
    } catch (error) {
        console.error('Erro ao processar PDF:', error);
        return `[PDF] - Erro no processamento: ${error.message}`;
    }
}
```

## 📝 Logs

O serviço gera logs estruturados:
- Inicialização do servidor
- Processamento de arquivos
- Erros de OCR
- Tempo de processamento
- Métricas de performance

## 🚨 Limitações

- Tamanho máximo de arquivo: 5MB
- Timeout de processamento: 30 segundos
- Memória limitada em ambientes serverless
- Tesseract deve estar instalado no sistema


## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 🆘 Suporte

Para suporte, abra uma issue no repositório ou entre em contato com a equipe de desenvolvimento.
