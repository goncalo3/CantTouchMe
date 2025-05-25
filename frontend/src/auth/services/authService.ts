import { requestLoginChallenge, sendLoginSignature, logout as logoutApi } from '@/auth/api/authApi';
import { signLoginChallenge } from '@/auth/crypto/login';
import type { User } from '@/models/user';
import router from '@/router';
import { userStore } from '@/store/userStore';

import { fetchNoteTitles } from '@/notes/api/notesApi';
import { decryptBlockTitle } from '@/notes/crypto/decryptTitle';
import type { CipherType } from '@/models/block';
import type {  NoteTitle, EncryptedTitle } from '@/models/title';
import { noteTitleStore } from '@/store/noteTitleStore';

/**
 * performs the full login flow:
 * 1. requests a login challenge and salt from the server
 * 2. signs the challenge using a key derived from the user's password
 * 3. sends the signed challenge to the server for authentication
 * 4. if successful, stores the user and decrypted note titles in state
 */
export async function loginWithPassword(email: string, password: string) {
    // reset any previously stored user or titles
    userStore.clearUser();
    noteTitleStore.clearNoteTitles();

    // step 1: request challenge and login salt from backend
    const { challenge, login_salt } = await requestLoginChallenge(email);

    // step 2: derive private key and sign the challenge
    const signedPayload = await signLoginChallenge(email, password, challenge, login_salt);

    // step 3: send signed challenge to backend and receive authenticated user data
    const loginResponse: User = await sendLoginSignature(signedPayload);

    // save authenticated user in store
    userStore.setUser(loginResponse);

    // step 4: decrypt and store note titles using the password and encryption type
    fetchAndDecryptTitles(password, loginResponse.encryption_type);

    return;
}

// Logs the user out:
export async function logout() {
    try {
        // call the backend logout endpoint to clear the HTTP-only cookie
        await logoutApi();
    } catch (error) {
        console.error('Failed to logout from backend:', error);
        // continue with client-side cleanup even if backend call fails
    }

    // clear the titles from the store
    noteTitleStore.clearNoteTitles();

    // clear the user from the store
    userStore.clearUser();
}

// attempts to retrieve the user from localStorage, otherwise redirects to login.
export async function getUser() {
    // get the user object from local storage
    const user = localStorage.getItem('user');

    // if the user object is not found, go to the /login page
    if (!user) {
        router.push('/login');
        return;
    }

    // parse the user object and return it
    return JSON.parse(user) as User;
}

// fetches all encrypted note titles from the server
// decrypts each title using the user's password
// saves the decrypted titles to the store
export async function fetchAndDecryptTitles(password: string, encryptionType: CipherType) {
  const user = userStore.getUser();
  const eTitles: EncryptedTitle[] = await fetchNoteTitles();
  
  for (const eTitle of eTitles) {
    try {
      const title : NoteTitle = await decryptBlockTitle(
        eTitle, 
        password, 
        user.encryption_salt,
        encryptionType
      );
      noteTitleStore.addNoteTitle(title);
    } catch (err) {
      console.warn('Failed to decrypt title:', err);
    }
  }
}
