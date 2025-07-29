import React, { useEffect, useState } from 'react';
import NavBar from '../components/NavBar';
import { apiFetch } from '../utils/api';
import { toJalali } from '../utils/jdate';

export default function TreatmentsPage() {
  const [list, setList] = useState([]);

  const load = () => {
    apiFetch('/sheep')
      .then(res => res.json())
      .then(sheep => {
        const rows = [];
        sheep.forEach(s => {
          (s.treatments || []).forEach(t => rows.push({ ear: s.earNumber1, ...t }));
        });
        setList(rows);
      });
  };

  useEffect(() => { load(); }, []);

  return (
    <>
      <NavBar />
      <main className="container pt-5 mt-4">
        <h2 className="mb-3">سوابق درمان</h2>
        <table className="table table-bordered table-hover" id="treatTable">
          <thead className="table-light"><tr><th>گوسفند</th><th>بیماری</th><th>توضیح درمان</th><th>تاریخ</th></tr></thead>
          <tbody>
            {list.map((t,i) => <tr key={i}><td>{t.ear}</td><td>{t.diseaseDescription}</td><td>{t.treatDescription}</td><td>{toJalali(t.date)}</td></tr>)}
          </tbody>
        </table>
      </main>
    </>
  );
}
