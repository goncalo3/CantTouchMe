import type { Block } from './block';

// represents a single block belonging to a note.
export type NoteBlock = {
  note_id: number | null;
  block: Block;
};

// represents a full chain of blocks for a given note.
export type NoteBlockChain = {
    note_id: string;
    blocks: Block[];
};

// represents a fully decrypted note, used in the frontend to render UI.
export type Note = {
    note_id: number;
    title: string;
    body: string;
    timestamp: string;
    hash: string;
    isIntegrityValid?: boolean; 
};