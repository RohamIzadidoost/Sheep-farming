// Shared scripts for dashboard and management pages

document.addEventListener('DOMContentLoaded', () => {
  const token = localStorage.getItem('token');
  if (!token && document.getElementById('loginForm') === null) {
    window.location.href = 'index.html';
    return;
  }

  if (document.getElementById('reminderList')) {
    loadDashboard();
  }
});

function loadDashboard() {
  // Example static data; replace with real API calls
  document.getElementById('statTotal').textContent = '120';
  document.getElementById('statPregnant').textContent = '25';
  document.getElementById('statVaccine').textContent = '8';
  document.getElementById('statOverdue').textContent = '3';

  const reminderList = document.getElementById('reminderList');
  const items = [
    {icon: 'fa-syringe', color: 'success', text: 'واکسیناسیون برای گوسفند شماره 12'},
    {icon: 'fa-baby', color: 'warning', text: 'بررسی آبستنی گوسفند شماره 8'},
    {icon: 'fa-cut', color: 'info', text: 'پشم‌چینی گوسفند شماره 4'}
  ];
  items.forEach(it => {
    const div = document.createElement('div');
    div.className = 'col-12 col-md-6 col-lg-4';
    div.innerHTML = `<div class="alert alert-${it.color} d-flex align-items-center" role="alert">
        <i class="fa ${it.icon} me-2"></i>
        <div>${it.text}</div>
      </div>`;
    reminderList.appendChild(div);
  });
}
