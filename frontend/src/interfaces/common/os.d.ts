import {
  OS_WINDOWS,
  OS_MAC,
  OS_LINUX,
} from '@/constants/os';

declare global {
  type OS = OS_WINDOWS | OS_MAC | OS_LINUX;
}
