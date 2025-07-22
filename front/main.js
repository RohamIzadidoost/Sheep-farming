import moment from "moment-jalaali";

moment.loadPersian({ dialect: "persian-modern", usePersianDigits: true });
console.log(moment().format("jYYYY/jMM/jDD"));

window.DateConverter = function(miladiDate) {
  return moment(miladiDate).format("jYYYY/jMM/jDD");
};

window.toGregorianStr = function(jdate) {
  const [jy, jm, jd] = jdate.split(/[-\/]/).map(Number);
  const g = window.jalaali.toGregorian(jy, jm, jd);
  return `${g.gy}-${String(g.gm).padStart(2, "0")}-${String(g.gd).padStart(2, "0")}`;
};

const statusEl = document.getElementById('status');
const loginForm = document.getElementById('loginForm');
const loginMessage = document.getElementById('loginMessage');
// const API_BASE = 'http://82.115.17.206:8080/api/v1';
const API_BASE = 'http://localhost:8080/api/v1' ; 

// Check connection to backend
fetch('http://localhost:8080')
  .then(res => {
    if (res.ok) {
      statusEl.textContent = 'ارتباط با سرور برقرار است';
    } else {
      statusEl.textContent = 'عدم ارتباط با سرور';
    }
  })
  .catch(() => {
    statusEl.textContent = 'عدم ارتباط با سرور';
  });

loginForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  loginMessage.textContent = 'در حال ورود...';
  try {
    const res = await fetch(`${API_BASE}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        email: document.getElementById('email').value,
        password: document.getElementById('password').value
      })
    });
    if (!res.ok) throw new Error('خطا در ورود');
    const data = await res.json();
    localStorage.setItem('token', data.token);
    loginMessage.textContent = 'ورود موفق';
    window.location.href = 'dashboard.html';
  } catch (err) {
    loginMessage.textContent = 'ورود ناموفق';
  }
});
