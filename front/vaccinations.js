const tokenVacc = localStorage.getItem('token');
if (!tokenVacc) {
    window.location.href = 'index.html';
}
const headersVacc = { 'Authorization': `Bearer ${tokenVacc}`, 'Content-Type': 'application/json' };

const tableVacc = document.querySelector('#vaccineTable tbody');
const formVacc = document.getElementById('vaccineForm');
const sheepSelect = document.getElementById('vaccineSheep');
const vaccineSelect = document.getElementById('vaccineName');

function toGregorianStr(jdate) {
    const [jy, jm, jd] = jdate.split('-').map(Number);
    const g = jalaali.toGregorian(jy, jm, jd);
    return `${g.gy}-${String(g.gm).padStart(2,'0')}-${String(g.gd).padStart(2,'0')}`;
}

function loadVaccinations() {
    fetch(`${API_BASE}/sheep`, { headers: headersVacc })
        .then(res => res.json())
        .then(list => {
            tableVacc.innerHTML = '';
            sheepSelect.innerHTML = '';
            list.forEach(s => {
                const opt = document.createElement('option');
                opt.value = s.id;
                opt.textContent = s.earNumber1;
                sheepSelect.appendChild(opt);
                (s.vaccinations || []).forEach(v => {
                    const tr = document.createElement('tr');
                    tr.innerHTML = `<td>${s.earNumber1}</td><td>${v.vaccine}</td><td>${v.date ? v.date.split('T')[0] : ''}</td><td>${v.vaccinator}</td><td>${v.description || ''}</td>`;
                    tableVacc.appendChild(tr);
                });
            });
        });
}

function loadVaccineDefs() {
    fetch(`${API_BASE}/vaccines`, { headers: headersVacc })
        .then(res => res.json())
        .then(list => {
            vaccineSelect.innerHTML = '';
            list.forEach(v => {
                const opt = document.createElement('option');
                opt.value = v.id;
                opt.textContent = v.name;
                vaccineSelect.appendChild(opt);
            });
        });
}

formVacc.addEventListener('submit', e => {
    e.preventDefault();
    const sheepID = sheepSelect.value;
    fetch(`${API_BASE}/sheep/${sheepID}`, { headers: headersVacc })
        .then(res => res.json())
        .then(sheep => {
            const vaccinations = sheep.vaccinations || [];
            vaccinations.push({
                vaccine: vaccineSelect.value,
                vaccinator: document.getElementById('vaccineVaccinator').value,
                description: document.getElementById('vaccineDesc').value,
                date: toGregorianStr(document.getElementById('vaccineDate').value)
            });
            return fetch(`${API_BASE}/sheep/${sheepID}`, {
                method: 'PUT',
                headers: headersVacc,
                body: JSON.stringify({ vaccinations })
            });
        })
        .then(() => {
            bootstrap.Modal.getInstance(document.getElementById('vaccineModal')).hide();
            formVacc.reset();
            loadVaccinations();
        });
});

loadVaccineDefs();
loadVaccinations();
