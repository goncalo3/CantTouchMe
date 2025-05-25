<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import logo from '@/assets/logo.svg'
import { loginWithPassword } from '@/auth/services/authService'
import { userStore } from '@/store/userStore'
import { renderAlert, showAlertWithRedirect } from '@/store/notifications';

const router = useRouter()

// check if user is already logged in
try {
   if (userStore.getUser() != null) {
     router.push('/');
   }
} catch (error) {
  // do nothing
} 

const email = ref('')
const password = ref('')

const handleLogin = async (e: Event) => {
  e.preventDefault()
  try {
    await loginWithPassword(email.value, password.value);
    showAlertWithRedirect({ message: 'Welcome back!', type: 'info'});
    router.push('/');
  } catch (err) {
    const error = err as Error;
    let message = 'Invalid email or password';
    
    // Handle specific error cases
    if (error.message.includes('Challenge not found')) {
      message = 'Invalid email or password';
    } else if (error.message.includes('Challenge expired')) {
      message = 'Login session expired. Please try again';
    } else if (error.message.includes('Challenge already used')) {
      message = 'Login attempt expired. Please try again';
    } else if (error.message.includes('Invalid signature')) {
      message = 'Invalid password';
    } else if (error.message.includes('User not found')) {
      message = 'Invalid email or password';
    }
    
    renderAlert({ message, type: 'error' });
  }
}
</script>

<template>
  <div class="mx-auto w-90 flex flex-col items-center gap-4">
    <a href="#" class="flex flex-col items-center gap-2">
      <div        class="flex aspect-square size-8 items-center justify-center">
        <img :src="logo" alt="Logo" class="size-8" />
      </div>
      <span class="font-semibold">Can't Touch Me!</span>
    </a>
    <Card class="mx-auto w-90">
      <CardHeader>
        <CardTitle class="text-2xl">
          Welcome Back!
        </CardTitle>
        <CardDescription>
          Sign in to continue using <strong>Can't Touch Me!</strong> to manage your notes securely.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form @submit="handleLogin" class="grid gap-4">
          <div class="grid gap-2">
            <Label for="email">Email</Label>
            <Input id="email" type="email" v-model="email" placeholder="han@solo.com" required />
          </div>
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="password">Password</Label>
            </div>
            <Input id="password" type="password" v-model="password" required />
          </div>
          <Button type="submit" class="w-full">
            Sign In
          </Button>
        </form>
        <div class="mt-4 text-center text-sm">
          Don't have an account?
          <router-link to="/register" class="underline">
            Sign up
          </router-link>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
