import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { apiFetch, setToken } from '../utils/api';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  const submit = (e) => {
    e.preventDefault();
    setMessage('در حال ورود...');
    apiFetch('/login', { method: 'POST', body: JSON.stringify({ email, password }) })
      .then(res => res.ok ? res.json() : Promise.reject())
      .then(data => {
        setToken(data.token);
        navigate('/dashboard');
      })
      .catch(() => setMessage('خطا در ورود'));
  };

  return (
    <div className="container d-flex align-items-center justify-content-center vh-100">
      <form onSubmit={submit} className="card p-4 shadow" style={{minWidth: '320px'}}>
        <h4 className="mb-3 text-center">ورود</h4>
        <input className="form-control mb-2" placeholder="ایمیل" value={email} onChange={e => setEmail(e.target.value)} />
        <input type="password" className="form-control mb-2" placeholder="گذرواژه" value={password} onChange={e => setPassword(e.target.value)} />
        {message && <div className="small text-center text-muted mb-2">{message}</div>}
        <button className="btn btn-primary w-100" type="submit">ورود</button>
      </form>
    </div>
  );
}
