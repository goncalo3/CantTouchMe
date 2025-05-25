import { randomBytes } from '@noble/hashes/utils';
import { hmac } from '@noble/hashes/hmac';
import { sha256, sha512 } from '@noble/hashes/sha2';
import { fromByteArray as toBase64 } from 'base64-js';
import { aes128cbcEncrypt, aes128ctrEncrypt } from './encryption';
import { getSessionKeys, derivePrivateKey } from '../../auth/crypto/keyDerivation';
import type { Block } from '@/models/block';
import type { User } from '@/models/user';
import { signBlock } from './signBlock';

// creates a secure, encrypted, authenticated, and signed record block.
// the block contains a title and body, encrypted and authenticated using keys derived from the user's password.
// it also includes a digital signature using the user's Ed25519 key for non-repudiation and tamper detection.
export async function createBlock(
	title: string,
  body: string,
  password: string,
  user: User,                     
  prevHashBase64: string
): Promise<Block> {
  // derive encryption and HMAC keys from password and salts
  const { encryptionKey, hmacKey } = await getSessionKeys(password, user);

  // generate fresh random IVs for both the title and body encryption
  const ivTitleBytes = randomBytes(16);
  const ivBodyBytes = randomBytes(16);

  // encrypt the title using the selected AES mode
  const cipherTitle = user.encryption_type === 'aes-128-cbc'
    ? await aes128cbcEncrypt(title, encryptionKey, ivTitleBytes)
    : await aes128ctrEncrypt(title, encryptionKey, ivTitleBytes);
  
  // encrypt the body using the selected AES mode
  const ciphertext = user.encryption_type === 'aes-128-cbc'
    ? await aes128cbcEncrypt(body, encryptionKey, ivBodyBytes)
    : await aes128ctrEncrypt(body, encryptionKey, ivBodyBytes);

  // prepare metadata
  const prev = prevHashBase64;
  const timestamp = new Date().toISOString().replace(/\.\d{3}Z$/, 'Z'); // RFC3339 format

  // construct input to HMAC: concatenate encrypted title and body
  const encoder = new TextEncoder();
  const macInput = encoder.encode(cipherTitle + ciphertext);

  // compute HMAC using the selected hash algorithm
  const macBytes = hmac(
    user.hmac_type === 'hmac-sha512' ? sha512 : sha256,
    hmacKey,
    macInput
  );

  // derive Ed25519 private key to sign the block
  const privateKey = await derivePrivateKey(password, user.login_salt);

  // build the block (with empty signature for now)
  const block: Block = {
    prev_hash: prev,
    iv: toBase64(ivBodyBytes),
    iv_title: toBase64(ivTitleBytes),
    cipher_title: cipherTitle,
    ciphertext: ciphertext,
    mac: toBase64(macBytes),
    signature: "",
    timestamp: timestamp,
  };

  // sign the block and return the finalized version
  return signBlock(block, privateKey);
}
