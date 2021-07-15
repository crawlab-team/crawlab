declare module '*.js';
declare module '*.jpg';
declare module '*.png';
declare module '*.svg';

declare global {
  type ElFormValidator = (rule: any, value: any, callback: any) => void;

  interface ElFormRule {
    required: boolean;
    trigger: string;
    validator: ElFormValidator;
  }
}
