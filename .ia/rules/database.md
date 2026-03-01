---
description: Convenções de banco de dados (GORM, PostgreSQL, migrações)
globs: internal/repository/**/*.go,internal/platform/database/**/*.go,db/**/*.sql
alwaysApply: false
---

# Banco de Dados

## Acesso

- Obter conexão: `database.DBFromContext(ctx)`
- Transaction-per-request: middleware injeta no context; commit/rollback baseado em status HTTP
- Mutação: `DatabaseWithTransaction`; leitura: `DatabaseWithoutTransaction`

## Repositórios

- Interface `Persistent` no pacote; implementação `datasource` (não exportado)
- `ErrNotFound` para registro inexistente
- `errors.Is(err, gorm.ErrRecordNotFound)` para detectar

## Migrações

- Formato: `NNNNNN_descricao.up.sql` e `.down.sql`
- Timezone: `TIMESTAMP WITH TIME ZONE`
- Índices: `idx_{table}_{column(s)}`
- Soft delete: coluna `deleted_at`

## GORM

- Entidades com tags `gorm:"..."`; `json:"-"` para não expor
- Hooks: `BeforeCreate` (UUID v7), `AfterFind` (UTC)

## Fixtures

- `db/fixtures/*.sql` para testes
- `dbtest.ResetWithFixtures(db, dir, "file.sql")` para isolamento
