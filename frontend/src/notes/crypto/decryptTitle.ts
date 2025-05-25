import { aes128cbcDecrypt, aes128ctrDecrypt } from './encryption';
import { deriveEncryptionKey } from '../../auth/crypto/keyDerivation';
import { toByteArray as fromBase64 } from 'base64-js';
import type { CipherType } from '@/models/block';
import type { NoteTitle, EncryptedTitle } from '@/models/title';


// decrypts the encrypted title of a note block.
//
// this function only decrypts the 'cipher_title' field to retrieve the plaintext title.
// it does not validate the integrity (MAC) — this is expected to be done later when the full note is decrypted.
export async function decryptBlockTitle(
  eTitle: EncryptedTitle,
  password: string,
  encryptionSalt: string,
  encryptionType: CipherType
): Promise<NoteTitle> {
  // derive encryption key from password + encryption salt
  const encryptionKey = await deriveEncryptionKey(password, encryptionSalt);

  // decode the IV from base64
  const iv = fromBase64(eTitle.iv_title);

  // decrypt the title using the selected AES mode
  const NoteTitle: NoteTitle = {
    note_id: eTitle.note_id,
    title: encryptionType === 'aes-128-cbc'
      ? await aes128cbcDecrypt(eTitle.cipher_title, encryptionKey, iv)
      : await aes128ctrDecrypt(eTitle.cipher_title, encryptionKey, iv),
    timestamp: eTitle.timestamp,
  };

  return NoteTitle    
}
