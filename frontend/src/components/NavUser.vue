<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

import {
  Avatar,
  AvatarFallback,
} from '@/components/ui/avatar'

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'
import {
  BadgeCheck,
  ChevronsUpDown,
  LogOut,
} from 'lucide-vue-next'
import { userStore } from '@/store/userStore'
import type { User } from '@/models/user'
import { logout } from '@/auth/services/authService'

const handleLogout = async () => {
  try {
    // Call the logout service which handles both backend and frontend cleanup
    await logout();
    
    // Redirect to login page
    router.push('/login');
  } catch (error) {
    console.error('Logout error:', error);
    // Even if logout fails, redirect to login for security
    router.push('/login');
  }
}

const handleAccount = () => {
  router.push('/account')  // Redirects to account page
}

let user: User | null = null;

try {
  user = userStore.getUser()
} catch (error) {
  router.push("/login")
}

const getInitials = (name: string) => {
  return name
    .split(' ')
    .map(part => part[0]?.toUpperCase() || '')
    .slice(0, 2)
    .join('')
}

const { isMobile } = useSidebar()
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <Avatar class="h-8 w-8 rounded-lg">
              <!-- <AvatarImage :src="user.avatar" :alt="user.name" /> -->
              <!-- TODO: Implement avatar logic later -->
              <AvatarFallback class="rounded-lg">
                {{ getInitials(user?.name ?? "NA") }}
              </AvatarFallback>
            </Avatar>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-semibold">{{ user?.name ?? "Username"}}</span>
              <span class="truncate text-xs">{{ user?.email ?? "email@example.com"}}</span>
            </div>
            <ChevronsUpDown class="ml-auto size-4" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
          :side="isMobile ? 'bottom' : 'right'"
          align="end"
          :side-offset="4"
        >
          <DropdownMenuLabel class="p-0 font-normal">
            <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
              <Avatar class="h-8 w-8 rounded-lg">
                <!-- <AvatarImage :src="user.avatar" :alt="user.name" /> -->
                <!-- TODO: Implement avatar logic later -->
                <AvatarFallback class="rounded-lg">
                  {{ getInitials(user?.name ?? "NA") }}
                </AvatarFallback>
              </Avatar>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">{{ user?.name ?? "username"}}</span>
                <span class="truncate text-xs">{{ user?.email ?? "email@example.com"}}</span>
              </div>
            </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
                <DropdownMenuItem @click="handleAccount">
                <BadgeCheck />
                Account
                </DropdownMenuItem>
            </DropdownMenuGroup>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="handleLogout">
                <LogOut />
                Log out
            </DropdownMenuItem>
            </DropdownMenuContent>
                </DropdownMenu>
                </SidebarMenuItem>
            </SidebarMenu>
</template>
