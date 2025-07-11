const token = localStorage.getItem('token');
if (!token) {
    window.location.href = 'index.html';
}

const headers = { 'Authorization': `Bearer ${token}` };

function loadStats() {
    fetch(`${API_BASE}/sheep`, { headers })
        .then(res => res.json())
        .then(data => {
            const total = data.length;
            const now = new Date();
            let birthsThisMonth = 0;
            const pregnant = data.filter(s => s.reproductionState === 'pregnant').length;
            data.forEach(s => {
                (s.lambings || []).forEach(l => {
                    const d = new Date(l.date);
                    if (d.getFullYear() === now.getFullYear() && d.getMonth() === now.getMonth()) {
                        birthsThisMonth += l.numBorn;
                    }
                });
            });
            const sick = data.filter(s => s.healthState === 'sick').length;
            const treated = data.filter(s => s.healthState === 'under_treatment').length;
            const stats = [
                { label: 'تعداد کل', value: total, icon: 'bi bi-emoji-smile' },
                { label: 'آبستن', value: pregnant, icon: 'bi bi-heart-fill' },
                { label: 'تولد این ماه', value: birthsThisMonth, icon: 'bi bi-gift' },
                { label: 'بیمار', value: sick, icon: 'bi bi-thermometer-half' },
                { label: 'در حال درمان', value: treated, icon: 'bi bi-hospital' }
            ];
            const container = document.getElementById('statsCards');
            container.innerHTML = '';
            stats.forEach(s => {
                const div = document.createElement('div');
                div.className = 'col-6 col-md-3';
                div.innerHTML = `<div class="card text-center shadow-sm">
                    <div class="card-body">
                        <i class="${s.icon} fs-1 text-primary"></i>
                        <h4 class="mt-2">${s.value}</h4>
                        <p class="mb-0">${s.label}</p>
                    </div>
                </div>`;
                container.appendChild(div);
            });
        });
}

function loadReminders() {
    fetch(`${API_BASE}/reminders`, { headers })
        .then(res => res.json())
        .then(data => {
            const container = document.getElementById('reminders');
            container.innerHTML = '';
            if (!data.length) {
                container.textContent = 'یادآوری برای نمایش وجود ندارد';
                return;
            }
            data.forEach(r => {
                const alert = document.createElement('div');
                alert.className = 'alert alert-warning d-flex align-items-center';
                alert.innerHTML = `<i class="bi bi-bell-fill me-2"></i><div>${r.message}</div>`;
                container.appendChild(alert);
            });
        });
}

loadStats();
loadReminders();
