<script setup lang="ts">
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
import { userStore } from '@/store/userStore'
import { userService } from '@/services/userService'
import type { User } from '@/models/user'
import { deleteUser } from '@/auth/api/authApi';
import { useRouter } from 'vue-router';
import { noteTitleStore } from '@/store/noteTitleStore';
import { showConfirm, showAlertWithRedirect, renderAlert } from '@/store/notifications';

const router = useRouter();
import { ref } from 'vue'

const user = ref<User>(userStore.getUser())
const isLoading = ref(false)
const successMessage = ref('')
const errorMessage = ref('')

async function saveAccount() {
  try {
    isLoading.value = true
    errorMessage.value = ''
    successMessage.value = ''

    const updatedUser = await userService.updateUser({
      name: user.value.name,
      email: user.value.email,
    })
    
    userStore.setUser(updatedUser)
    successMessage.value = 'Account updated successfully'
  } catch (error: any) {
    let message = 'Failed to update account';
    
    if (error.response?.data?.error) {
      if (error.response.data.error.includes('Email already in use')) {
        message = 'This email is already in use';
      } else if (error.response.data.error.includes('invalid email format')) {
        message = 'Please enter a valid email address';
      }
    }
    
    errorMessage.value = message;
  } finally {
    isLoading.value = false
  }
}

function confirmAndDeleteAccount() {
  showConfirm('Are you sure you want to delete your account? This action is irreversible and all your notes and data will be permanently erased.')
    .then((confirmed) => {
      if (confirmed) {
        deleteUser()
          .then(() => {
            showAlertWithRedirect({ message: 'Account deleted successfully', type: 'info' });
            userStore.clearUser();
            noteTitleStore.clearNoteTitles();
            router.push('/login');
          })
          .catch((error) => {
            let message = 'Failed to delete account';
            
            if (error.message.includes('Unauthorized')) {
              message = 'Your session has expired. Please log in again';
            }
            
            renderAlert({ message, type: 'error' });
          });
      }
    });
}

function goBack() {
  router.push('/');
}
</script>

<template>
  <div class="mx-auto w-90 flex flex-col items-center gap-4">
    <Button variant="ghost" class="absolute top-4 left-4" @click="goBack">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="15 18 9 12 15 6"></polyline>
      </svg>
    </Button>
    <a href="#" class="flex flex-col items-center gap-2">      <div class="flex aspect-square size-8 items-center justify-center">
        <img :src="logo" alt="Logo" class="size-8" />
      </div>
      <span class="font-semibold">Can't Touch Me!</span>
    </a>

    <Card class="mx-auto w-90">
      <CardHeader>
        <CardTitle class="text-2xl">
          Account Settings
        </CardTitle>
        <CardDescription>
          Manage your personal information and security preferences.
        </CardDescription>
      </CardHeader>

      <CardContent>
        <div v-if="successMessage" class="mb-4 p-3 bg-green-100 text-green-700 rounded">
          {{ successMessage }}
        </div>
        <div v-if="errorMessage" class="mb-4 p-3 bg-red-100 text-red-700 rounded">
          {{ errorMessage }}
        </div>

        <div class="grid gap-4">
          <div class="grid gap-2">
            <Label for="name">Name</Label>
            <Input id="name" v-model="user.name" required />
          </div>

          <div class="grid gap-2">
            <Label for="email">Email</Label>
            <Input id="email" v-model="user.email" type="email" required />
          </div>

          <div class="grid gap-2">
            <Label>Password</Label>
            <div class="flex items-center justify-between">
              <span>••••••••</span>
            </div>
          </div>

          <!-- Configurações de Segurança -->
          <div class="grid gap-2">
            <Label>Encryption Algorithm</Label>
            <span class="text-sm">{{ user.encryption_type }}</span>
          </div>

          <div class="grid gap-2">
            <Label>HMAC Algorithm</Label>
            <span class="text-sm">{{ user.hmac_type }}</span>
          </div>

          <Button class="w-full" @click="saveAccount" :disabled="isLoading">
            {{ isLoading ? 'Saving...' : 'Save Changes' }}
          </Button>
            <Button variant="destructive" class="w-full" @click="confirmAndDeleteAccount">
                Delete Account
            </Button>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
