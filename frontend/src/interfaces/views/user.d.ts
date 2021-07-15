import {ROLE_ADMIN, ROLE_NORMAL} from '@/constants/user';

declare global {
  interface User {
    _id?: string;
    username?: string;
    password?: string;
    role?: string;
    email?: string;
  }

  type UserRole = ROLE_ADMIN | ROLE_NORMAL;
}
