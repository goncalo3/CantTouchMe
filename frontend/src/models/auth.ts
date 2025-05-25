// payload sent during user registration.
// all salts and the public key are base64-encoded.
export type RegistrationPayload = {
  name: string;
  email: string;
  hmac_type: 'hmac-sha256' | 'hmac-sha512';
  encryption_type: 'aes-128-cbc' | 'aes-128-ctr';
  login_salt: string;
  encryption_salt: string;
  hmac_salt: string;
  public_key: string;
};

// response received from the server when requesting a login challenge.
export type ChallengeResponse = {
  challenge: string;
  login_salt: string;
  expires_at: string;
};

// payload sent when responding to a login challenge.
export type LoginRequestPayload = {
  email: string;
  challenge: string;
  signature: string;
};
