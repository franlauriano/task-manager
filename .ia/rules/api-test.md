---
description: Testes de API com Venom (YAML)
globs: api_test/**/*.yml internal/transport/**/*_test.go
alwaysApply: false
---

# Testes de API (Venom)

Para pirâmide de testes, escopo por camada e critérios de valor, ver **test-strategy.md**.

## Onde fica a chamada

Os testes Venom são executados em `internal/transport/*_handler_test.go`. Cada `TestXxx` usa `testenv.Setup` com `WithAPITest` e chama `env.RunAPISuite(t, suitePath)` para rodar os YAML. Um teste Go mapeia para uma ou mais suites YAML (ex: `TestCreateTask` → `success/tasks/create/basic.yml`). Para Testcontainers e testenv.Setup, ver **test-strategy.md**.

## Estrutura dos YAML

- Sucesso: `api_test/success/{resource}/{op}/basic.yml`, `edge_cases.yml`, `corner_cases.yml`
- Falha: `api_test/failure/{resource}/{op}/bad_request.yml`, `validation_errors.yml`, `not_found.yml`

## Categorias

| Categoria | Uso |
|-----------|-----|
| basic | Fluxo padrão, valores normais |
| edge_cases | Limites (255 chars, page=99999) |
| corner_cases | Múltiplos fatores, Unicode |
| bad_request | HTTP 400 — JSON inválido, UUID errado |
| validation_errors | HTTP 422 — campos vazios/longos |
| not_found | HTTP 404 |

## Sintaxe

```yaml
vars:
  task_uuid: from: result.bodyjson.uuid
assertions:
  - result.statuscode ShouldEqual 200
  - result.bodyjson.uuid ShouldNotBeNil
  - result.bodyjson ShouldContainKey "title"
```

## Rotas de Mutação

Requirem `Content-Type: application/json` e `Accept: application/json`.
