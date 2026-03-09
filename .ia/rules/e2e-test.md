---
description: Testes E2E com Playwright
globs: e2e/**/*.spec.ts
alwaysApply: false
---

# Testes E2E (Playwright)

Para pirâmide de testes, escopo por camada e critérios de valor, ver **test-strategy.md**.

## Tecnologias

- **Playwright** — E2E contra UI + API real (ver project-conventions para stack)
- **Não** usar Cypress

## Estrutura de Diretórios

```
e2e/
  fixtures/
    base.ts          — extend test com Page Objects e helpers
  pages/
    tasks.page.ts    — Page Objects por recurso
    teams.page.ts
  support/
    api.ts           — helpers para chamadas à API (setup de dados)
  tasks/
    create.spec.ts   — specs por recurso/operação
    list.spec.ts
  teams/
    create.spec.ts
  global-setup.ts    — chama make seed antes da suite
  playwright.config.ts
```

## Nomenclatura

| Elemento | Padrão | Exemplo |
|----------|--------|---------|
| Arquivo spec | `{recurso}/{op}.spec.ts` | `tasks/create.spec.ts` |
| Page Object | `{recurso}.page.ts` | `tasks.page.ts` |
| Describe | Um por feature por arquivo | `describe('create task', ...)` |
| Test | `should + verbo + resultado` | `should create and display a new task` |

## Seletores (ordem de preferência)

1. **getByRole** — acessível, estável
2. **getByLabel** — formulários
3. **getByText** — conteúdo visível
4. **getByTestId** — último recurso (data-testid)

```ts
// ✅ BOM
await page.getByRole('button', { name: 'New Task' }).click();
await page.getByLabel('Title').fill('My task');
await page.getByRole('button', { name: 'Save' }).click();
await expect(page.getByText('My task')).toBeVisible();

// ❌ EVITAR
await page.locator('.btn-primary').click();
await page.$('#title').fill('My task');
```

## Asserções

- `expect(locator).toBeVisible()`, `toHaveText()`, `toHaveURL()`, `toHaveCount()`
- **Proibido** `page.waitForTimeout()` — usar `expect` com auto-retry do Playwright

## Page Objects e Fixtures

- Page Object: `pages/{resource}.page.ts` com locators e métodos (`goto()`, `create()`)
- Fixture: `fixtures/base.ts` estende `test` com `tasksPage`, `api` etc.
- Imports nos specs: `from '../fixtures/base'`, nunca direto de `@playwright/test`

```ts
// e2e/tasks/create.spec.ts
import { test, expect } from '../fixtures/base';

test('should create and display a new task', async ({ tasksPage, page }) => {
  await tasksPage.goto();
  await tasksPage.create('Buy groceries', 'Milk, bread');
  await expect(page.getByText('Buy groceries')).toBeVisible();
});
```

## Setup

- Pré-requisitos: DB, migrações, API e UI rodando (comandos em **project-conventions**)
- Seed: `make seed` no globalSetup; dados específicos via `support/api.ts`

## Regras

- Testes independentes (sem ordem de execução)
- Setup complexo via API, validação via UI
- **Proibido** `page.evaluate()` salvo necessidade extrema
- Arquivos `.ts` — não usar `.tsx` nos specs
- Payloads e status codes: **api-style.mdc**
