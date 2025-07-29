import React, { useEffect, useState } from 'react';
import NavBar from '../components/NavBar';
import { apiFetch } from '../utils/api';

export default function DashboardPage() {
  const [stats, setStats] = useState([]);
  const [reminders, setReminders] = useState({});

  useEffect(() => {
    apiFetch('/sheep')
      .then(res => res.json())
      .then(list => {
        const total = list.length;
        const pregnant = list.filter(s => s.reproductionState === 'pregnant').length;
        const sick = list.filter(s => s.healthState === 'sick').length;
        const treated = list.filter(s => s.healthState === 'under_treatment').length;
        let birthsThisMonth = 0;
        const now = new Date();
        list.forEach(s => {
          if (Array.isArray(s.lambings)) {
            s.lambings.forEach(l => {
              const d = new Date(l.date);
              if (d.getMonth() === now.getMonth() && d.getFullYear() === now.getFullYear()) 
                birthsThisMonth += l.numBorn;
            });
          }
        });
        setStats([
          { label: 'تعداد کل', value: total, icon: 'bi bi-emoji-smile' },
          { label: 'آبستن', value: pregnant, icon: 'bi bi-heart-fill' },
          { label: 'متولد این ماه', value: birthsThisMonth, icon: 'bi bi-baby' },
          { label: 'بیمار', value: sick, icon: 'bi bi-emoji-frown' },
          { label: 'تحت درمان', value: treated, icon: 'bi bi-hospital' }
        ]);
      });

    apiFetch('/reminders')
      .then(res => res.json())
      .then(list => {
        const groups = {};
        if (Array.isArray(list)) {
          list.forEach(r => {
            if (r && r.type) {
              if (!groups[r.type]) groups[r.type] = [];
              groups[r.type].push(r.message);
            }
          });
        }
        setReminders(groups);
      });
  }, []);

  return (
    <>
      <NavBar />
      <main className="container pt-5 mt-4">
        <div className="row g-3" id="statsCards">
          {stats.map(s => (
            <div className="col-6 col-md-3" key={s.label}>
              <div className="card text-center shadow-sm">
                <div className="card-body">
                  <i className={`${s.icon} fs-1 text-primary`}></i>
                  <h4 className="mt-2">{s.value}</h4>
                  <p className="mb-0">{s.label}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
        <h2 className="mt-4">یادآورها</h2>
        <div id="reminders" className="mt-3">
          {Object.keys(reminders).length === 0 && <p>یادآوری برای نمایش وجود ندارد</p>}
          {Object.keys(reminders).map(k => (
            <div className="card mb-2" key={k}>
              <div className="card-body">
                <strong>{k}</strong>
                <ul className="mb-0">
                  {Array.isArray(reminders[k]) && reminders[k].map((m,i) => <li key={i}>{m}</li>)}
                </ul>
              </div>
            </div>
          ))}
        </div>
      </main>
    </>
  );
}
