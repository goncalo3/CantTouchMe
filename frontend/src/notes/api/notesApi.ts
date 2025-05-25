import api from '@/lib/api';
import type { NoteBlock } from '@/models/note';
import type { Block } from '@/models/block';
import type { EncryptedTitle } from '@/models/title';

// creates a new note
// note: assumes the user is authenticated and token is set as httpOnly cookie
export async function createNote(noteBlock: Block): Promise<[number, string]> {
  try {
    const res = await api.post('/notes/new', noteBlock);
    // return the note id and timestamp as a tuple
    return [res.data.note_id, res.data.timestamp];
  } catch (error: any) {
    const errorMessage = error.response?.data || error.message || 'Failed to create record';
    throw new Error(errorMessage);
  }
}


// sends a new encrypted record block to the backend
// note: assumes the user is authenticated and token is set as httpOnly cookie
export async function editNote(noteBlock: NoteBlock): Promise<string> {
  try {
    const res = await api.post('/notes/edit', noteBlock);
    return res.data.timestamp; // Return the timestamp from the response
  } catch (error: any) {
    const errorMessage = error.response?.data || error.message || 'Failed to save record';
    throw new Error(errorMessage);
  }
}

// fetches the latest encrypted record block for the currently authenticated user
// note: assumes token is send as httpOnly cookie
export async function fetchNotes(noteId: number): Promise<Block> {
  try {
    const res = await api.post('/notes/get', { note_id: noteId });
    return res.data as Block;
  } catch (error: any) {
    // Extract the error message from the backend response
    const errorMessage = error.response?.data || error.message || 'Failed to fetch note';
    throw new Error(errorMessage);
  }
}

// fetches all note titles for the currently authenticated user
// note: assumes token is sent as an httpOnly cookie
export async function fetchNoteTitles(): Promise<EncryptedTitle[]> {
  try {
    const res = await api.get('/notes/titles');
    return res.data as EncryptedTitle[];
  } catch (error: any) {
    const errorMessage = error.response?.data || error.message || 'Failed to fetch note titles';
    throw new Error(errorMessage);
  }
}

// deletes a note by id
// note: assumes the user is authenticated and token is set as httpOnly cookie
export async function deleteNote(noteId: number): Promise<void> {
  try {
    await api.delete(`/notes/delete`, { data: { note_id: noteId } });
  } catch (error: any) {
    const errorMessage = error.response?.data || error.message || 'Failed to delete note';
    throw new Error(errorMessage);
  }
}