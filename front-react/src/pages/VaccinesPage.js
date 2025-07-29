import React, { useEffect, useState } from 'react';
import NavBar from '../components/NavBar';
import { apiFetch } from '../utils/api';

export default function VaccinesPage() {
  const [list, setList] = useState([]);
  const [form, setForm] = useState({ name: '', interval: '' });

  const load = () => {
    apiFetch('/vaccines')
      .then(res => res.json())
      .then(setList);
  };

  useEffect(() => { load(); }, []);

  const submit = (e) => {
    e.preventDefault();
    apiFetch('/vaccines', { method: 'POST', body: JSON.stringify({ name: form.name, intervalMonths: parseInt(form.interval,10) }) })
      .then(() => { setForm({ name:'', interval:'' }); load(); });
  };

  return (
    <>
      <NavBar />
      <main className="container pt-5 mt-4">
        <h2 className="mb-3">لیست واکسن‌ها</h2>
        <table className="table table-bordered table-hover" id="vaccineDefTable">
          <thead className="table-light"><tr><th>نام</th><th>دوره (ماه)</th></tr></thead>
          <tbody>{list.map(v => <tr key={v.id}><td>{v.name}</td><td>{v.intervalMonths}</td></tr>)}</tbody>
        </table>
        <form onSubmit={submit} className="card p-3">
          <div className="row g-2">
            <div className="col"><input className="form-control" placeholder="نام" value={form.name} onChange={e=>setForm({...form,name:e.target.value})} /></div>
            <div className="col"><input type="number" className="form-control" placeholder="دوره" value={form.interval} onChange={e=>setForm({...form,interval:e.target.value})} /></div>
            <div className="col-auto"><button className="btn btn-success" type="submit">افزودن</button></div>
          </div>
        </form>
      </main>
    </>
  );
}
