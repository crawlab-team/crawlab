interface LoginForm {
  username: string;
  password: string;
  confirmPassword?: string;
  email?: string;
}

interface LoginRules {
  username: ElFormRule[];
  password: ElFormRule[];
  confirmPassword: ElFormRule[];
}
