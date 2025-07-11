const tokenVaxDef = localStorage.getItem('token');
if (!tokenVaxDef) {
    window.location.href = 'index.html';
}
const headersVaxDef = { 'Authorization': `Bearer ${tokenVaxDef}`, 'Content-Type': 'application/json' };

const tableVaxDef = document.querySelector('#vaccineDefTable tbody');
const formVaxDef = document.getElementById('vaccineDefForm');

function loadVaccineDefs() {
    fetch(`${API_BASE}/vaccines`, { headers: headersVaxDef })
        .then(res => res.json())
        .then(list => {
            tableVaxDef.innerHTML = '';
            list.forEach(v => {
                const tr = document.createElement('tr');
                tr.innerHTML = `<td>${v.name}</td><td>${v.intervalMonths}</td>`;
                tableVaxDef.appendChild(tr);
            });
        });
}

formVaxDef.addEventListener('submit', e => {
    e.preventDefault();
    const body = JSON.stringify({
        name: document.getElementById('defName').value,
        intervalMonths: parseInt(document.getElementById('defInterval').value, 10)
    });
    fetch(`${API_BASE}/vaccines`, { method: 'POST', headers: headersVaxDef, body })
        .then(() => {
            bootstrap.Modal.getInstance(document.getElementById('vaccineDefModal')).hide();
            formVaxDef.reset();
            loadVaccineDefs();
        });
});

loadVaccineDefs();
