<script setup lang="ts">
import { ref, computed, type ComputedRef } from 'vue'
import { userStore } from '@/store/userStore' 
import { noteTitleStore } from '@/store/noteTitleStore'
import SearchForm from '@/components/SearchForm.vue'
import NavUser from '@/components/NavUser.vue'
import Button from '@/components/ui/button/Button.vue'
import { Collapsible } from '@/components/ui/collapsible'
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarHeader, SidebarMenu, SidebarMenuButton, SidebarMenuItem, type SidebarProps, SidebarRail } from '@/components/ui/sidebar'
import { NotebookPen } from 'lucide-vue-next'
import logo from '@/assets/logo.svg'
import type { NoteTitle } from '@/models/title'

const props = defineProps<SidebarProps>()

const user = userStore.getUser()

const searchQuery = ref('')

// Get notes from noteTitleStore instead of API
const notesTitles: ComputedRef<NoteTitle[]> = computed(() => {
  return noteTitleStore.getNoteTitles();
})

// Filter notes based on search query
const filteredNotes = computed(() => {
  const query = searchQuery.value.toLowerCase()
  if (!query) return notesTitles.value
  
  return notesTitles.value.filter(note => 
    note.title.toLowerCase().includes(query)
  )
})

const selectedNoteId = ref<number | null>(null)
const emit = defineEmits(['selectNote', 'newNote', 'noteCreated'])

function handleSearchQuery(query: string) {
  searchQuery.value = query
}

function selectNote(noteId: number) {
  selectedNoteId.value = noteId
  emit('selectNote', noteId)
}

function createNewNote() {
  emit('newNote')
}
</script>


<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <a href="#">              <div class="flex aspect-square size-8 items-center justify-center">
                <img :src="logo" alt="Logo" class="size-8" />
              </div>
              <div class="flex flex-col gap-0.5 leading-none">
                <span class="font-semibold">Can't Touch Me!</span>
              </div>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
      <SearchForm @search="handleSearchQuery" />
      <div class="border-t border-border my-3"></div>
      <Button variant="default" size="sm" class="w-full" @click="createNewNote">
        <NotebookPen class="size-4" />
        New Note
      </Button>
    </SidebarHeader>    <SidebarContent>
      <SidebarGroup>
        <SidebarMenu>
          <div v-if="filteredNotes.length === 0" class="px-2 py-4 text-sm text-muted-foreground text-center">
            No notes found
          </div>
          <Collapsible v-else v-for="note in filteredNotes" :key="note.note_id" class="group/collapsible w-full">
            <SidebarMenuItem>
              <SidebarMenuButton 
                @click="selectNote(note.note_id)"
                :class="{
                  'bg-accent text-accent-foreground': selectedNoteId === note.note_id,
                  'hover:bg-muted': selectedNoteId !== note.note_id
                }"
                class="transition-colors duration-200"
              >
                <div class="w-full overflow-hidden text-ellipsis whitespace-nowrap">
                  {{ note.title }}
                </div>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </Collapsible>
        </SidebarMenu>
      </SidebarGroup>
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="user" />
    </SidebarFooter>    <SidebarRail />
  </Sidebar>
</template>