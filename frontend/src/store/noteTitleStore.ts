import type { NoteTitle } from '@/models/title';
import { reactive } from 'vue';

const state = reactive({
  noteTitles: [] as NoteTitle[],
});

export const noteTitleStore = {
  addNoteTitle(noteTitle: NoteTitle) {
    const noteTitlesData = localStorage.getItem('noteTitles');
    if (noteTitlesData) {
      try {
        const parsedNoteTitles = JSON.parse(noteTitlesData) as NoteTitle[];
        state.noteTitles = parsedNoteTitles;
      } catch (error) {
        console.error('Failed to parse noteTitles from localStorage:', error);
        state.noteTitles = [];
      }
    }

    state.noteTitles.push(noteTitle);
    localStorage.setItem('noteTitles', JSON.stringify(state.noteTitles));
  },

  getNoteTitles(): NoteTitle[] {
    if (state.noteTitles.length > 0) {
      return state.noteTitles;
    }

    const noteTitlesData = localStorage.getItem('noteTitles');
    if (noteTitlesData) {
      try {
        const parsedNoteTitles = JSON.parse(noteTitlesData) as NoteTitle[];
        state.noteTitles = parsedNoteTitles;
        return parsedNoteTitles;
      } catch (error) {
        console.error('Failed to parse noteTitles from localStorage:', error);
      }
    }
    return [];
  },

  clearNoteTitleById(noteId: number) {
    const noteTitlesData = localStorage.getItem('noteTitles');
    if (noteTitlesData) {
      try {
        const parsedNoteTitles = JSON.parse(noteTitlesData) as NoteTitle[];
        state.noteTitles = parsedNoteTitles.filter(note => note.note_id !== noteId);
        localStorage.setItem('noteTitles', JSON.stringify(state.noteTitles));
      } catch (error) {
        console.error('Failed to parse noteTitles from localStorage:', error);
      }
    }
  },

  clearNoteTitles() {
    localStorage.removeItem('noteTitles');
    state.noteTitles = [];
  },
};
