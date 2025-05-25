import { pbkdf2Async } from '@noble/hashes/pbkdf2';
import { sha256, sha512 } from '@noble/hashes/sha2';
import { toByteArray as fromBase64 } from 'base64-js';
import type { User } from '@/models/user';

// derives a 32-byte Ed25519 private key from a password and base64-encoded salt
export async function derivePrivateKey(password: string, saltBase64: string): Promise<Uint8Array> {
  const salt = fromBase64(saltBase64);
  const passwordBytes = new TextEncoder().encode(password);
  return await pbkdf2Async(sha256, passwordBytes, salt, {c: 100_000, dkLen: 32,});
}

// derives a 128-bit (16-byte) AES encryption key from password and base64-encoded salt
export async function deriveEncryptionKey(
  password: string,
  encryptionSaltBase64: string
): Promise<Uint8Array> {
  const salt = fromBase64(encryptionSaltBase64);
  const passwordBytes = new TextEncoder().encode(password);
  return await pbkdf2Async(sha256, passwordBytes, salt, { c: 100_000, dkLen: 16 }); 
}

// derives an HMAC key (32 or 64 bytes depending on hash function) from password and salt
export async function deriveHMACKey(
  password: string,
  hmacSaltBase64: string,
  hashType: 'hmac-sha256' | 'hmac-sha512'
): Promise<Uint8Array> {
  const salt = fromBase64(hmacSaltBase64);
  const passwordBytes = new TextEncoder().encode(password);
  const hashFn = hashType === 'hmac-sha512' ? sha512 : sha256;
  const dkLen = hashType === 'hmac-sha512' ? 64 : 32;

  return await pbkdf2Async(hashFn, passwordBytes, salt, { c: 100_000, dkLen });
}

// derives both encryption and HMAC keys from the user's password and stored salts
export async function getSessionKeys(
  password: string,
  user: User
): Promise<{
  encryptionKey: Uint8Array;
  hmacKey: Uint8Array;
}> {
  const encryptionKey = await deriveEncryptionKey(password, user.encryption_salt);
  const hmacKey = await deriveHMACKey(password, user.hmac_salt, user.hmac_type);
  return { encryptionKey, hmacKey };
}
