export interface AuthUser {
  id: string;
  email: string;
  name: string;
  role: 'SELLER' | 'BUYER' | 'SYSTEM_ADMIN' | 'GLOBAL_ADMIN';
  active: boolean;
}

const TOKEN_KEY = '2pshop_token';
const USER_KEY = '2pshop_user';

export function getToken(): string | null {
  if (typeof window === 'undefined') return null;
  return window.localStorage.getItem(TOKEN_KEY);
}

export function getStoredUser(): AuthUser | null {
  if (typeof window === 'undefined') return null;
  const raw = window.localStorage.getItem(USER_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as AuthUser;
  } catch {
    return null;
  }
}

export function storeSession(token: string, user: AuthUser) {
  if (typeof window === 'undefined') return;
  window.localStorage.setItem(TOKEN_KEY, token);
  window.localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearSession() {
  if (typeof window === 'undefined') return;
  window.localStorage.removeItem(TOKEN_KEY);
  window.localStorage.removeItem(USER_KEY);
}

export const ROLE_LABELS: Record<AuthUser['role'], string> = {
  SELLER: 'Vendedor',
  BUYER: 'Cliente',
  SYSTEM_ADMIN: 'Administrador do Sistema',
  GLOBAL_ADMIN: 'Administrador Global',
};

// Where each role lands as their "home" after login.
export function roleHomePath(role: AuthUser['role']): string {
  switch (role) {
    case 'SELLER':
      return '/seller';
    case 'SYSTEM_ADMIN':
      return '/admin/system';
    case 'GLOBAL_ADMIN':
      return '/admin/global';
    default:
      return '/products';
  }
}
