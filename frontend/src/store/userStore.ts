import type { User } from '@/models/user';
import { reactive } from 'vue';

const state = reactive({
  user: null as User | null,
});

export const userStore = {
  setUser(user: User) {
    localStorage.setItem('user', JSON.stringify(user));
    state.user = user;
  },

  getUser(): User {
    if (state.user) {
      return state.user;
    }

    const userData = localStorage.getItem('user');
    if (userData) {
      try {
        const parsedUser = JSON.parse(userData) as User;
        state.user = parsedUser;
        return parsedUser;
      } catch (error) {
        console.error('Failed to parse user from localStorage:', error);
      }
    }
    throw new Error('User not authenticated');
  },

  clearUser() {
    localStorage.removeItem('user');
    state.user = null;
  },
};
