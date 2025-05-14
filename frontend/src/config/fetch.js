/**
 * Simple API utility to make HTTP requests
 * This is a basic implementation. In a real app, you might want to:
 * - Add authentication headers
 * - Handle token refreshing
 * - Add error logging/reporting
 * - etc.
 */
import { API_URL, apiUrl } from './api'

/**
 * Generic fetch wrapper with error handling
 * @param {string} endpoint - API endpoint
 * @param {Object} options - Fetch options
 * @returns {Promise<any>} - Response data
 */
export async function fetchApi(endpoint, options = {}) {
  const url = apiUrl(endpoint)
  
  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
      // You can add authorization headers here if needed
    },
  }
  
  try {
    const response = await fetch(url, { ...defaultOptions, ...options })
    
    if (!response.ok) {
      throw new Error(`API error: ${response.status} ${response.statusText}`)
    }
    
    // Parse JSON if the response has content
    if (response.status !== 204) { // 204 = No Content
      return await response.json()
    }
    
    return null
  } catch (error) {
    // Log the error and rethrow
    console.error('API request failed:', error)
    throw error
  }
}

/**
 * Convenience methods for common HTTP methods
 */
export const api = {
  get: (endpoint) => fetchApi(endpoint, { method: 'GET' }),
  
  post: (endpoint, data) => fetchApi(endpoint, {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  
  put: (endpoint, data) => fetchApi(endpoint, {
    method: 'PUT',
    body: JSON.stringify(data),
  }),
  
  patch: (endpoint, data) => fetchApi(endpoint, {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
  
  delete: (endpoint) => fetchApi(endpoint, { method: 'DELETE' }),
}

export default api
