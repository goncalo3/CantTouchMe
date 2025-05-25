// import the configured Axios instance with base URL from .env
import api from '@/lib/api';
import type { User } from '@/models/user';
import type {
  RegistrationPayload,
  ChallengeResponse,
  LoginRequestPayload,
} from '@/models/auth';

// sends a registration request to the backend
export async function sendRegistrationData(payload: RegistrationPayload): Promise<any> {
  try {
    const res = await api.post('/auth/register', payload);
    return res.data;
  } catch (error: any) {
    const errorMessage = error.response?.data || error.message || 'Registration failed';
    throw new Error(errorMessage);
  }
}

// sends the user's email to the backend to request a login challenge and login salt
export async function requestLoginChallenge(email: string): Promise<ChallengeResponse> {
  try {
    const res = await api.post('/auth/challenge', { email });
    return res.data as ChallengeResponse;
  } catch (error: any) {
    const errorMessage = error.response?.data || error.message || 'Failed to get challenge from server';
    throw new Error(errorMessage);
  }
}

// sends the signed challenge to the backend for login
export async function sendLoginSignature(payload: LoginRequestPayload): Promise<User> {
  try {
    const res = await api.post('/auth/login', payload);
    return res.data.user as User;
  } catch (error: any) {
    const errorMessage = error.response?.data || error.message || 'Login failed';
    throw new Error(errorMessage);
  }
}

// sends a logout request to the backend to clear the HTTP-only cookie
export async function logout(): Promise<void> {
  try {
    await api.post('/auth/logout');
  } catch (error: any) {
    const errorMessage = error.response?.data || error.message || 'Logout failed';
    throw new Error(errorMessage);
  }
}

// sends a delete request to the backend to delete a user by ID
export async function deleteUser(): Promise<void> {
  try {
    await api.delete('/auth/delete');
  } catch (error: any) {
    const errorMessage = error.response?.data || error.message || 'Failed to delete user';
    throw new Error(errorMessage);
  }
}
