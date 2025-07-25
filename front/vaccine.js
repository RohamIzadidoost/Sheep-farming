const tokenVax = localStorage.getItem('token');
if (!tokenVax) {
    window.location.href = 'index.html';
}
const headersVax = { 'Authorization': `Bearer ${tokenVax}`, 'Content-Type': 'application/json' };

const tableVax = document.querySelector('#vaccineTable tbody');
const formVax = document.getElementById('vaccineForm');

function loadVaccines() {
    fetch(`${API_BASE}/vaccines`, { headers: headersVax })
        .then(res => res.json())
        .then(list => {
            tableVax.innerHTML = '';
            list.forEach(v => {
                const tr = document.createElement('tr');
                tr.innerHTML = `<td>${v.sheepID}</td><td>${v.name}</td><td>${DateConverter(v.date)}</td><td>${DateConverter(v.nextDose)}</td><td>${v.description || ''}</td>`;
                tableVax.appendChild(tr);
            });
        });
}

formVax.addEventListener('submit', e => {
    e.preventDefault();
    const body = JSON.stringify({
        name: document.getElementById('vaccineName').value,
        sheepID: document.getElementById('vaccineSheep').value,
        date: toGregorianStr(document.getElementById('vaccineDate').value),
        nextDose: toGregorianStr(document.getElementById('vaccineNext').value),
        description: document.getElementById('vaccineNote').value
    });
    fetch(`${API_BASE}/vaccines`, { method: 'POST', headers: headersVax, body })
        .then(() => {
            bootstrap.Modal.getInstance(document.getElementById('vaccineModal')).hide();
            formVax.reset();
            loadVaccines();
        });
});

loadVaccines();
