import { getToken } from './auth';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID || 'tenant_01';

export function api(path: string, init?: RequestInit): Promise<Response> {
  const url = `${API_URL}/api/v1${path}`;
  const headers = new Headers(init?.headers);

  headers.set('Content-Type', 'application/json');
  headers.set('X-Tenant-ID', TENANT_ID);

  const token = getToken();
  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  // Propaga trace context se disponível (W3C Trace Context)
  const traceparent = getTraceparent();
  if (traceparent) {
    headers.set('traceparent', traceparent);
  }

  return fetch(url, {
    ...init,
    headers,
  });
}

function getTraceparent(): string | null {
  // Simplificado: em produção, integrar com OTel propagator
  if (typeof window !== 'undefined') {
    // Poderia ler de um span ativo
    return null;
  }
  return null;
}
