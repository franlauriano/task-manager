---
description: Padrões React (TypeScript, Vite, TanStack Query, Tailwind)
globs: ui/**/*.{ts,tsx}
alwaysApply: false
---

# Padrões React

## Stack

- React 19, TypeScript, Vite
- TanStack Query (React Query), React Router 7
- Tailwind CSS 4, Lucide React
- Zod para validação

## Estrutura

```
ui/
  api/           — Funções de fetch
  hooks/         — Custom hooks (useQuery, etc.)
  pages/         — Componentes de página
  App.tsx
  main.tsx
```

## Componentes

- Functional components
- Default export para páginas
- Named export para hooks e utilitários

## Dados

- TanStack Query: `useQuery` com `queryKey` e `queryFn`
- Fetch em `api/`; hooks em `hooks/` encapsulam `useQuery`

## Estilo

- Tailwind: classes utilitárias (`flex`, `gap-4`, `rounded-2xl`)
- Path alias: `@/` para `src/` ou raiz do projeto

## Imports

- Ordem: React, externos, internos (`@/`)
