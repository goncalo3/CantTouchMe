import { getSessionKeys } from '../../auth/crypto/keyDerivation';
import { aes128cbcDecrypt, aes128ctrDecrypt } from './encryption';
import { hmac } from '@noble/hashes/hmac';
import { sha256, sha512 } from '@noble/hashes/sha2';
import { toByteArray as fromBase64 } from 'base64-js';
import type { Block, HashType } from '@/models/block';
import type { User } from '@/models/user';

// validates the integrity of a block using HMAC.
// compares the MAC stored in the block with the MAC freshly computed using the derived key.
function validateMac(block: Block, hmacKey: Uint8Array, hashType: HashType): boolean {
  const encoder = new TextEncoder();

  // prepare the input for HMAC: cipher_title + ciphertext
  const macInput = encoder.encode(block.cipher_title + block.ciphertext);

  // recompute expected MAC using provided key and hash algorithm
  const expectedMac = hmac(
    hashType === 'hmac-sha512' ? sha512 : sha256,
    hmacKey,
    macInput
  );

  const actualMac = fromBase64(block.mac);
  
  // constant-time comparison to prevent timing attacks
  if (expectedMac.length !== actualMac.length) {
    return false;
  }
  
  for (let i = 0; i < expectedMac.length; i++) {
    if (expectedMac[i] !== actualMac[i]) {
      return false;
    }
  }
  
  return true;
}

// decrypts the body of a block after verifying its integrity.
export async function decryptBodyFromBlock(
  block: Block,
  password: string,
  user : User,
): Promise<{ body: string, isIntegrityValid: boolean }> {
  // derive encryption and integrity keys from the password
  const { encryptionKey, hmacKey } = await getSessionKeys(password, user);

  // validate the block's MAC before decrypting
  const macValid = validateMac(block, hmacKey, user.hmac_type);

  // decrypt the body (ciphertext) regardless of MAC validity
  const iv = fromBase64(block.iv);
  const decryptedBody = user.encryption_type === 'aes-128-cbc'
    ? await aes128cbcDecrypt(block.ciphertext, encryptionKey, iv)
    : await aes128ctrDecrypt(block.ciphertext, encryptionKey, iv);
    
  // return both the decrypted content and MAC verification result
  return {
    body: decryptedBody,
    isIntegrityValid: macValid
  };
}