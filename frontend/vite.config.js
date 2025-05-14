import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig(({ command, mode }) => {
  // Load env variables based on mode (development or production)
  const env = loadEnv(mode, process.cwd(), '')
  
  // Determine API URL based on environment
  const apiUrl = env.ENVIRONMENT === 'development' 
    ? 'http://localhost:3000'
    : 'https://canttouchme.goncalo3.pt/api'

  return {
    plugins: [vue()],
    define: {
      // Make API URL available in client code based on environment
      'import.meta.env.API_URL': JSON.stringify(apiUrl)
    }
  }
})
