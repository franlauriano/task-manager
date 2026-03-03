# BUG-[NUM]: [Título curto — o que o usuário vê de errado]

## Metadados

- **ID:** BUG-[NUM]
- **Severidade:** Alta / Média / Baixa
- **Status:** Aberto
- **Requisito afetado:** RF-[XX]
- **Arquivo:** `[caminho/do/arquivo.ext:linha]`

## Descrição

[Escreva como uma pessoa de QA descreveria o problema para alguém que nunca viu o sistema. Foque no que o usuário experimenta: o que ele fez, o que apareceu na tela, por que isso é um problema. Evite termos técnicos aqui — guarde-os para a seção "Causa Raiz". Se houver impacto no fluxo completo, mencione.]

## Evidências

![BUG-[NUM]](./screenshots/bug-[num].png)

> [Descrição do que a evidência mostra — ex: "Tela em branco ao acessar /teams" ou "Lista exibindo tasks erradas após clicar em filtro 'To Do'".]

## Ambiente

- Página / Rota: `[ex: /teams/:uuid]`
- Componente / Seção: `[ex: Tela de detalhes do team → lista de tasks]`
- Navegador: [ex: Chrome, Chromium]

## Passos para Reproduzir

1. [Pré-condição, se houver — ex: "Ter pelo menos um team com tasks de status variados"]
2. [Navegar para a tela — ex: "Acessar `/teams/{uuid}`"]
3. [Ação que desencadeia o bug — ex: "Clicar no botão de filtro 'To Do'"]
4. [O que observar — ex: "Verificar a lista de tasks exibida"]

## Comportamento Esperado

[O que o usuário deveria ver ou conseguir fazer. Escreva do ponto de vista do usuário, não do código.]

## Comportamento Atual

[O que acontece de fato. Seja específico: que mensagem apareceu, que dado foi exibido, o que não funcionou. Evite generalidades como "não funciona".]

## Causa Raiz (para o dev)

[Explique brevemente onde está o problema no código. Pode incluir um trecho reduzido para facilitar a localização — mas sem precisar explicar a teoria completa. O dev vai ler o código; aqui basta apontar onde está.]

```[linguagem]
// [arquivo:linha] — trecho com o problema destacado
[código problemático]
```

## Correção Sugerida (para o dev)

```[linguagem]
// trecho corrigido
[código corrigido]
```
