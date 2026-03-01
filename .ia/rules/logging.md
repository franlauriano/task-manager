---
description: Convenções de logging (slog, NDJSON)
globs: internal/**/*.go
alwaysApply: false
---

# Logging

## Aplicação

- **slog** como logger padrão
- Configuração via TOML: `[logger] level = "info"`

## Níveis — quando usar

| Nível | Uso |
|-------|-----|
| **error** | Falhas em operações; erros que exigem atenção; algo quebrou e o fluxo falhou |
| **warn** | Situações anômalas mas recuperáveis; fallback usado (ex: cache miss); degradação |
| **info** | Eventos normais do fluxo; início/fim de operações; estado da aplicação (ex: "Database connection opened") |
| **debug** | Detalhes para diagnóstico; desenvolvimento; troubleshooting (não em produção por padrão) |

## Logs com contexto

Usar `slog.ErrorContext`, `slog.InfoContext`, etc. quando a função tiver acesso ao `context.Context`. Permite rastreio de erros. Usar `slog.Error`, `slog.Info` (sem Context) apenas quando o contexto não estiver disponível (ex: inicialização, setup).

```go
// ✅ BOM — com contexto (handler tem r.Context())
slog.ErrorContext(r.Context(), "error decoding JSON body for create task", "error", err)

// ⚠️ Aceitável — sem contexto (função não recebe ctx)
slog.Info("Database connection opened")
```

## log.Fatal

- **Permitido** apenas no fluxo de inicialização da aplicação (ex: `cmd/main.go`), para validar configurações necessárias ao funcionamento: conexão com banco, config obrigatórias, etc.
- **Jamais** usar em fluxo de funcionalidades (handlers, usecases, repositórios). Retornar erro e deixar o chamador decidir.

## Evitar

- Logs em excesso em hot paths
