import type { Block } from '@/models/block';
import { sha256 } from 'js-sha256';
import { fromByteArray as toBase64 } from 'base64-js';

// computes a base64-encoded SHA-256 hash of a block's contents.
//
// this can be used to uniquely identify a block and ensure integrity.
export function blockHash(block: Block): string {
  // create a string representation of the block
  const blockString = JSON.stringify({
    prev_hash: block.prev_hash,
    iv: block.iv,
    iv_title: block.iv_title,
    cipher_title: block.cipher_title,
    ciphertext: block.ciphertext,
    mac: block.mac,
    signature: block.signature,
    timestamp: block.timestamp
  });

  // compute the SHA-256 hash of the block string
  const hashBuffer = new Uint8Array(sha256.arrayBuffer(blockString));

  // convert the hash to a Base64 string using base64-js
  return toBase64(hashBuffer);
}