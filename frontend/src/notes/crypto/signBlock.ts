import type { Block } from '@/models/block';
import * as ed from '@noble/ed25519';
import { fromByteArray as toBase64 } from 'base64-js';

// signs a record block using Ed25519 and returns the signed block.
//
// this ensures the block's integrity and authenticity by binding its contents
// to a digital signature that can later be verified with the user's public key.
//
// it also guarantees non-repudiation: only the user with the correct password
// (from which the private key is derived) can produce a valid signature,
// making it cryptographically infeasible to deny authorship of the block.
export async function signBlock(block: Block, privateKey: Uint8Array): Promise<Block> {
  const encoder = new TextEncoder();

  // prepare the data to sign by concatenating critical fields
  const dataToSign = encoder.encode(
    block.prev_hash +
    block.iv +
    block.iv_title +
    block.cipher_title +
    block.ciphertext +
    block.mac +
    block.timestamp
  );

  // generate Ed25519 signature using the user's private key
  const signature = await ed.signAsync(dataToSign, privateKey);

  // store the base64-encoded signature in the block
  block.signature = toBase64(signature);

  return block;
}