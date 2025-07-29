const API_BASE = process.env.REACT_APP_API_BASE || 'http://localhost:8080/api/v1';

export function getToken() {
  return localStorage.getItem('token') || '';
}

export function setToken(t) {
  localStorage.setItem('token', t);
}

export function apiFetch(path, options = {}) {
  const token = getToken();
  const headers = { 'Content-Type': 'application/json', ...options.headers };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  return fetch(`${API_BASE}${path}`, { ...options, headers });
}
