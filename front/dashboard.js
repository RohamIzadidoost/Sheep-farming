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
            const pregnant = data.filter(s => s.reproductionState === 'pregnant').length;
            const sick = data.filter(s => s.healthState === 'sick').length;
            const treated = data.filter(s => s.healthState === 'under_treatment').length;
            let birthsThisMonth = 0;
            const now = new Date();
            data.forEach(s => {
                (s.lambings || []).forEach(l => {
                    const d = new Date(l.date);
                    if(d.getMonth() === now.getMonth() && d.getFullYear() === now.getFullYear()) birthsThisMonth += l.numBorn;
                });
            });
            const stats = [
                { label: 'تعداد کل', value: total, icon: 'bi bi-emoji-smile' },
                { label: 'آبستن', value: pregnant, icon: 'bi bi-heart-fill' },
                { label: 'متولد این ماه', value: birthsThisMonth, icon: 'bi bi-baby' },
                { label: 'بیمار', value: sick, icon: 'bi bi-emoji-frown' },
                { label: 'تحت درمان', value: treated, icon: 'bi bi-hospital' }
            ];
            const container = document.getElementById('statsCards');
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
            if (!data.length) {
                container.textContent = 'یادآوری برای نمایش وجود ندارد';
                return;
            }
            const groups = {};
            data.forEach(r => {
                if (!groups[r.type]) groups[r.type] = [];
                groups[r.type].push(r.message);
            });
            Object.keys(groups).forEach(k => {
                const card = document.createElement('div');
                card.className = 'card mb-2';
                const body = document.createElement('div');
                body.className = 'card-body';
                body.innerHTML = `<strong>${k}</strong><ul class="mb-0"></ul>`;
                groups[k].forEach(m => {
                    const li = document.createElement('li');
                    li.textContent = m;
                    body.querySelector('ul').appendChild(li);
                });
                card.appendChild(body);
                container.appendChild(card);
            });
        });
}

loadStats();
loadReminders();
