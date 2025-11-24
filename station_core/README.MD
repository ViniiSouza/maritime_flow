# station_core

Serviço unificado de Plataforma/Central (P/C) do sistema de fluxo marítimo.

## Funcionalidade

O serviço implementa a funcionalidade descrita nos requisitos:
- **Requisito 3**: T pergunta a P se seus helipads/docks estão disponíveis para receber H/S
- **Requisito 4**: T pergunta a C se seus helipads/docks estão disponíveis para receber H/S

O serviço mantém estado em memória sobre disponibilidade de slots e verifica se um slot específico está livre ou em uso.

## Endpoint

### `POST /slots`

Verifica o estado de um slot específico (helipad ou dock).

**Request:**
```json
{
  "slot_number": 0,
  "slot_type": "dock" | "helipad"
}
```

**Response:**
```json
{
  "state": "free" | "in_use"
}
```

**Exemplo:**
```bash
curl -X POST http://localhost:8000/slots \
  -H "Content-Type: application/json" \
  -d '{"slot_number": 0, "slot_type": "helipad"}'
```

## Estado em Memória

O serviço mantém em memória o estado de todos os slots:
- **Helipads**: slots 0, 1, 2 (padrão: todos livres)
- **Docks**: slots 0, 1, 2, 3, 4 (padrão: todos livres)

Estado: `True` = livre, `False` = em uso

## Instalação e Execução

```bash
pip install -r requirements.txt
python main.py
```

O serviço estará disponível em `http://localhost:8000`
