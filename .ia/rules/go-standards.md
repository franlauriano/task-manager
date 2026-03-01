---
description: Padrões e estilo de código Go (sem HTTP, testes ou DB)
globs: **/*.go
alwaysApply: false
---

# Padrões Go

## Error Handling

```go
// ✅ BOM — wrapping nas fronteiras
return fmt.Errorf("task repo create: %w", err)

// ✅ BOM — comparação segura (compatível com wrapping)
if errors.Is(err, errs.ErrNotFound) { ... }
var validErr *errors.ValidationErrors
if errors.As(err, &validErr) { ... }

// ❌ RUIM
return err
if err == errs.ErrNotFound { ... }
```

## Imports

- Stdlib primeiro; linha em branco; terceiros; linha em branco; `taskmanager/...`
- Alias para evitar colisão: `taskEntity "taskmanager/internal/entity/task"`, `errs "taskmanager/internal/platform/errors"`

## Structs e Validação

- Validação em métodos `Validate()`, `ValidateTransitionTo()` na entity
- Hooks GORM: `BeforeCreate` (UUID v7), `AfterFind` (UTC)
- Tags GORM em structs de entidade; `json:"-"` para não expor internals

## Nomenclatura

- Pacotes: singular, lowercase (`task`, `dto`, `middleware`)
- Interfaces: substantivo ou verbo passado (`Persistent`, `TaskRepository`)
- Constantes: PascalCase com prefixo (`StatusTodo`, `StatusInProgress`)
- Variáveis: camelCase (`taskUUID`, `statusFilter`)

## Context

- Passar `context.Context` como primeiro parâmetro em operações de I/O
- Nunca armazenar context em structs

## Evitar

- `log.Fatal` em fluxo de funcionalidades — permitido só na inicialização (config, conexão DB). Retornar erro.
