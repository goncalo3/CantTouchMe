import api from '@/lib/api';
import type { User } from '@/models/user';

export interface UpdateUserRequest {
  name?: string;
  email?: string;
}

export const userService = {
  async updateUser(data: UpdateUserRequest): Promise<User> {
    const response = await api.put('/auth/update', data);
    return response.data.user;
  },
}; 