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
import {
   Select,
   SelectContent,
   SelectItem,
   SelectTrigger,
   SelectValue,
} from '@/components/ui/select' 
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import logo from '@/assets/logo.svg'

import { registerUser } from '@/auth/crypto/register'
import type { RegistrationPayload } from '@/models/auth'
import { sendRegistrationData } from '@/auth/api/authApi'
import { renderAlert, showAlertWithRedirect } from '@/store/notifications';

const name = ref('')
const email = ref('')
const password = ref('')
const encryptionType = ref<'aes-128-cbc' | 'aes-128-ctr' | ''>('')
const hmacType = ref<'hmac-sha256' | 'hmac-sha512' | ''>('')

const router = useRouter()

const handleRegister = async (e: Event) => {
  e.preventDefault()

  if (!name.value || !email.value || !password.value || !encryptionType.value || !hmacType.value) {
    renderAlert({ message: 'Please fill in all fields', type: 'error' });
    return
  }

  // check if password longer than 8 characters, upper and lower case letters, numbers and special characters
  const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/
  if (!passwordRegex.test(password.value)) {
    renderAlert({ message: 'Password must be at least 8 characters long and contain upper and lower case letters, numbers, and special characters', type: 'error' });
    return
  }


  try {
    const payload: RegistrationPayload = await registerUser(
      name.value,
      email.value,
      password.value,
      hmacType.value,
      encryptionType.value
    )
    await sendRegistrationData(payload)
    showAlertWithRedirect({ message: 'Account created successfully! Please sign in', type: 'info' });
    router.push('/login')
  } catch (err: any) {
    let message = 'Registration failed';
    
    // Handle specific error cases
    if (err.message.includes('Email already registered')) {
      message = 'This email is already registered';
    } else if (err.message.includes('invalid email format')) {
      message = 'Please enter a valid email address';
    } else if (err.message.includes('salt must be encoded in base64')) {
      message = 'An error occurred during registration. Please try again';
    } else if (err.message.includes('hmac type must be')) {
      message = 'Please select a valid HMAC algorithm';
    } else if (err.message.includes('encryption type must be')) {
      message = 'Please select a valid encryption algorithm';
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
          Welcome!
        </CardTitle>
        <CardDescription>
          Sign up to access <strong>Can't Touch Me!</strong> and write your first secure note.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form @submit="handleRegister" class="grid gap-4">
         <div class="grid gap-2">
            <Label for="name">Name</Label>
            <Input id="name" v-model="name" placeholder="Han Solo" required />
          </div>
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
         
         <div class="grid gap-2">
            <Label for="tipo-cifra">Type of Cipher:</Label>
            <Select v-model="encryptionType">
               <SelectTrigger class="w-full">
                  <SelectValue placeholder="Choose your prefered cypher" />
               </SelectTrigger>
                  <SelectContent>
                     <SelectItem value="aes-128-cbc">AES-128-CBC</SelectItem>
                     <SelectItem value="aes-128-ctr">AES-128-CTR</SelectItem>
                  </SelectContent>
            </Select>
         </div>

         <div class="grid gap-2">
            <Label for="tipo-hmac">Type of HMAC:</Label>
            <Select v-model="hmacType">
               <SelectTrigger class="w-full">
                  <SelectValue placeholder="Choose your prefered HMAC" />
               </SelectTrigger>
                  <SelectContent>
                     <SelectItem value="hmac-sha256">HMAC-SHA256</SelectItem>
                     <SelectItem value="hmac-sha512">HMAC-SHA512</SelectItem>
                  </SelectContent>
            </Select>
         </div>
          <Button type="submit" class="w-full">
            Sign Up
          </Button>
        </form>
        <div class="mt-4 text-center text-sm">
          Already have an account?
          <router-link to="/login" class="underline">
            Sign in
          </router-link>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
