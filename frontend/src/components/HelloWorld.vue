<script setup>
import { ref, onMounted } from 'vue'
import { API_URL, apiUrl } from '../config/api'

defineProps({
  msg: String,
})

const apiStatus = ref('Loading...')
const apiResponse = ref(null)
const loading = ref(true)
const error = ref(null)

// Function to check API status
const checkApiStatus = async () => {
  loading.value = true
  error.value = null
  
  try {
    // Make a direct fetch request to get text response
    const url = apiUrl('/')
    const response = await fetch(url)
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status} ${response.statusText}`)
    }
    
    // Get response as text instead of JSON
    const textResponse = await response.text()
    apiStatus.value = 'Connected'
    apiResponse.value = textResponse
  } catch (err) {
    apiStatus.value = 'Error'
    error.value = err.message
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  checkApiStatus()
})
</script>

<template>
  <h1>{{ msg }}</h1>
  
  <div class="api-status-card">
    <h2>API Status</h2>
    <div class="endpoint">
      <strong>Endpoint:</strong> {{ API_URL }}
    </div>
    
    <div class="status" :class="apiStatus.toLowerCase()">
      <strong>Status:</strong> {{ apiStatus }}
      <button @click="checkApiStatus" :disabled="loading" class="refresh-btn">
        {{ loading ? 'Loading...' : 'Refresh' }}
      </button>
    </div>
    
    <div v-if="error" class="error">
      {{ error }}
    </div>
    
    <div v-if="apiResponse && !error" class="response">
      <h3>Response</h3>
      <div class="text-response">{{ apiResponse }}</div>
    </div>
  </div>
</template>

<style scoped>
.api-status-card {
  background-color: #2b2b2b;
  border-radius: 8px;
  padding: 20px;
  margin: 20px 0;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  max-width: 600px;
}

.endpoint {
  margin-bottom: 10px;
  word-break: break-all;
}

.status {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.status.connected {
  color: #198754;
}

.status.error {
  color: #dc3545;
}

.status.loading {
  color: #6c757d;
}

.refresh-btn {
  background-color: #0d6efd;
  border: none;
  border-radius: 4px;
  color: white;
  padding: 5px 10px;
  margin-left: 10px;
  cursor: pointer;
}

.refresh-btn:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}

.error {
  background-color: #f8d7da;
  border: 1px solid #f5c2c7;
  border-radius: 4px;
  color: #842029;
  padding: 10px;
  margin-bottom: 15px;
}

.response {
  background-color: #2b2b2b;
  border-radius: 4px;
  padding: 15px;
  overflow-x: auto;
  color: #ffffff;
  border: 1px solid #444;
}

.text-response {
  white-space: pre-wrap;
  font-family: 'Courier New', monospace;
  line-height: 1.5;
  color: #ffffff;
  font-size: 14px;
}
</style>
