const tokenSheep = localStorage.getItem('token');
if (!tokenSheep) {
    window.location.href = 'index.html';
}
const headersSheep = { 'Authorization': `Bearer ${tokenSheep}`, 'Content-Type': 'application/json' };

const tableBody = document.querySelector('#sheepTable tbody');
const form = document.getElementById('sheepForm');
let editingId = null;

function loadSheep() {
    fetch(`${API_BASE}/sheep`, { headers: headersSheep })
        .then(res => res.json())
        .then(list => {
            tableBody.innerHTML = '';
            list.forEach(s => {
                const tr = document.createElement('tr');
                const age = new Date().getFullYear() - new Date(s.dateOfBirth).getFullYear();
                tr.innerHTML = `
                    <td>${s.name}</td>
                    <td>${age}</td>
                    <td>${s.gender === 'male' ? 'نر' : 'ماده'}</td>
                    <td>${s.id}</td>
                    <td>${s.breedingDate ? 'آبستن' : '-'}</td>
                    <td>
                        <button class="btn btn-sm btn-primary me-1" onclick="editSheep('${s.id}')"><i class="bi bi-pencil"></i></button>
                        <button class="btn btn-sm btn-danger" onclick="deleteSheep('${s.id}')"><i class="bi bi-trash"></i></button>
                    </td>`;
                tableBody.appendChild(tr);
            });
        });
}

window.editSheep = function(id) {
    fetch(`${API_BASE}/sheep/${id}`, { headers: headersSheep })
        .then(res => res.json())
        .then(s => {
            editingId = id;
            document.getElementById('sheepName').value = s.name;
            document.getElementById('sheepGender').value = s.gender;
            document.getElementById('sheepDob').value = s.dateOfBirth.split('T')[0];
            document.getElementById('sheepTag').value = s.id;
            new bootstrap.Modal(document.getElementById('sheepModal')).show();
        });
}

window.deleteSheep = function(id) {
    if (!confirm('حذف شود؟')) return;
    fetch(`${API_BASE}/sheep/${id}`, { method: 'DELETE', headers: headersSheep })
        .then(() => loadSheep());
}

form.addEventListener('submit', e => {
    e.preventDefault();
    const body = JSON.stringify({
        name: document.getElementById('sheepName').value,
        gender: document.getElementById('sheepGender').value,
        dateOfBirth: document.getElementById('sheepDob').value
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
