const baseURL = 'http://localhost:8080/api/v1';

function loadHeader() {
  const headerEl = document.getElementById('header');
  if (!headerEl) return;
  fetch('header.html')
    .then(r => r.text())
    .then(html => { headerEl.innerHTML = html; });
}

function checkConnection() {
  const statusEl = document.getElementById('status');
  if (!statusEl) return;
  fetch('http://localhost:8080/')
    .then(res => {
      statusEl.textContent = res.ok ? 'ارتباط با سرور برقرار است' : 'عدم ارتباط با سرور';
    })
    .catch(() => { statusEl.textContent = 'عدم ارتباط با سرور'; });
}

function loginPage() {
  checkConnection();
  const form = document.getElementById('loginForm');
  const message = document.getElementById('loginMessage');
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    message.textContent = 'در حال ورود...';
    try {
      const res = await fetch(baseURL + '/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: document.getElementById('email').value,
          password: document.getElementById('password').value
        })
      });
      if (!res.ok) throw new Error();
      const data = await res.json();
      localStorage.setItem('token', data.token);
      message.textContent = 'ورود موفق';
      window.location.href = 'create_sheep.html';
    } catch (err) {
      message.textContent = 'ورود ناموفق';
    }
  });
}

function toGregorianDate(jdate) {
  const [jy,jm,jd] = jdate.split('-').map(Number);
  const g = jalaali.toGregorian(jy, jm, jd);
  const pad = n => n.toString().padStart(2,'0');
  return `${g.gy}-${pad(g.gm)}-${pad(g.gd)}`;
}

function authHeaders() {
  const token = localStorage.getItem('token') || '';
  return { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` };
}

function createSheepPage() {
  loadHeader();
  const form = document.getElementById('sheepForm');
  const msg = document.getElementById('sheepMsg');
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    msg.textContent = 'در حال ارسال...';
    const data = {
      name: document.getElementById('name').value,
      gender: document.getElementById('gender').value,
      dateOfBirth: toGregorianDate(document.getElementById('dob').value)
    };
    try {
      const res = await fetch(baseURL + '/sheep', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify(data)
      });
      if (!res.ok) throw new Error();
      msg.textContent = 'ثبت شد';
      form.reset();
    } catch (err) {
      msg.textContent = 'خطا در ثبت';
    }
  });
}

function createVaccinePage() {
  loadHeader();
  const form = document.getElementById('vaccineForm');
  const msg = document.getElementById('vaccineMsg');
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    msg.textContent = 'در حال ارسال...';
    const data = {
      name: document.getElementById('vname').value,
      intervalMonths: parseInt(document.getElementById('interval').value)
    };
    try {
      const res = await fetch(baseURL + '/vaccines', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify(data)
      });
      if (!res.ok) throw new Error();
      msg.textContent = 'ثبت شد';
      form.reset();
    } catch (err) {
      msg.textContent = 'خطا در ثبت';
    }
  });
}

function timersPage() {
  loadHeader();
  const list = document.getElementById('timersList');
  fetch(baseURL + '/reminders', { headers: authHeaders() })
    .then(res => res.json())
    .then(data => {
      list.innerHTML = data.map(r => `<li>${r.message}</li>`).join('');
    })
    .catch(() => { list.textContent = 'خطا در دریافت'; });
}

function sheepListPage() {
  loadHeader();
  const list = document.getElementById('sheepList');
  const genderFilter = document.getElementById('genderFilter');
  const ageFilter = document.getElementById('ageFilter');

  function loadSheep() {
    fetch(baseURL + '/sheep', { headers: authHeaders() })
      .then(res => res.json())
      .then(data => {
        const now = new Date();
        const filtered = data.filter(s => {
          let pass = true;
          if (genderFilter.value && s.gender !== genderFilter.value) pass = false;
          if (ageFilter.value) {
            const dob = new Date(s.dateOfBirth);
            const age = (now - dob) / (365.25*24*3600*1000);
            if (age < parseFloat(ageFilter.value)) pass = false;
          }
          return pass;
        });
        list.innerHTML = filtered.map(s => `<li>${s.name} - ${s.gender}</li>`).join('');
      })
      .catch(() => { list.textContent = 'خطا در دریافت'; });
  }

  genderFilter.addEventListener('change', loadSheep);
  ageFilter.addEventListener('change', loadSheep);
  loadSheep();
}

const page = document.body.id;
if (page === 'loginPage') loginPage();
if (page === 'createSheepPage') createSheepPage();
if (page === 'createVaccinePage') createVaccinePage();
if (page === 'timersPage') timersPage();
if (page === 'sheepListPage') sheepListPage();