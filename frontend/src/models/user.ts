import type { HashType, CipherType } from './block';

// defines the structure of a User object used throughout the application
export interface User {
    id: number;
    name: string;
    email: string;
    public_key: string;
    encryption_salt: string;
    hmac_salt: string;
    hmac_type: HashType;
    login_salt: string;
    encryption_type: CipherType;
}
