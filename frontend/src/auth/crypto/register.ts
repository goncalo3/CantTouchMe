import * as ed from '@noble/ed25519';
import { randomBytes } from '@noble/hashes/utils';
import { fromByteArray as toBase64 } from 'base64-js';
import { derivePrivateKey } from './keyDerivation';
import type { RegistrationPayload } from '@/models/auth';

// generate a random 32-byte salt
function generateSalt(): Uint8Array {
    return randomBytes(32);
}

// generates a deterministic Ed25519 key pair and returns the full registration payload
export async function registerUser(
  name: string,
  email: string,
  password: string,
  hmacType: "hmac-sha256" | "hmac-sha512",
  encryptionType: "aes-128-cbc" | "aes-128-ctr"
): Promise<RegistrationPayload> {
  // generate all salts
  const loginSalt = generateSalt();
  const encryptionSalt = generateSalt();
  const hmacSalt = generateSalt();

  // derive the Ed25519 private key from password + loginSalt (PBKDF2)
  const privateKey = await derivePrivateKey(password, toBase64(loginSalt));

  // derive the Ed25519 public key from the private key
  const publicKey = await ed.getPublicKeyAsync(privateKey);

  // return a registration payload with all required values encoded in base64
  return {
    name,
    email,
    hmac_type: hmacType,
    encryption_type: encryptionType,
    login_salt: toBase64(loginSalt),
    encryption_salt: toBase64(encryptionSalt),
    hmac_salt: toBase64(hmacSalt),
    public_key: toBase64(publicKey),
  } as RegistrationPayload;
}
