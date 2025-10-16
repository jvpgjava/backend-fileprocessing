#!/bin/bash

echo "🧪 Testando API Backend File Processing..."
echo "=========================================="

# Verificar se o servidor está rodando
echo "1. Verificando se o servidor está rodando..."
if ! curl -s http://localhost:9091/api/v1/health > /dev/null; then
    echo "❌ Servidor não está rodando. Execute 'make run' primeiro."
    exit 1
fi
echo "✅ Servidor está rodando!"

echo -e "\n2. Testando health check..."
curl -s http://localhost:9091/api/v1/health | jq .

echo -e "\n3. Testando status detalhado..."
curl -s http://localhost:9091/api/v1/status | jq .

echo -e "\n4. Testando tipos suportados..."
curl -s http://localhost:9091/api/v1/files/supported-types | jq .

echo -e "\n5. Criando arquivo de teste..."
echo "Este é um arquivo de teste para verificar se a API está funcionando corretamente." > test.txt

echo -e "\n6. Testando processamento de arquivo de texto..."
curl -X POST -F "file=@test.txt" http://localhost:9091/api/v1/files/process | jq .

# Limpar arquivo de teste
rm test.txt

echo -e "\n✅ Testes concluídos!"
echo "Para testar com PDFs ou imagens, use:"
echo "curl -X POST -F \"file=@seu-arquivo.pdf\" http://localhost:9091/api/v1/files/process"
echo -e "\n📚 Documentação Swagger disponível em:"
echo "http://localhost:9091/swagger/index.html"
