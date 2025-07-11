const tokenSheep = localStorage.getItem('token');
if (!tokenSheep) {
    window.location.href = 'index.html';
}
const headersSheep = { 'Authorization': `Bearer ${tokenSheep}`, 'Content-Type': 'application/json' };

const tableBody = document.querySelector('#sheepTable tbody');
const form = document.getElementById('sheepForm');
let editingId = null;
let currentSheep = null;

function toGregorianStr(jdate) {
    const [jy,jm,jd] = jdate.split('-').map(Number);
    const g = jalaali.toGregorian(jy,jm,jd);
    return `${g.gy}-${String(g.gm).padStart(2,'0')}-${String(g.gd).padStart(2,'0')}`;
}

function loadSheep() {
    fetch(`${API_BASE}/sheep`, { headers: headersSheep })
        .then(res => res.json())
        .then(list => {
            tableBody.innerHTML = '';
            list.forEach(s => {
                const tr = document.createElement('tr');
                const diffMs = Date.now() - new Date(s.dateOfBirth).getTime();
                const years = Math.floor(diffMs / (365*24*60*60*1000));
                const months = Math.floor((diffMs % (365*24*60*60*1000)) / (30*24*60*60*1000));
                const age = `${years} سال ${months} ماه`;
                tr.innerHTML = `
                    <td>${s.earNumber1}</td>
                    <td>${s.earNumber2 || ''}</td>
                    <td>${s.earNumber3 || ''}</td>
                    <td>${age}</td>
                    <td>${s.gender === 'male' ? 'نر' : 'ماده'}</td>
                    <td>${s.reproductionState}</td>
                    <td>${s.healthState}</td>
                    <td>${s.fatherGen || ''}</td>
                    <td><button class="btn btn-sm btn-info" onclick="showSheep('${s.id}')">مشاهده</button></td>`;
                tableBody.appendChild(tr);
            });
        });
}

window.editSheep = function(id) {
    fetch(`${API_BASE}/sheep/${id}`, { headers: headersSheep })
        .then(res => res.json())
        .then(s => {
            editingId = id;
            document.getElementById('sheepGender').value = s.gender;
            document.getElementById('sheepDob').value = s.dateOfBirth.split('T')[0];
            document.getElementById('ear1').value = s.earNumber1;
            document.getElementById('ear2').value = s.earNumber2 || '';
            document.getElementById('ear3').value = s.earNumber3 || '';
            document.getElementById('gen').value = s.fatherGen || '';
            new bootstrap.Modal(document.getElementById('sheepModal')).show();
        });
}

window.deleteSheep = function(id) {
    if (!confirm('حذف شود؟')) return;
    fetch(`${API_BASE}/sheep/${id}`, { method: 'DELETE', headers: headersSheep })
        .then(() => loadSheep());
}

window.showSheep = function(id) {
    fetch(`${API_BASE}/sheep/${id}`, { headers: headersSheep })
        .then(res => res.json())
        .then(s => {
            currentSheep = s;
            const detail = document.getElementById('sheepDetails');
            detail.innerHTML = `
                <img src="${s.photoUrl || 'https://cdn.jsdelivr.net/gh/twitter/twemoji/assets/svg/1f411.svg'}" class="img-thumbnail mb-3" style="max-width:150px">
                <div>گوش 1: ${s.earNumber1}</div>
                <div>گوش 2: ${s.earNumber2 || ''}</div>
                <div>گوش 3: ${s.earNumber3 || ''}</div>
                <div>شماره پلاک: ${s.neckNumber || ''}</div>
                <div>تاریخ تولد: ${s.dateOfBirth.split('T')[0]}</div>
                <div>وزن تولد: ${s.birthWeight}</div>
                <div>نژاد: ${s.fatherGen}</div>`;
            document.getElementById('stateReproduction').value = s.reproductionState;
            document.getElementById('stateHealth').value = s.healthState;
            fetch(`${API_BASE}/vaccines`, { headers: headersSheep })
                .then(r => r.json())
                .then(vlist => {
                    const select = document.getElementById('detailVaccine');
                    select.innerHTML = '';
                    vlist.forEach(v => {
                        const opt = document.createElement('option');
                        opt.value = v.id;
                        opt.textContent = v.name;
                        select.appendChild(opt);
                    });
                    const vaccList = document.getElementById('vaccList');
                    vaccList.innerHTML = '';
                    (s.vaccinations || []).forEach(v => {
                        const li = document.createElement('li');
                        li.textContent = `${v.vaccine} - ${v.date.split('T')[0]}`;
                        vaccList.appendChild(li);
                    });
                    const treatList = document.getElementById('treatList');
                    treatList.innerHTML = '';
                    (s.treatments || []).forEach(t => {
                        const li = document.createElement('li');
                        li.textContent = `${t.diseaseDescription} - ${t.date.split('T')[0]}`;
                        treatList.appendChild(li);
                    });
                    const lambList = document.getElementById('lambList');
                    lambList.innerHTML = '';
                    (s.lambings || []).forEach(l => {
                        const li = document.createElement('li');
                        li.textContent = `${l.date.split('T')[0]} - ${l.numBorn}`;
                        lambList.appendChild(li);
                    });
                    bootstrap.Modal.getOrCreateInstance(document.getElementById('detailModal')).show();
                });
        });
}

document.getElementById('stateForm').addEventListener('submit', e => {
    e.preventDefault();
    if (!currentSheep) return;
    fetch(`${API_BASE}/sheep/${currentSheep.id}`, {
        method: 'PUT',
        headers: headersSheep,
        body: JSON.stringify({
            reproductionState: document.getElementById('stateReproduction').value,
            healthState: document.getElementById('stateHealth').value
        })
    }).then(() => {
        bootstrap.Modal.getInstance(document.getElementById('detailModal')).hide();
        loadSheep();
    });
});

document.getElementById('treatmentForm').addEventListener('submit', e => {
    e.preventDefault();
    if (!currentSheep) return;
    const body = JSON.stringify({
        diseaseDescription: document.getElementById('diseaseDesc').value,
        treatDescription: document.getElementById('treatDesc').value,
        date: toGregorianStr(document.getElementById('treatDate').value)
    });
    fetch(`${API_BASE}/sheep/${currentSheep.id}/treatments`, {
        method: 'POST',
        headers: headersSheep,
        body
    }).then(() => {
        bootstrap.Modal.getInstance(document.getElementById('detailModal')).hide();
        loadSheep();
    });
});

document.getElementById('vaccForm').addEventListener('submit', e => {
    e.preventDefault();
    if (!currentSheep) return;
    const body = JSON.stringify({
        vaccine: document.getElementById('detailVaccine').value,
        vaccinator: document.getElementById('detailVaccinator').value,
        description: document.getElementById('detailVDesc').value,
        date: toGregorianStr(document.getElementById('detailVDate').value)
    });
    fetch(`${API_BASE}/sheep/${currentSheep.id}/vaccinations`, {
        method: 'POST',
        headers: headersSheep,
        body
    }).then(() => {
        bootstrap.Modal.getInstance(document.getElementById('detailModal')).hide();
        loadSheep();
    });
});

document.getElementById('lambForm').addEventListener('submit', e => {
    e.preventDefault();
    if (!currentSheep) return;
    const males = parseInt(document.getElementById('lambMale').value,10)||0;
    const females = parseInt(document.getElementById('lambFemale').value,10)||0;
    const numDead = parseInt(document.getElementById('lambDead').value,10)||0;
    const sexes = [];
    for(let i=0;i<males;i++) sexes.push('male');
    for(let i=0;i<females;i++) sexes.push('female');
    const body = JSON.stringify({
        date: toGregorianStr(document.getElementById('lambDate').value),
        numBorn: males + females,
        sexes,
        numDead
    });
    fetch(`${API_BASE}/sheep/${currentSheep.id}/lambings`, {
        method: 'POST',
        headers: headersSheep,
        body
    }).then(() => {
        bootstrap.Modal.getInstance(document.getElementById('detailModal')).hide();
        loadSheep();
    });
});

form.addEventListener('submit', e => {
    e.preventDefault();
    const body = JSON.stringify({
        gender: document.getElementById('sheepGender').value,
        dateOfBirth: document.getElementById('sheepDob').value,
        earNumber1: document.getElementById('ear1').value,
        earNumber2: document.getElementById('ear2').value,
        earNumber3: document.getElementById('ear3').value,
        fatherGen: document.getElementById('gen').value
    });
    const method = editingId ? 'PUT' : 'POST';
    const url = editingId ? `${API_BASE}/sheep/${editingId}` : `${API_BASE}/sheep`;
    fetch(url, { method, headers: headersSheep, body })
        .then(() => {
            bootstrap.Modal.getInstance(document.getElementById('sheepModal')).hide();
            editingId = null;
            form.reset();
            loadSheep();
        });
});

loadSheep();
