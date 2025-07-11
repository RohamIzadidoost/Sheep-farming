const tokenTr = localStorage.getItem('token');
if (!tokenTr) {
    window.location.href = 'index.html';
}
const headersTr = { 'Authorization': `Bearer ${tokenTr}` };

const tableTr = document.querySelector('#treatTable tbody');

function loadTreatments() {
    fetch(`${API_BASE}/sheep`, { headers: headersTr })
        .then(res => res.json())
        .then(list => {
            tableTr.innerHTML = '';
            list.forEach(s => {
                if (!s.treatments) return;
                s.treatments.forEach(t => {
                    const tr = document.createElement('tr');
                    tr.innerHTML = `<td>${s.earNumber1}</td><td>${t.diseaseDescription}</td><td>${t.treatDescription}</td><td>${t.date ? t.date.split('T')[0] : ''}</td>`;
                    tableTr.appendChild(tr);
                });
            });
        });
}

loadTreatments();
