---
description: Convenções da API REST (status codes, request/response body, Chi, DTOs)
globs: internal/transport/**/*.go
alwaysApply: false
---

# Estilo da API REST

## Princípios

- **DTOs em responses**: Sempre construir respostas via DTOs (`ToTaskResponse`, `ToPaginatedTasksResponse`, etc.). Nunca expor entidades diretamente. Garante que só saem os dados desejados e evita vazamento de campos internos.
- **UUID em vez de ID**: A API usa apenas UUID (v7) como identificador público. IDs internos (auto-increment) não são expostos. Evita ataques de enumeração (adivinhar IDs sequenciais para descobrir recursos).

## Status Codes

| Código | Uso |
|--------|-----|
| 200 | Sucesso em todas as operações (create, update, delete, list, retrieve) |
| 400 | Bad Request — JSON inválido, UUID inválido, campo obrigatório ausente |
| 404 | Not Found — recurso inexistente |
| 415 | Unsupported Media Type — Content-Type diferente de application/json |
| 422 | Unprocessable Entity — validação de domínio (ValidationErrors) |
| 500 | Internal Server Error — erro inesperado |

## Request Body (JSON)

### Tasks

| Endpoint | Body |
|----------|------|
| POST /api/tasks | `{ "title": string, "description": string }` |
| PUT /api/tasks/{uuid} | `{ "title": string, "description": string }` |
| POST /api/tasks/{uuid}/status | `{ "status": "to_do" | "in_progress" | "done" | "canceled" }` |

### Teams

| Endpoint | Body |
|----------|------|
| POST /api/teams | `{ "name": string, "description": string }` |
| POST /api/teams/{uuid}/tasks | `{ "task_uuid": string }` |
| DELETE /api/teams/{uuid}/tasks/{task_uuid} | (sem body) |

### Headers

- Mutação: `Content-Type: application/json` obrigatório
- Recomendado: `Accept: application/json`

## Response Body (JSON)

### Sucesso — Recurso único

```json
{
  "uuid": "uuid-v7-string",
  "title": "string",
  "description": "string",
  "status": "to_do" | "in_progress" | "done" | "canceled",
  "started_at": "RFC3339 ou null",
  "finished_at": "RFC3339 ou null",
  "created_at": "RFC3339",
  "updated_at": "RFC3339"
}
```

Team: `uuid`, `name`, `description`, `created_at`, `updated_at`. Com tasks: inclui `tasks: []`.

### Sucesso — Lista paginada

```json
{
  "page": 1,
  "items_per_page": 10,
  "total_items": 42,
  "total_pages": 5,
  "items": [ /* array de recursos */ ]
}
```

### Sucesso — Sem body

Delete, UpdateStatus, AssociateTask, DisassociateTask retornam 200 com body vazio `[]`.

### Erro 400 (BadRequestError)

```json
{
  "message": "string",
  "field": "string (opcional)"
}
```

### Erro 422 (ValidationErrors)

```json
{
  "errors": [
    {
      "field": "title",
      "message": "title is required",
      "type": "required",
      "code": "REQUIRED",
      "params": { "max": 255 }
    }
  ]
}
```

| Campo | Significado |
|-------|-------------|
| field | Nome do campo que falhou na validação |
| message | Mensagem legível para o usuário |
| type | Tipo do erro: `required`, `invalid`, `format`, `max_length`, etc. |
| code | Código único para identificação programática (ex: `REQUIRED`, `MAX_LENGTH`) |
| params | Parâmetros extras (ex: `{ "max": 255 }` para limite de caracteres) |

### Erro 404, 500

Body vazio (null).

### Erro 415

```json
{
  "message": "Content-Type must be application/json"
}
```

## Query Params

| Param | Uso | Default |
|-------|-----|---------|
| page | Paginação | 1 |
| limit | Itens por página | config (ex: 10) |
| status | Filtro (tasks): to_do, in_progress, done, canceled | (todos) |

## Handlers e Rotas

- Assinatura: `(int, []byte)`; middleware escreve na response
- Prefixo `/api`; mutação com `RequireContentTypeJSON` + `DatabaseWithTransaction`
- Path params: `{uuid}`, `{task_uuid}`
