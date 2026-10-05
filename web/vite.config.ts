import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath } from 'node:url'
import { defineConfig, loadEnv } from 'vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const proxy = (target: string, prefix: string) => ({
    target,
    changeOrigin: true,
    rewrite: (path: string) => path.replace(new RegExp(`^${prefix}`), ''),
  })

  return {
    plugins: [react(), tailwindcss()],
    resolve: {
      alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
    },
    server: {
      proxy: {
        '/b2b': proxy(env.B2B_PROXY_TARGET, '/b2b'),
        '/mis': proxy(env.MIS_PROXY_TARGET, '/mis'),
        '/b2c': proxy(env.B2C_PROXY_TARGET, '/b2c'),
      },
    },
  }
})
