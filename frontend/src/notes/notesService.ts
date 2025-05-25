import { fetchNotes } from './api/notesApi';
import { decryptBodyFromBlock } from './crypto/decryptBody';
import { userStore } from '@/store/userStore';
import { aes128cbcDecrypt, aes128ctrDecrypt } from './crypto/encryption';
import { deriveEncryptionKey } from '../auth/crypto/keyDerivation';
import { toByteArray } from 'base64-js';
import type { Note } from '@/models/note';
import { blockHash } from './crypto/blockHash';

/**
 * Fetches and decrypts a complete note using the user's encryption settings
 * @param password - The user's password for decryption
 * @param noteId - The ID of the note to fetch and decrypt
 * @returns Promise<Note> - The decrypted note with title and body
 */
export async function fetchAndDecryptNote(password: string, noteId: number): Promise<Note> {
  try {
    // Get user encryption settings
    const user = userStore.getUser();

    // Fetch the encrypted note block (now returns single block)
    const block = await fetchNotes(noteId);

    // Decrypt the title from the block
    const encryptionKey = await deriveEncryptionKey(password, user.encryption_salt);
    const titleIv = toByteArray(block.iv_title);

    const decryptedTitle = user.encryption_type === 'aes-128-cbc'
      ? await aes128cbcDecrypt(block.cipher_title, encryptionKey, titleIv)
      : await aes128ctrDecrypt(block.cipher_title, encryptionKey, titleIv);

    // Decrypt the body from the block and get integrity status
    const { body: decryptedBody, isIntegrityValid } = await decryptBodyFromBlock(
      block,
      password,
      user,
    );

    // Get timestamp from the block
    const timestamp = block.timestamp || new Date().toISOString();
    
    // Return the decrypted note with integrity status
    const decryptedNote: Note = {
      note_id: noteId,
      title: decryptedTitle.trim() || 'Untitled Note',
      body: decryptedBody,
      timestamp: timestamp,
      hash: blockHash(block),
      isIntegrityValid,
    };

    return decryptedNote;

  } catch (error) {
    console.error('Failed to fetch and decrypt note:', error);
    throw new Error(`Failed to decrypt note: ${error instanceof Error ? error.message : 'Unknown error'}`);
  }
}