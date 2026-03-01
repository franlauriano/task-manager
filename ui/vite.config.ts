import path from 'path'
import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import type { IncomingMessage, ServerResponse } from 'http'

function onProxyError(_err: Error, _req: IncomingMessage, res: ServerResponse) {
  res.writeHead(503)
  res.end()
}

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, '.'),
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8090',
        bypass(req) {
          if (req.url && /\.[a-z]+$/i.test(req.url)) return req.url
        },
        onError: onProxyError,
      },
      '/healthcheck': {
        target: 'http://localhost:8090',
        onError: onProxyError,
      },
    },
  },
})
