// represents the encrypted form of a note title, as stored in the backend
export type EncryptedTitle = {
  note_id: number;
  cipher_title: string;
  timestamp: string;
  iv_title: string;
}

// represents the decrypted form of a note title, as shown in the UI
export type NoteTitle = {
    note_id: number;
    title: string;
    timestamp: string;
};