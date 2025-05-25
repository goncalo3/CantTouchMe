// this type mirrors the backend model for a note block.
// each block securely stores a single encrypted note entry.
export type Block = {
  prev_hash: string;           
  iv: string;             
  iv_title: string;        
  cipher_title: string;    
  ciphertext: string;     
  mac: string;
  signature: string;
  timestamp: string;             
};

// supported HMAC hashing algorithms for block integrity
export type HashType = 'hmac-sha256' | 'hmac-sha512';

// supported AES encryption modes for note encryption
export type CipherType = 'aes-128-cbc' | 'aes-128-ctr'; 