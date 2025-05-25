<script lang="ts">
export const iframeHeight = '800px'
export const description = 'A sidebar with collapsible submenus.'
</script>

<script setup lang="ts">
import AppSidebar from '@/components/Sidebar.vue'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Separator } from '@/components/ui/separator'
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from '@/components/ui/sidebar'

import { nextTick, ref } from 'vue'
import Button from '@/components/ui/button/Button.vue'
import Toggle from '@/components/Toggle.vue'
import Modal from '@/components/Modal.vue'
import { Input } from '@/components/ui/input'
import { createNote, editNote, deleteNote } from '@/notes/api/notesApi'
import type { Block } from '@/models/block'
import { createBlock } from '@/notes/crypto/createBlock'
import { userStore } from '@/store/userStore'
import { noteTitleStore } from '@/store/noteTitleStore'
import type { NoteTitle } from '@/models/title'
import type { Note } from '@/models/note'
import { fetchAndDecryptNote } from '@/notes/notesService'
import { blockHash } from '@/notes/crypto/blockHash'
import { fromByteArray as toBase64 } from 'base64-js';
import { showConfirm, renderAlert } from '@/store/notifications';


const showPasswordModal = ref(false)
const showDecryptPasswordModal = ref(false)
const password = ref('')
const decryptPassword = ref('')
const passwordError = ref('')
const decryptPasswordError = ref('')
const selectedNoteId = ref<number | null>(null)
const isSaving = ref(false)
const isDecrypting = ref(false)

const noteData = ref({
  id: null as number | null,
  title: '',
  body: '',
  hash: '',
  isIntegrityValid: true, // initialize with true for new notes
})

const isEditingTitle = ref(true)
const noteInputRef = ref<HTMLTextAreaElement | null>(null)
const titleInputRef = ref<HTMLInputElement | null>(null)

const emit = defineEmits(['noteCreated'])

// Initialize the component in "new note" state
function initializeNewNote() {
  noteData.value = {
    id: null,
    title: '',
    body: '',
    hash: toBase64(new Uint8Array(32)), // all-zero hash for the first block (base64-encoded)
    isIntegrityValid: true
  }
  isEditingTitle.value = true
  nextTick(() => {
    titleInputRef.value?.focus()
  })
}

// Initialize on component mount
initializeNewNote()

function saveNote() {
  if (noteData.value.body.trim() === '') {
    renderAlert({ message: 'Please write something before saving.', type: 'error' });
    return
  }
  // Clear password and error when opening modal
  password.value = ''
  passwordError.value = ''
  showPasswordModal.value = true
}

async function handlePasswordSubmit() {
  if (!password.value) {
    passwordError.value = 'Password is required'
    return
  }
  isSaving.value = true
  try {
    // Get user data for encryption settings
    const user = userStore.getUser();
    // Force a tick to ensure loading state is rendered
    await nextTick();
    // create the note title object with a default title if empty
    const noteTitle: NoteTitle = {
        note_id: noteData.value.id || 0, // Use 0 for new notes
        title: noteData.value.title.trim() || 'Untitled Note',
        timestamp: ""
    };      // Update note data with default title if empty
      noteData.value.title = noteData.value.title.trim() || 'Untitled Note';
      
      // Create block record
      const block: Block = await createBlock(
        noteData.value.title,
        noteData.value.body,
        password.value,
        user,
        noteData.value.hash
      );

    /// New notes are created with id 0
    if (!noteData.value.id || noteData.value.id === 0) {
      const [noteId, timestamp] = await createNote(block);

      noteTitle.timestamp = timestamp; // Set the timestamp for the new note title

      noteData.value.id = noteId; // Update noteData with the new ID

      noteTitle.note_id = noteId;

      // Save note title to store
      noteTitleStore.addNoteTitle(noteTitle);      // Clear the form after successful creation
      noteData.value = {
        id: null,
        title: '',
        body: '',
        hash: toBase64(new Uint8Array(32)),
        isIntegrityValid: true // Nova nota é sempre válida
      }
      isEditingTitle.value = true;
    } else {
      // create the NoteBlock object
      const noteBlock = {
        note_id: noteData.value.id,
        block: block,
      };

      noteTitle.note_id = noteData.value.id; // Ensure note_id is set for existing notes
      // Update existing note
      try {
        const timestamp = await editNote(noteBlock);
        noteTitle.timestamp = timestamp; // Set the timestamp for the edited note title

        // update the title store
        noteTitleStore.clearNoteTitleById(noteData.value.id);
        noteTitleStore.addNoteTitle(noteTitle);

        noteData.value.hash = blockHash(block); // Update hash after saving
      } catch (error) {
        console.error('Error editing note:', error);
        // Close the modal and show error in a notification instead
        showPasswordModal.value = false
        password.value = ''
        passwordError.value = ''
        
        // Show specific error message from backend
        const errorMessage = error instanceof Error ? error.message : 'Failed to edit note. Please try again.'
        renderAlert({ message: errorMessage, type: 'error' });
        return; // Exit early to prevent continuing with the success flow
      }
    }

    // Emit event to refresh sidebar
    emit('noteCreated')
      // Reset state
    showPasswordModal.value = false
    password.value = ''
    passwordError.value = ''
  } catch (error) {
    console.error('Error saving note:', error)
    passwordError.value = 'Failed to save note. Please try again.'
  } finally {
    isSaving.value = false
  }
}

function startEditingTitle() {
  isEditingTitle.value = true
  nextTick(() => {
    titleInputRef.value?.focus()
  })
}

function saveTitle() {
  isEditingTitle.value = false
}

function handleSelectNote(id_nota: number) {
  selectedNoteId.value = id_nota
  // Clear password and error when opening modal
  decryptPassword.value = ''
  decryptPasswordError.value = ''
  showDecryptPasswordModal.value = true
}

async function handleDecryptPasswordSubmit() {
  if (!decryptPassword.value) {
    decryptPasswordError.value = 'Password is required'
    return
  }

  if (!selectedNoteId.value) {
    decryptPasswordError.value = 'No note selected'
    return
  }
  isDecrypting.value = true
  try {
    // Force a tick to ensure loading state is rendered
    await nextTick();
    // Fetch and decrypt the note
    const decryptedNote: Note = await fetchAndDecryptNote(decryptPassword.value, selectedNoteId.value)
      // Populate the form with the decrypted note data
    noteData.value = {
      id: decryptedNote.note_id,
      title: decryptedNote.title,
      body: decryptedNote.body,
      hash: decryptedNote.hash,
      isIntegrityValid: decryptedNote.isIntegrityValid ?? true
    };
      // Reset modal state
    showDecryptPasswordModal.value = false
    decryptPassword.value = ''
    decryptPasswordError.value = ''
    selectedNoteId.value = null
    
    // Set editing state
    isEditingTitle.value = false
    
  } catch (error) {
    console.error('Error decrypting note:', error)
    // Close the modal and show error in a notification instead
    showDecryptPasswordModal.value = false
    decryptPassword.value = ''
    decryptPasswordError.value = ''
    selectedNoteId.value = null
    
    // Show error message
    const errorMessage = error instanceof Error ? error.message : 'Failed to decrypt note. Please check your password and try again.'
    renderAlert({ message: errorMessage, type: 'error' })
  } finally {
    isDecrypting.value = false
  }
}

function handleNewNote() {
  initializeNewNote()
}

async function handleDeleteNote() {
  if (!noteData.value.id) {
    renderAlert({ message: 'No note selected to delete.', type: 'error' });
    return
  }

  const confirmDelete = await showConfirm('Are you sure you want to delete this note?');
  if (!confirmDelete) return

  try {
    await deleteNote(noteData.value.id)
    noteTitleStore.clearNoteTitleById(noteData.value.id)    // Reset note data
    noteData.value = {
      id: null,
      title: '',
      body: '',
      hash: '',
      isIntegrityValid: true // Nota nova é sempre válida
    }

    renderAlert({ message: 'Note deleted successfully.', type: 'info' });
  } catch (error) {
    console.error('Error deleting note:', error)
    renderAlert({ message: 'Failed to delete note. Please try again.', type: 'error' });
  }
}
</script>

<template>
  <SidebarProvider>
    <AppSidebar 
      @selectNote="handleSelectNote"
      @newNote="handleNewNote"
    />
    <SidebarInset>
      <header class="flex h-16 shrink-0 items-center gap-2 border-b px-4">
        <SidebarTrigger class="-ml-1" />
        <Separator orientation="vertical" class="mr-2 h-4" />
        <div class="flex flex-1 justify-between items-center">
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem class="hidden md:block">
                <BreadcrumbLink href="#">
                  Notes
                </BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator class="hidden md:block" />
              <BreadcrumbItem>
                <BreadcrumbPage>
                  {{ noteData.title || 'New Note' }}
                </BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
          <Toggle class="justify-end" />
        </div>
      </header>
      <div class="flex flex-1 flex-col gap-4 p-4">
        <h1 
          v-if="!isEditingTitle && noteData.title"
          @click="startEditingTitle"
          class="text-4xl font-bold cursor-pointer"
        >
          {{ noteData.title }}
        </h1>
        <input
          v-if="isEditingTitle || !noteData.title"
          ref="titleInputRef"
          v-model="noteData.title"
          @blur="saveTitle"
          @keyup.enter="saveTitle"
          class="text-4xl font-bold bg-transparent outline-none"
          placeholder="Enter title..."
        />
        
        <div v-if="noteData.id" class="flex items-center justify-between mb-6 pb-3 border-b border-border/40">
          <div class="flex items-center gap-3">
            <div 
              class="w-3 h-3 rounded-full"
              :class="noteData.isIntegrityValid === true ? 'bg-green-500' : 'bg-red-500'"
            ></div>
            <span class="text-sm font-medium text-foreground">
              Integrity Status:
            </span>
            <span 
              class="text-sm font-semibold px-2 py-1 rounded-md"
              :class="noteData.isIntegrityValid === true 
                ? 'text-green-700 bg-green-100 dark:text-green-300 dark:bg-green-900/30' 
                : 'text-red-700 bg-red-100 dark:text-red-300 dark:bg-red-900/30'"
            >
              {{ noteData.isIntegrityValid === true ? 'Valid' : 'Invalid' }}
            </span>
          </div>
        </div>

        <textarea 
          v-model="noteData.body" 
          ref="noteInputRef"
          class="w-full flex-1 resize-none outline-none bg-transparent" 
          placeholder="Start writing..."
        ></textarea>
            <div class="flex justify-end">
          <Button variant="default" size="sm" @click="saveNote">
            Save
          </Button>
          <Button variant="ghost" size="sm" class="ml-2" @click="handleDeleteNote">
            Remove
          </Button>
        </div>
      </div>
    </SidebarInset>

    <Modal v-if="showPasswordModal" @close="showPasswordModal = false">
      <template #title>Enter Password</template>
      <template #description>Please enter your password to save this note</template>
      
      <form @submit.prevent="handlePasswordSubmit" class="space-y-4">
        <div class="space-y-2">          <Input
            v-model="password"
            type="password"
            placeholder="Enter your password"
            :class="{ 'border-destructive': passwordError }"
          />
          <p v-if="passwordError" class="text-sm text-destructive">{{ passwordError }}</p>
        </div>
      </form>      <template #footer>
        <Button variant="outline" @click="showPasswordModal = false" :disabled="isSaving">Cancel</Button>
        <Button type="submit" @click="handlePasswordSubmit" :disabled="isSaving">
          <template v-if="isSaving">
            <span class="inline-block animate-spin mr-2">⌛</span>
            Saving...
          </template>
          <template v-else>Save</template>
        </Button>
      </template>
    </Modal>

    <Modal v-if="showDecryptPasswordModal" @close="showDecryptPasswordModal = false">
      <template #title>Enter Password</template>
      <template #description>Please enter your password to decrypt and load this note</template>
      
      <form @submit.prevent="handleDecryptPasswordSubmit" class="space-y-4">
        <div class="space-y-2">
          <Input
            v-model="decryptPassword"
            type="password"
            placeholder="Enter your password"
            :class="{ 'border-destructive': decryptPasswordError }"
          />
          <p v-if="decryptPasswordError" class="text-sm text-destructive">{{ decryptPasswordError }}</p>
        </div>
      </form>      <template #footer>
        <Button variant="outline" @click="showDecryptPasswordModal = false" :disabled="isDecrypting">Cancel</Button>
        <Button type="submit" @click="handleDecryptPasswordSubmit" :disabled="isDecrypting">
          <template v-if="isDecrypting">
            <span class="inline-block animate-spin mr-2">⌛</span>
            Decrypting...
          </template>
          <template v-else>Decrypt</template>
        </Button>
      </template>
    </Modal>
  </SidebarProvider>
</template>