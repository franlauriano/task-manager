---
description: Testes unitários Go (entity, usecase, mocks)
globs: internal/**/*_test.go
alwaysApply: false
---

# Testes Unitários Go

## Build Tag

```go
//go:build test
```

Rodar: `go test -tags=test ./...`

## Table-Driven (obrigatório)

Todos os testes devem seguir o padrão table-driven.

```go
tests := []struct {
    name    string
    input   X
    wantErr error  // usar error, nunca bool
}{
    {"caso sucesso", validInput, nil},
    {"caso vazio", emptyInput, &errs.ValidationErrors{...}},
}
for _, tc := range tests {
    t.Run(tc.name, func(t *testing.T) {
        // ...
    })
}
```

## wantErr

Funções que retornam erro: usar `wantErr error` no struct de teste. **Não** usar `wantErr bool` — o tipo `error` permite comparar com `nil` e com erros específicos (`assert.CompareErrors`).

## Helpers e Cleanup

- `t.Helper()` em helpers para stack trace correto
- `t.Cleanup(fn)` para cleanup (ordem LIFO)
- Usecase com mocks: `defer taskRepo.SetPersist(originalPersist)` após `SetPersist(mock)`

## Asserções

- Usar `assert.CompareErrors(got, want)` para erros
- `cmp.Diff` para structs (google/go-cmp)

## Mocks

- Repository: `SetPersist(&MockPersistent{FnCreate: ...})`
- Restaurar original no final do teste

---

Para escopo por camada, infraestrutura (Testcontainers, testenv) e critérios de valor, ver **test-strategy.md**.
