# Tarefa X.0: [Título da Tarefa]

<critical>Ler os arquivos de prd.md e techspec.md da pasta pai (../prd.md e ../techspec.md), se você não ler esses arquivos sua tarefa será invalidada</critical>

## Visão Geral

[Breve descrição da tarefa]

<requirements>
[Lista de requisitos obrigatórios]
</requirements>

## Subtarefas

- [ ] X.1 [Descrição da subtarefa]
- [ ] X.2 [Descrição da subtarefa]

## Detalhes de Implementação

[Seções relevantes da spec técnica **NÃO PRECISA MOSTRAR TODA A IMPLEMENTAÇÃO, APENAS REFERENCIE A techspec.md**]

## Critérios de Sucesso

- [Resultados mensuráveis]
- [Requisitos de qualidade]

## Testes da Tarefa

### Cenários Obrigatórios
- [ ] Happy path (fluxo principal com dados válidos)
- [ ] Validações de entrada (campos obrigatórios, limites de tamanho)
- [ ] Erros esperados (not found, conflito de estado, dados inválidos)

### Cenários Adicionais (obrigatórios quando aplicáveis à tarefa)
> **Regra:** se o cenário é aplicável à funcionalidade sendo implementada, ele **deve** ser coberto por testes. Avaliar quais se aplicam e incluir todos os pertinentes.

- [ ] Valores limite (string vazia, max length 255, zero, negativos)
- [ ] Paginação extrema (page=0, page=99999, limit=0)
- [ ] Caracteres especiais e Unicode
- [ ] Múltiplos fatores combinados (ex: filtro + paginação + valor limite)
- [ ] Transições de estado inválidas (se aplicável)

### NÃO testar
- Funções privadas — serão testadas indiretamente pelas funções públicas que as utilizam
- Comportamento de frameworks/libs de terceiros — assume-se que já foram testados pelos mantenedores

<critical>SEMPRE CRIE E EXECUTE OS TESTES DA TAREFA ANTES DE CONSIDERÁ-LA FINALIZADA</critical>
<critical>Cada teste deve ter um "porquê" claro — se não consegue justificar o valor do teste, não o crie</critical>

## Arquivos relevantes
- [Arquivos relevantes desta tarefa]