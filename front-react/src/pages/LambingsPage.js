import React, { useEffect, useState } from 'react';
import NavBar from '../components/NavBar';
import { apiFetch } from '../utils/api';
import JalaliDatePicker from '../components/JalaliDatePicker';
import { toJalali } from '../utils/jdate';

export default function LambingsPage() {
  const [list, setList] = useState([]);
  const [from, setFrom] = useState('');
  const [to, setTo] = useState('');

  const load = () => {
    const params = new URLSearchParams();
    if (from) params.append('from', from);
    if (to) params.append('to', to);
    const qs = params.toString() ? `?${params.toString()}` : '';
    apiFetch(`/lambings${qs}`)
      .then(res => res.json())
      .then(setList);
  };

  useEffect(() => { load(); }, []);

  return (
    <>
      <NavBar />
      <main className="container pt-5 mt-4">
        <div className="mb-3">
          <form onSubmit={e=>{e.preventDefault();load();}} className="row g-2">
            <div className="col-md-4">
              <JalaliDatePicker value={from} onChange={setFrom} />
            </div>
            <div className="col-md-4">
              <JalaliDatePicker value={to} onChange={setTo} />
            </div>
            <div className="col-md-4">
              <button className="btn btn-primary w-100" type="submit">فیلتر</button>
            </div>
          </form>
        </div>
        <table className="table table-bordered" id="lambTable">
          <thead className="table-light">
            <tr>
              <th>گوسفند</th>
              <th>تاریخ</th>
              <th>تعداد تولد</th>
              <th>جنسیت‌ها</th>
              <th>تلفات</th>
            </tr>
          </thead>
          <tbody>
            {list.map((l,i) => (
              <tr key={i}>
                <td>{l.sheepEar || ''}</td>
                <td>{toJalali(l.date)}</td>
                <td>{l.numBorn}</td>
                <td>نر:{l.numMaleBorn} ماده:{l.numFemaleBorn}</td>
                <td>{l.numDead}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </main>
    </>
  );
}
