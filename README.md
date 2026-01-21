# Task Management API

API REST desenvolvida em Go para gerenciamento de tarefas (tasks) e equipes (teams), seguindo princípios de Clean Architecture e arquitetura em camadas.

## Índice

- [Documentação](#documentação)
- [Sobre o Projeto](#sobre-o-projeto)
- [Tecnologias](#tecnologias)
- [Pré-requisitos](#pré-requisitos)
- [Instalação](#instalação)
- [Configuração](#configuração)
- [Como Executar](#como-executar)
- [Testes](#testes)
- [Comandos Úteis](#comandos-úteis)

## Documentação

Documentação adicional do projeto está disponível na pasta `doc/`:

- **[Arquitetura](doc/architecture.md)**: Documentação completa sobre arquitetura do sistema, estrutura de diretórios, fluxo de dados e padrões utilizados.

## Sobre o Projeto

Esta API permite gerenciar tarefas e equipes, com funcionalidades como:

- **Tarefas (Tasks)**: Criação, listagem, atualização, exclusão e gerenciamento de status
- **Equipes (Teams)**: Criação, listagem, recuperação e associação/desassociação de tarefas
- **Status de Tarefas**: Estados `to_do`, `in_progress`, `done` e `canceled`
- **Relacionamentos**: Tarefas podem ser associadas a equipes
- **Paginação**: Suporte a paginação em listagens
- **Soft Delete**: Exclusão lógica de registros

## Tecnologias

- **Linguagem**: Go 1.25
- **Framework Web**: Chi (go-chi/chi)
- **ORM**: GORM
- **Banco de Dados**: PostgreSQL 18
- **Configuração**: pelletier/go-toml/v2 (TOML parsing com expansão de variáveis de ambiente)
- **Logging**: slog (structured logging)
- **Testes**: Go testing package, Testcontainers (PostgreSQL em testes) e Venom (testes de API)
- **Containerização**: Docker Compose (postgres, migrate)
- **Migrações**: SQL direto (up/down)

## Pré-requisitos

- Go 1.25 ou superior
- Docker e Docker Compose
- Make (opcional, mas recomendado)
- Air (para live reload durante desenvolvimento)

## Instalação

1. Clone o repositório:
```bash
git clone <repository-url>
cd task-manager
```

2. Instale as dependências:
```bash
make deps
# ou
go mod download
go mod tidy
```

> **Dica**: Execute `make help` para ver todos os comandos disponíveis.

3. Instale o Air (para desenvolvimento com live reload):
```bash
go install github.com/air-verse/air@v1.64.0
```

## Configuração

A aplicação utiliza arquivos de configuração TOML e variáveis de ambiente.

**Para a aplicação:**
- `etc/config.toml`: Configuração principal (criado a partir de `etc/config.toml.example` se não existir)
- `etc/.env`: Variáveis de ambiente (criado a partir de `etc/.env.exemple` se não existir). Ajuste `DATABASE_*`, `SERVER_*`, etc.

**Para testes:**
- `etc/config_test.toml`: Configuração TOML para testes
- `etc/.env.test`: Variáveis de ambiente para testes (carregadas automaticamente antes do TOML)

A aplicação utiliza `pelletier/go-toml/v2` para parsing de arquivos TOML com suporte a expansão de variáveis de ambiente:
- `${VAR_NAME}`: Expande para o valor da variável de ambiente
- `${VAR_NAME:-default}`: Expande para o valor da variável de ambiente, ou usa `default` se a variável não estiver definida

## Como Executar

### 1. Subir o banco de dados

```bash
make db-up
```

### 2. Executar migrações

```bash
make migrate
```

### 3. (Opcional) Popular com dados de desenvolvimento

```bash
make seed
```

### 4. Executar a aplicação

**Modo desenvolvimento (com live reload):**
```bash
make run-dev
```

**Modo desenvolvimento (sem live reload):**
```bash
make run
```

**Compilar binário:**
```bash
make build
```

A API estará disponível em `http://localhost:8090`

## Testes

Os testes utilizam **Testcontainers** para subir um PostgreSQL automaticamente. Não é necessário subir banco de testes nem executar migrações manualmente.

### Testes unitários e de integração

Execute os testes:

```bash
make test
```

### Cobertura de Testes

Gere relatório de cobertura:

```bash
make coverage
```

Isso irá:
- Executar os testes com cobertura
- Gerar `var/coverage.out` e `var/coverage.html`
- Exibir a cobertura total no terminal

> **Nota**: Para informações sobre testes de integração, consulte a [documentação de arquitetura](doc/architecture.md).

## Comandos Úteis

O projeto inclui um `Makefile` com os seguintes comandos:

| Comando | Descrição |
|---------|-----------|
| `make help` | Mostra ajuda com todos os comandos disponíveis |
| `make deps` | Baixa e organiza dependências |
| `make build` | Compila o binário da aplicação |
| `make clean` | Remove arquivos gerados (binários, coverage, etc.) |
| `make db-up` | Sobe o banco de dados PostgreSQL da aplicação |
| `make db-down` | Para o banco de dados da aplicação |
| `make migrate` | Executa migrações do banco da aplicação |
| `make migrate-down` | Reverte migrações do banco da aplicação |
| `make seed` | Executa o seed (`db/seed/populate.sql`) via Docker; depende de `migrate` |
| `make run` | Executa a aplicação (sem live reload) |
| `make run-dev` | Executa a aplicação com live reload (requer Air) |
| `make test` | Executa testes unitários e de integração (usa Testcontainers) |
| `make coverage` | Gera relatório de cobertura |

## Licença

Este projeto é proprietário.