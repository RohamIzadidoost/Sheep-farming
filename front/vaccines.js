const tokenVDef = localStorage.getItem('token');
if (!tokenVDef) {
    window.location.href = 'index.html';
}
const headersVDef = { 'Authorization': `Bearer ${tokenVDef}`, 'Content-Type': 'application/json' };

const vacDefTable = document.querySelector('#vacDefTable tbody');
const vacDefForm = document.getElementById('vacDefForm');

function loadVacDefs() {
    fetch(`${API_BASE}/vaccines`, { headers: headersVDef })
        .then(res => res.json())
        .then(list => {
            vacDefTable.innerHTML = '';
            list.forEach(v => {
                const tr = document.createElement('tr');
                tr.innerHTML = `<td>${v.name}</td><td>${v.intervalMonths}</td>`;
                vacDefTable.appendChild(tr);
            });
        });
}

vacDefForm.addEventListener('submit', e => {
    e.preventDefault();
    const body = JSON.stringify({
        name: document.getElementById('vacDefName').value,
        intervalMonths: parseInt(document.getElementById('vacDefInterval').value, 10)
    });
    fetch(`${API_BASE}/vaccines`, { method: 'POST', headers: headersVDef, body })
        .then(() => {
            bootstrap.Modal.getInstance(document.getElementById('vacDefModal')).hide();
            vacDefForm.reset();
            loadVacDefs();
        });
});

loadVacDefs();
