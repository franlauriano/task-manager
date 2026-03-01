---
description: Visão geral e convenções do projeto Task Manager
alwaysApply: true
---

# Task Manager — Convenções do Projeto

## Arquitetura

Clean Architecture: Transport → Usecase → Repository → Entity. Dependências fluem para dentro.

## Estrutura de Diretórios

```
cmd/                — Entry point
internal/
  entity/           — Modelos de domínio
  usecase/          — Lógica de negócio
  repository/       — Persistência
  transport/        — HTTP, handlers, DTOs
  platform/         — DB, cache, logger, http utils, errors
  config/
  paths/
db/
  migrate/          — Migrações SQL (up/down)
  seed/             — Dados iniciais (desenvolvimento)
  fixtures/         — Dados para testes
api_test/           — Testes Venom (YAML)
e2e/                — Testes E2E (Playwright)
ui/                 — Frontend React (Vite, TypeScript)
  api/              — Chamadas HTTP à API
  hooks/            — Custom hooks (TanStack Query)
  pages/            — Componentes de página
  App.tsx, main.tsx — Entry point
etc/                — Configuração (.env, config.toml)
```

## Regra de Dependências

- **Entity**: sem deps internas
- **Repository**: entity + platform/database
- **Usecase**: entity + repository
- **Transport**: usecase + entity (tipos) + platform/http
- **Platform**: biblioteca interna; sem imports de negócio

## Comandos

### Build e dependências

| Comando | Descrição |
|---------|-----------|
| `make deps` | Baixa e organiza dependências Go |
| `make build` | Compila binário em `bin/task-api` |
| `make clean` | Remove binários e arquivos de cobertura |

### Execução

| Comando | Descrição |
|---------|-----------|
| `make run` | Roda a API (sem live reload) |
| `make run-dev` | Roda com Air (live reload) |
| `make run-ui` | Roda o frontend React (Vite) |

### Infraestrutura (Docker)

| Comando | Descrição |
|---------|-----------|
| `make db-up` | Sobe PostgreSQL |
| `make db-down` | Para PostgreSQL |
| `make redis-up` | Sobe Redis |
| `make redis-down` | Para Redis |
| `make run-docker` | Sobe PostgreSQL + Redis |

### Banco de dados

| Comando | Descrição |
|---------|-----------|
| `make migrate` | Executa migrações (up) |
| `make migrate-down` | Reverte migrações |
| `make seed` | Carrega seed (requer db-up + migrate) |

### Testes

| Comando | Descrição |
|---------|-----------|
| `make test` | Roda testes (`-tags=test`) |
| `make coverage` | Gera relatório em `var/coverage.html` |

### Outros

| Comando | Descrição |
|---------|-----------|
| `make help` | Lista todos os comandos |

## Nomenclatura

- Handlers: `CreateXxx`, `UpdateXxx`, `ListXxx`, `RetrieveByUUID`
- Pacotes: singular (`task`, `team`, `dto`)
- Migrações: `NNNNNN_descricao.up.sql` / `.down.sql`
