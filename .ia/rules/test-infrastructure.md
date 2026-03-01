---
description: Infraestrutura de testes (Testcontainers, testenv, fixtures)
globs: internal/**/*_test.go
alwaysApply: false
---

# Infraestrutura de Testes (Go)

Infraestrutura para **testes do backend Go** — unitários (entity, usecase), integração (repository) e API (transport com Venom). Não se aplica a testes E2E (Playwright); ver e2e-test.mdc.

## Table-driven (obrigatório)

Todos os testes Go (unitários e API/transport) devem seguir o padrão table-driven: `tests := []struct {...}` com `t.Run()` por caso.

## Testcontainers — um container por pacote

Cada pacote que precisa de DB/Redis cria **um único container** no `TestMain`, compartilhado por todos os testes daquele pacote. Isolamento entre pacotes: `repository/task`, `repository/team` e `transport` rodam em processos separados, cada um com seu próprio PostgreSQL/Redis. Evita interferência e permite paralelismo (`go test -p N`).

## testenv.Setup

Para testes que precisam de DB, Redis, HTTP ou Venom (ex: transport): usar `testenv.Setup(t, opts...)`. Options: `WithDatabase`, `WithRedis`, `WithHTTPServer`, `WithAPITest`. Containers do `TestMain` são passados via `WithDatabase(container)`, `WithRedis(container)`. Cleanup automático via `t.Cleanup()`. Para resetar dados entre subtestes: `dbtest.ResetWithFixtures(env.DB, dir, "fixture.sql")` e `env.FlushRedis()`.
