import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/components/Login.vue'
import Register from '@/components/Register.vue'
import Home from '@/pages/NoteForm.vue'
import Account from '@/pages/Account.vue'


const routes = [
  { path: '/', redirect: '/home' },
  { path: '/login', name: 'Login', component: Login },
  { path: '/register', name: 'Register', component: Register },
  { path: '/home', name: 'Home', component: Home },
  { path: '/account', name: 'Account', component: Account },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guard to check authentication
router.beforeEach((to, _, next) => {
  const publicPages = ['/login', '/register'];
  const authRequired = !publicPages.includes(to.path);
  
  try {
    const user = localStorage.getItem('user');
    const isLoggedIn = !!user;

    if (authRequired && !isLoggedIn) {
      // Not logged in, redirect to login page
      next('/login');
    } else if (isLoggedIn && publicPages.includes(to.path)) {
      // Logged in and trying to access login/register page, redirect to home
      next('/home');
    } else {
      // Otherwise proceed normally
      next();
    }
  } catch (error) {
    // If there's any error (like corrupted localStorage), redirect to login
    localStorage.removeItem('user');
    localStorage.removeItem('token');
    next('/login');
  }
});

export default router
