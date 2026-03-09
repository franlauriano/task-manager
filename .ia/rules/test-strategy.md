---
description: Referência central de testes — backend Go (unitário, integração, API) + E2E Playwright, escopo por camada, valor do teste, infraestrutura
globs:
  - internal/**/*_test.go
  - api_test/**/*.yml
  - e2e/**/*.spec.ts
alwaysApply: false
---

# Estratégia de Testes

## Stack por tipo

| Stack | Tipos de teste |
|-------|----------------|
| **Backend Go** | Unitário, Integração, API |
| **Frontend** | E2E (Playwright) — UI + API real |

## Pirâmide

| Tipo | Stack | Onde | Uso |
|------|-------|------|-----|
| **Unitário** | Go | `internal/entity/**/*_test.go`, `internal/usecase/**/*_test.go` | Entity, Usecase (mocks), lógica isolada |
| **Integração** | Go | `internal/repository/**/*_test.go` | Repository com DB real (Testcontainers) |
| **API** | Go (Venom) | `api_test/**/*.yml` + `internal/transport/*_handler_test.go` | Contratos HTTP, status codes, DTOs |
| **E2E** | TypeScript (Playwright) | `e2e/**/*.spec.ts` | Fluxos críticos de usuário, UI + API |

Regras específicas: **unit-test.md**, **api-test.md**, **e2e-test.md**.

## O que testar por camada (backend Go)

| Camada | Testar | NÃO testar |
|--------|--------|------------|
| **Entity** | Validações (`Validate()`), regras de negócio, transições de estado (`ValidateTransitionTo`), formatação | Funções privadas (cobertas pelas públicas), comportamento de libs (GORM hooks, tags) |
| **Usecase** | Orquestração, tratamento de erros (not found, validação), lógica condicional | Funções privadas (cobertas pelas públicas), comportamento de libs |
| **Repository** | Queries complexas, filtros, paginação, constraints do banco | CRUD básico provido pelo GORM, conexão com banco |
| **Transport** | Parsing de request, status codes, formato de response (DTOs), validação de headers | Routing do chi, middleware padrão, serialização JSON da stdlib |

## Quando NÃO criar um teste

- **Funções privadas**: serão testadas indiretamente pelas funções públicas que as utilizam
- **Frameworks/libs de terceiros**: assume-se que já foram testados pelos mantenedores (GORM, chi, slog, etc.)
- **Duplicação de camada**: se um teste de integração/API já cobre o cenário end-to-end, não duplique com unitário do mesmo caminho

## Como identificar um teste de alto valor

Um teste é **valioso** quando:
1. Protege contra regressões em lógica de negócio
2. Documenta um comportamento não-óbvio (regra complexa, cenário de erro)
3. Valida um cenário que já causou bug ou é propenso a erro

Um teste é de **baixo valor** quando:
1. Testa comportamento já garantido por libs/frameworks
2. Duplica exatamente o que outro teste já cobre
3. Está acoplado à implementação (quebra quando refatora sem mudar comportamento)

---

## Infraestrutura (testes Go)

Infraestrutura para **testes do backend Go** — unitários (entity, usecase), integração (repository) e API (transport com Venom). Não se aplica a testes E2E (Playwright); ver **e2e-test.md**.

### Table-driven (obrigatório)

Todos os testes Go (unitários e API/transport) devem seguir o padrão table-driven: `tests := []struct {...}` com `t.Run()` por caso.

### Testcontainers — um container por pacote

Cada pacote que precisa de DB/Redis cria **um único container** no `TestMain`, compartilhado por todos os testes daquele pacote. Isolamento entre pacotes: `repository/task`, `repository/team` e `transport` rodam em processos separados, cada um com seu próprio PostgreSQL/Redis. Evita interferência e permite paralelismo (`go test -p N`).

### testenv.Setup

Para testes que precisam de DB, Redis, HTTP ou Venom (ex: transport): usar `testenv.Setup(t, opts...)`. Options: `WithDatabase`, `WithRedis`, `WithHTTPServer`, `WithAPITest`. Containers do `TestMain` são passados via `WithDatabase(container)`, `WithRedis(container)`. Cleanup automático via `t.Cleanup()`. Para resetar dados entre subtestes: `dbtest.ResetWithFixtures(env.DB, dir, "fixture.sql")` e `env.FlushRedis()`.
