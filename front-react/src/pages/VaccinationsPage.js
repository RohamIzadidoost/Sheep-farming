import React, { useEffect, useState } from 'react';
import NavBar from '../components/NavBar';
import { apiFetch } from '../utils/api';
import { toJalali, toGregorian } from '../utils/jdate';

export default function VaccinationsPage() {
  const [sheepList, setSheepList] = useState([]);
  const [vaccList, setVaccList] = useState([]);
  const [form, setForm] = useState({ sheep:'', vaccine:'', date:'', vaccinator:'', desc:'' });

  const load = () => {
    apiFetch('/sheep').then(res => res.json()).then(list => setSheepList(list));
    apiFetch('/vaccines').then(res => res.json()).then(setVaccList);
  };
  useEffect(() => { load(); }, []);

  const submit = (e) => {
    e.preventDefault();
    apiFetch(`/sheep/${form.sheep}`, { method:'GET' })
      .then(res => res.json())
      .then(sheep => {
        const vaccinations = sheep.vaccinations || [];
        vaccinations.push({ vaccine: form.vaccine, vaccinator: form.vaccinator, description: form.desc, date: toGregorian(form.date) });
        return apiFetch(`/sheep/${form.sheep}`, { method:'PUT', body: JSON.stringify({ vaccinations })});
      })
      .then(() => { setForm({ sheep:'', vaccine:'', date:'', vaccinator:'', desc:'' }); load(); });
  };

  return (
    <>
      <NavBar />
      <main className="container pt-5 mt-4">
        <h2 className="mb-3">سوابق واکسن</h2>
        <table className="table table-bordered table-hover" id="vaccineTable">
          <thead className="table-light"><tr><th>گوسفند</th><th>واکسن</th><th>تاریخ</th><th>واکسیناتور</th><th>توضیح</th></tr></thead>
          <tbody>
            {sheepList.map(s => (s.vaccinations||[]).map((v,i) => (
              <tr key={s.id + '-' + i}><td>{s.earNumber1}</td><td>{v.vaccine}</td><td>{toJalali(v.date)}</td><td>{v.vaccinator}</td><td>{v.description}</td></tr>
            )))}
          </tbody>
        </table>
        <form onSubmit={submit} className="card p-3">
          <div className="row g-2">
            <div className="col"><select className="form-select" value={form.sheep} onChange={e=>setForm({...form,sheep:e.target.value})}>{sheepList.map(s => <option key={s.id} value={s.id}>{s.earNumber1}</option>)}</select></div>
            <div className="col"><select className="form-select" value={form.vaccine} onChange={e=>setForm({...form,vaccine:e.target.value})}>{vaccList.map(v => <option key={v.id} value={v.id}>{v.name}</option>)}</select></div>
            <div className="col"><input type="date" className="form-control" value={form.date} onChange={e=>setForm({...form,date:e.target.value})} /></div>
            <div className="col"><input className="form-control" placeholder="واکسیناتور" value={form.vaccinator} onChange={e=>setForm({...form,vaccinator:e.target.value})}/></div>
            <div className="col"><input className="form-control" placeholder="توضیح" value={form.desc} onChange={e=>setForm({...form,desc:e.target.value})}/></div>
            <div className="col-auto"><button className="btn btn-success" type="submit">افزودن</button></div>
          </div>
        </form>
      </main>
    </>
  );
}
