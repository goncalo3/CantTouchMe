import axios from 'axios';
import router from '@/router';
import { userStore } from '@/store/userStore';
import { noteTitleStore } from '@/store/noteTitleStore';
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  withCredentials: true,
});

// response interceptor: handle 401 errors
api.interceptors.response.use(
  response => response,
  (error) => {
    if (error.response?.status === 401) {
      console.warn('Session expired or unauthorized.');

      // Note: We can't use the logout service here due to circular dependency
      // (this api module is used by auth service), so we manually clear the stores
      
      // clear user data from local storage
      userStore.clearUser();

      // clear the titles
      noteTitleStore.clearNoteTitles();

      // redirect to login
      router.push('/login');

      // reject with a custom error so it can be handled elsewhere if needed
      return Promise.reject(new Error('SessionExpired'));
    }

    return Promise.reject(error);
  }
);

export default api;
