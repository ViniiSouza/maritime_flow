# Audit Manager

Serviço responsável por consumir eventos de auditoria do RabbitMQ e registrá-los no OpenSearch (banco NoSQL).

## Descrição

O Audit Manager é um serviço Python que:
- Consome eventos de auditoria da fila do RabbitMQ
- Processa e valida os eventos recebidos
- Armazena os eventos no OpenSearch para análise e auditoria

## Estrutura dos Eventos

O Audit Manager suporta dois tipos de eventos que chegam na fila `audit`:

### 1. Events (H/S -> RabbitMQ)
Eventos de chegada/partida de veículos:

```json
{
  "vehicle_type": "ship" | "helicopter",
  "vehicle_uuid": "823429efabc9283",
  "structure_type": "platform" | "central",
  "structure_uuid": "acb432efab98234d",
  "timestamp": 3094870293,
  "event": "arrived" | "departed",
  "slot_number": 0,
  "tower_id": "adbae9438ff92a"
}
```

### 2. Requests (T -> RabbitMQ)
Requisições de slots das torres:

```json
{
  "vehicle_type": "ship" | "helicopter",
  "vehicle_uuid": "823429efabc9283",
  "structure_type": "platform" | "central",
  "structure_uuid": "acb432efab98234d",
  "timestamp": 3094870293,
  "result": "allowed" | "denied",
  "slot_number": 0
}
```

**Nota:** O campo `timestamp` pode ser um número (Unix timestamp) ou uma string ISO.

## Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```env
# RabbitMQ Configuration
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USERNAME=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_AUDIT_QUEUE=audit
RABBITMQ_EXCHANGE=requests
RABBITMQ_ROUTING_KEY=requests

# OpenSearch Configuration
OPENSEARCH_HOST=localhost
OPENSEARCH_PORT=9200
OPENSEARCH_USERNAME=admin
OPENSEARCH_PASSWORD=FSFcH8Cs26n7
OPENSEARCH_USE_SSL=true
OPENSEARCH_VERIFY_CERTS=false
OPENSEARCH_INDEX=audit-events
```

### Descrição das Variáveis

#### RabbitMQ
- `RABBITMQ_HOST`: Host do servidor RabbitMQ
- `RABBITMQ_PORT`: Porta do servidor RabbitMQ (padrão: 5672)
- `RABBITMQ_USERNAME`: Usuário do RabbitMQ
- `RABBITMQ_PASSWORD`: Senha do RabbitMQ
- `RABBITMQ_AUDIT_QUEUE`: Nome da fila de auditoria
- `RABBITMQ_EXCHANGE`: Nome do exchange (padrão: "requests")
- `RABBITMQ_ROUTING_KEY`: Routing key para binding (padrão: "requests")

#### OpenSearch
- `OPENSEARCH_HOST`: Host do servidor OpenSearch
- `OPENSEARCH_PORT`: Porta do servidor OpenSearch (padrão: 9200)
- `OPENSEARCH_USERNAME`: Usuário do OpenSearch
- `OPENSEARCH_PASSWORD`: Senha do OpenSearch
- `OPENSEARCH_USE_SSL`: Usar SSL para conexão (true/false)
- `OPENSEARCH_VERIFY_CERTS`: Verificar certificados SSL (true/false)
- `OPENSEARCH_INDEX`: Nome do índice onde os eventos serão armazenados

## Instalação

### Requisitos
- Python 3.11 ou superior
- pip

### Instalação Local

1. Clone o repositório e navegue até o diretório:
```bash
cd audit_manager
```

2. Crie um ambiente virtual (recomendado):
```bash
python -m venv venv
source venv/bin/activate  # Linux/Mac
# ou
venv\Scripts\activate  # Windows
```

3. Instale as dependências:
```bash
pip install -r requirements.txt
```

4. Configure as variáveis de ambiente (crie um arquivo `.env`)

5. Execute o serviço:
```bash
python main.py
```

## Execução com Docker

### Build da imagem:
```bash
docker build -t audit-manager .
```

### Execução:
```bash
docker run --env-file .env audit-manager
```

### Execução com docker-compose (exemplo):
```yaml
version: '3.8'
services:
  audit-manager:
    build: .
    env_file:
      - .env
    restart: unless-stopped
    depends_on:
      - rabbitmq
      - opensearch
```

## Funcionamento

1. **Conexão**: O serviço conecta-se ao RabbitMQ e OpenSearch na inicialização
2. **Criação de Índice**: Se o índice não existir no OpenSearch, ele é criado automaticamente com o mapeamento apropriado
3. **Consumo**: O serviço fica aguardando eventos na fila do RabbitMQ
4. **Processamento**: Quando um evento é recebido:
   - O JSON é decodificado e validado
   - Campos obrigatórios são verificados
   - Timestamp é normalizado
   - Um campo `indexed_at` é adicionado com a data/hora de indexação
5. **Armazenamento**: O evento é indexado no OpenSearch
6. **Confirmação**: A mensagem é confirmada no RabbitMQ após processamento bem-sucedido

## Estrutura do Índice OpenSearch

O índice criado automaticamente possui o seguinte mapeamento:

```json
{
  "mappings": {
    "properties": {
      "vehicle_type": {"type": "keyword"},
      "vehicle_uuid": {"type": "keyword"},
      "structure_type": {"type": "keyword"},
      "structure_uuid": {"type": "keyword"},
      "timestamp": {"type": "date"},
      "event": {"type": "keyword"},
      "result": {"type": "keyword"},
      "slot_number": {"type": "integer"},
      "tower_id": {"type": "keyword"},
      "event_type": {"type": "keyword"},
      "indexed_at": {"type": "date"}
    }
  }
}
```

## Tratamento de Erros

- **Erros de conexão**: O serviço tenta reconectar automaticamente
- **Erros de processamento**: Mensagens inválidas são rejeitadas e não reenfileiradas
- **Erros de armazenamento**: Logs detalhados são gerados para diagnóstico

## Logs

O serviço gera logs estruturados com:
- Timestamp
- Nível de log
- Mensagem descritiva
- Detalhes de eventos processados

## Shutdown Graceful

O serviço suporta shutdown graceful através de sinais SIGINT e SIGTERM:
- Encerra o consumo de mensagens
- Finaliza processamento de mensagens em andamento
- Fecha conexões de forma segura

## Desenvolvimento

### Estrutura do Código

- `main.py`: Arquivo principal com a classe `AuditManager`
- `requirements.txt`: Dependências do projeto
- `Dockerfile`: Imagem Docker para containerização

### Dependências Principais

- `pika`: Cliente RabbitMQ para Python
- `opensearch-py`: Cliente oficial do OpenSearch
- `python-dotenv`: Gerenciamento de variáveis de ambiente
- `requests`: Requisições HTTP (usado pelo opensearch-py)

## Troubleshooting

### Erro de conexão com RabbitMQ
- Verifique se o RabbitMQ está rodando
- Confirme as credenciais e host/porta
- Verifique se a fila e exchange existem

### Erro de conexão com OpenSearch
- Verifique se o OpenSearch está rodando
- Confirme as credenciais (usuário: admin, senha: FSFcH8Cs26n7)
- Verifique configurações SSL se aplicável

### Eventos não estão sendo processados
- Verifique se há mensagens na fila do RabbitMQ
- Confirme que o routing key está correto
- Verifique os logs para erros de processamento

## Licença

Este projeto faz parte do sistema Maritime Flow.
