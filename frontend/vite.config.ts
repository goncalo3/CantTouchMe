import path from 'node:path'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig, loadEnv } from 'vite'
import type { UserConfig, ConfigEnv } from 'vite'

// Configuration with environment handling and production defaults
export default defineConfig(({}: ConfigEnv): UserConfig => {
  const env = loadEnv('', process.cwd(), '')
  
  return {
    plugins: [vue(), tailwindcss()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
    define: {
      'import.meta.env.VITE_API_URL': JSON.stringify(
        env.ENVIRONMENT === 'production' 
          ? 'https://canttouchme.goncalo3.pt/api'
          : `http://localhost:${env.API_PORT || 3000}`
      ),
    },
  }
})