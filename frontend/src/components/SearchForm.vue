<script setup lang="ts">
import { Label } from '@/components/ui/label'
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarInput,
} from '@/components/ui/sidebar'
import { Search } from 'lucide-vue-next'
import { ref, watch } from 'vue'

const searchQuery = ref('')
const emit = defineEmits(['search'])

watch(searchQuery, (newQuery) => {
  emit('search', newQuery)
})
</script>

<template>
  <form @submit.prevent>
    <SidebarGroup class="py-0">
      <SidebarGroupContent class="relative">
        <Label for="search" class="sr-only">
          Search
        </Label>
        <SidebarInput
          id="search"
          v-model="searchQuery"
          placeholder="Search your notes..."
          class="pl-8"
        />
        <Search class="pointer-events-none absolute left-2 top-1/2 size-4 -translate-y-1/2 select-none opacity-50" />
      </SidebarGroupContent>
    </SidebarGroup>
  </form>
</template>