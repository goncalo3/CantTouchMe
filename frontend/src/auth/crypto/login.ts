import * as ed from '@noble/ed25519';
import { toByteArray as fromBase64 } from 'base64-js';
import { fromByteArray as toBase64 } from 'base64-js';
import { derivePrivateKey } from './keyDerivation';

// signs the challenge using the Ed25519 private key derived from password + salt
export async function signLoginChallenge(
  email: string,
  password: string,
  challengeBase64: string,
  loginSaltBase64: string
): Promise<{
  email: string;
  challenge: string;
  signature: string;
}> {
  // derive the Ed25519 private key deterministically from the user's password and login salt
  const privateKey = await derivePrivateKey(password, loginSaltBase64);

  // decode the base64-encoded challenge into a byte array
  const challengeBytes = fromBase64(challengeBase64);
  
  // sign the challenge using the derived Ed25519 private key
  const signature = await ed.signAsync(challengeBytes, privateKey);

  // return the required fields for login
  return {
    email,
    challenge: challengeBase64,
    signature: toBase64(signature),
  };
}
