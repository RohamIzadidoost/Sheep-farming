import React, { useEffect, useState } from 'react';
import NavBar from '../components/NavBar';
import JalaliDatePicker from '../components/JalaliDatePicker';
import { apiFetch } from '../utils/api';
import { toJalali } from '../utils/jdate';

export default function SheepPage() {
  const [list, setList] = useState([]);
  const [filters, setFilters] = useState({ gender: '', minAge: '', maxAge: '' });
  const [form, setForm] = useState({ gender: 'male', dob: '', weight: '', ear1: '', ear2: '', ear3: '', gen: '' });
  const [editingId, setEditingId] = useState(null);

  useEffect(() => {
    load();
  }, []);

  const load = () => {
    const params = new URLSearchParams();
    if(filters.gender) params.append('gender', filters.gender);
    if(filters.minAge) params.append('minAgeDays', parseInt(filters.minAge) * 30);
    if(filters.maxAge) params.append('maxAgeDays', parseInt(filters.maxAge) * 30);
    const qs = params.toString() ? `?${params.toString()}` : '';
    apiFetch(`/sheep${qs}`)
      .then(res => res.json())
      .then(setList);
  };

  const submit = (e) => {
    e.preventDefault();
    const body = JSON.stringify({
      gender: form.gender,
      dateOfBirth: form.dob,
      birthWeight: parseFloat(form.weight) || 0,
      earNumber1: form.ear1,
      earNumber2: form.ear2,
      earNumber3: form.ear3,
      fatherGen: form.gen
    });
    const method = editingId ? 'PUT' : 'POST';
    const url = editingId ? `/sheep/${editingId}` : '/sheep';
    apiFetch(url, { method, body }).then(() => { setEditingId(null); setForm({ gender:'male', dob:'', weight:'', ear1:'', ear2:'', ear3:'', gen:'' }); load(); });
  };

  return (
    <>
      <NavBar />
      <main className="container pt-5 mt-4">
        <div className="d-flex justify-content-between align-items-center mb-3">
          <h2 className="m-0">لیست گوسفندان</h2>
        </div>
        <form onSubmit={e => {e.preventDefault();load();}} className="row g-2 mb-3" id="filterForm">
          <div className="col-md-3">
            <select className="form-select" value={filters.gender} onChange={e => setFilters(f => ({...f, gender:e.target.value}))}>
              <option value="">جنسیت</option>
              <option value="male">نر</option>
              <option value="female">ماده</option>
            </select>
          </div>
          <div className="col-md-3"><input type="number" className="form-control" placeholder="حداقل سن (ماه)" value={filters.minAge} onChange={e => setFilters(f=>({...f,minAge:e.target.value}))} /></div>
          <div className="col-md-3"><input type="number" className="form-control" placeholder="حداکثر سن (ماه)" value={filters.maxAge} onChange={e => setFilters(f=>({...f,maxAge:e.target.value}))} /></div>
          <div className="col-md-3"><button className="btn btn-primary w-100" type="submit">فیلتر</button></div>
        </form>
        <table className="table table-bordered table-hover" id="sheepTable">
          <thead className="table-light">
            <tr>
              <th>گوش1</th>
              <th>گوش2</th>
              <th>گوش3</th>
              <th>تولد</th>
              <th>جنس</th>
              <th>وضعیت تولیدمثل</th>
              <th>وضعیت سلامت</th>
              <th>نژاد</th>
            </tr>
          </thead>
          <tbody>
            {list.map(s => {
              const diffMs = Date.now() - new Date(s.dateOfBirth).getTime();
              const years = Math.floor(diffMs / (365*24*60*60*1000));
              const months = Math.floor((diffMs % (365*24*60*60*1000)) / (30*24*60*60*1000));
              return (
                <tr key={s.id} onClick={() => {setEditingId(s.id); setForm({gender:s.gender,dob:s.dateOfBirth.split('T')[0],weight:s.birthWeight||'',ear1:s.earNumber1,ear2:s.earNumber2||'',ear3:s.earNumber3||'',gen:s.fatherGen||''});}} style={{cursor:'pointer'}}>
                  <td>{s.earNumber1}</td>
                  <td>{s.earNumber2}</td>
                  <td>{s.earNumber3}</td>
                  <td>{toJalali(s.dateOfBirth)} ({years}سال {months}ماه)</td>
                  <td>{s.gender==='male'?'نر':'ماده'}</td>
                  <td>{s.reproductionState}</td>
                  <td>{s.healthState}</td>
                  <td>{s.fatherGen}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
        <form onSubmit={submit} className="card p-3 mb-5" id="sheepForm">
          <h5 className="mb-3">{editingId ? 'ویرایش گوسفند' : 'افزودن گوسفند'}</h5>
          <div className="row g-2">
            <div className="col-md-4"><select className="form-select" value={form.gender} onChange={e=>setForm({...form,gender:e.target.value})}><option value="male">نر</option><option value="female">ماده</option></select></div>
              <div className="col-md-4"><JalaliDatePicker value={form.dob} onChange={d=>setForm({...form,dob:d})} /></div>
            <div className="col-md-4"><input type="number" step="0.1" className="form-control" value={form.weight} onChange={e=>setForm({...form,weight:e.target.value})} placeholder="وزن تولد" /></div>
            <div className="col-md-4"><input className="form-control" value={form.ear1} onChange={e=>setForm({...form,ear1:e.target.value})} placeholder="گوش 1" required/></div>
            <div className="col-md-4"><input className="form-control" value={form.ear2} onChange={e=>setForm({...form,ear2:e.target.value})} placeholder="گوش 2" /></div>
            <div className="col-md-4"><input className="form-control" value={form.ear3} onChange={e=>setForm({...form,ear3:e.target.value})} placeholder="گوش 3" /></div>
            <div className="col-md-4"><input className="form-control" value={form.gen} onChange={e=>setForm({...form,gen:e.target.value})} placeholder="نژاد" /></div>
            <div className="col-md-4"><button className="btn btn-success w-100" type="submit">ذخیره</button></div>
          </div>
        </form>
      </main>
    </>
  );
}
