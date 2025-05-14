/**
 * API configuration
 * 
 * This file centralizes API URL configuration based on the current environment.
 * - When ENVIRONMENT=development: Uses http://localhost:3000
 * - When ENVIRONMENT=production: Uses https://canttouchme.goncalo3.pt/api
 * 
 * The API_URL is automatically set in vite.config.js based on the ENVIRONMENT variable.
 */

// Get the API base URL from environment variables
export const API_URL = import.meta.env.API_URL

/**
 * Creates a URL for an API endpoint
 * @param {string} endpoint - The API endpoint without leading slash
 * @returns {string} - The complete API URL
 */
export const apiUrl = (endpoint) => {
  const path = endpoint.startsWith('/') ? endpoint : `/${endpoint}`
  return `${API_URL}${path}`
}

export default {
  API_URL,
  apiUrl
}
