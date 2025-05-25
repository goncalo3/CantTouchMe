import { cbc, ctr } from '@noble/ciphers/aes.js';
import { fromByteArray as toBase64, toByteArray as fromBase64 } from 'base64-js';


// AES-CBC requires the plaintext to be a multiple of 16 bytes (block size)
//
// to ensure this, we apply PKCS#7-style padding during encryption:
//   - the number of padding bytes (N) is calculated as (16 - plaintext.length % 16)
//   - N bytes are added, each with the value N
//
// during decryption, we reverse this process:
//   - the last byte of the decrypted data tells us how many bytes to remove
//   - we slice off the padding to recover the original plaintext
//
// this ensures correct block alignment and compatibility with standard AES-CBC implementations.


// Adds PKCS#7 padding to the input to make its length a multiple of 16 bytes
function padPKCS7(data: Uint8Array): Uint8Array {
  const blockSize = 16;
  const padLen = blockSize - (data.length % blockSize);
  const padding = new Uint8Array(padLen).fill(padLen);
  return new Uint8Array([...data, ...padding]);
}

// Removes PKCS#7 padding from the decrypted data
function unpadPKCS7(data: Uint8Array): Uint8Array {
  const padLen = data[data.length - 1];
  return data.slice(0, data.length - padLen);
}

// AES-128-CBC encryption
export async function aes128cbcEncrypt(
  text: string,
  key: Uint8Array,
  iv: Uint8Array
): Promise<string> {
  const cipher = cbc(key, iv);
  const plaintextBytes = new TextEncoder().encode(text);
  const padded = padPKCS7(plaintextBytes);
  const ciphertextBytes = cipher.encrypt(padded);
  return toBase64(ciphertextBytes);
}

// AES-128-CBC decryption
export async function aes128cbcDecrypt(
  base64Ciphertext: string,
  key: Uint8Array,
  iv: Uint8Array
): Promise<string> {
  const cipher = cbc(key, iv);
  const ciphertextBytes = fromBase64(base64Ciphertext);
  const decryptedPadded = cipher.decrypt(ciphertextBytes);
  const decrypted = unpadPKCS7(decryptedPadded);
  return new TextDecoder().decode(decrypted);
}

// AES-128-CTR encryption
export async function aes128ctrEncrypt(
  text: string,
  key: Uint8Array,
  nonce: Uint8Array
): Promise<string> {
  const cipher = ctr(key, nonce);
  const plaintextBytes = new TextEncoder().encode(text);
  const ciphertextBytes = cipher.encrypt(plaintextBytes);
  return toBase64(ciphertextBytes);
}

// AES-128-CTR decryption
export async function aes128ctrDecrypt(
  base64Ciphertext: string,
  key: Uint8Array,
  nonce: Uint8Array
): Promise<string> {
  const cipher = ctr(key, nonce);
  const ciphertextBytes = fromBase64(base64Ciphertext);
  const plaintextBytes = cipher.decrypt(ciphertextBytes);
  return new TextDecoder().decode(plaintextBytes);
}
