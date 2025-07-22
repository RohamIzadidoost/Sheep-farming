import moment from "moment-jalaali";
moment.loadPersian({ dialect: "persian-modern", usePersianDigits: true });

function DateConverter(miladiDate) {
  return `${moment(miladiDate).format("jYYYY/jMM/jDD")}`;
}

const tokenSheep = localStorage.getItem("token");
if (!tokenSheep) {
  window.location.href = "index.html";
}
const headersSheep = {
  Authorization: `Bearer ${tokenSheep}`,
  "Content-Type": "application/json",
};

const tableBody = document.querySelector("#sheepTable tbody");
const form = document.getElementById("sheepForm");
let editingId = null;
let currentSheep = null;

function toGregorianStr(jdate) {
  const [jy, jm, jd] = jdate.split("-").map(Number);
  const g = jalaali.toGregorian(jy, jm, jd);
  return `${g.gy}-${String(g.gm).padStart(2, "0")}-${String(g.gd).padStart(2, "0")}`;
}

function loadSheep() {
  const params = new URLSearchParams();
  const g = document.getElementById("filterGender").value;
  if (g) params.append("gender", g);
  const minA = document.getElementById("filterMinAge").value;
  if (minA) params.append("minAgeDays", parseInt(minA) * 30);
  const maxA = document.getElementById("filterMaxAge").value;
  if (maxA) params.append("maxAgeDays", parseInt(maxA) * 30);
  const qs = params.toString() ? `?${params.toString()}` : "";
  fetch(`${API_BASE}/sheep${qs}`, { headers: headersSheep })
    .then((res) => res.json())
    .then((list) => {
      tableBody.innerHTML = "";
      list.forEach((s) => {
        const tr = document.createElement("tr");
        // Use DateConverter for Jalali date display
        const dobJalali = DateConverter(s.dateOfBirth);
        const diffMs = Date.now() - new Date(s.dateOfBirth).getTime();
        const years = Math.floor(diffMs / (365 * 24 * 60 * 60 * 1000));
        const months = Math.floor(
          (diffMs % (365 * 24 * 60 * 60 * 1000)) / (30 * 24 * 60 * 60 * 1000),
        );
        const age = `${years} سال ${months} ماه`;
        tr.innerHTML = `
                    <td>${s.earNumber1}</td>
                    <td>${s.earNumber2 || ""}</td>
                    <td>${s.earNumber3 || ""}</td>
                    <td>${dobJalali}</td>
                    <td>${s.gender === "male" ? "نر" : "ماده"}</td>
                    <td>${s.reproductionState}</td>
                    <td>${s.healthState}</td>
                    <td>${s.fatherGen || ""}</td>
                    <td><button class="btn btn-sm btn-info" onclick="showSheep('${s.id}')">مشاهده</button></td>`;
        tableBody.appendChild(tr);
      });
    });
}

window.editSheep = function (id) {
  fetch(`${API_BASE}/sheep/${id}`, { headers: headersSheep })
    .then((res) => res.json())
    .then((s) => {
      editingId = id;
      document.getElementById("sheepGender").value = s.gender;
      document.getElementById("sheepDob").value = s.dateOfBirth.split("T")[0];
      document.getElementById("birthWeight").value = s.birthWeight || "";
      document.getElementById("ear1").value = s.earNumber1;
      document.getElementById("ear2").value = s.earNumber2 || "";
      document.getElementById("ear3").value = s.earNumber3 || "";
      document.getElementById("gen").value = s.fatherGen || "";
      new bootstrap.Modal(document.getElementById("sheepModal")).show();
    });
};

window.deleteSheep = function (id) {
  if (!confirm("حذف شود؟")) return;
  fetch(`${API_BASE}/sheep/${id}`, {
    method: "DELETE",
    headers: headersSheep,
  }).then(() => loadSheep());
};

window.showSheep = function (id) {
  fetch(`${API_BASE}/sheep/${id}`, { headers: headersSheep })
    .then((res) => res.json())
    .then((s) => {
      currentSheep = s;
      const detail = document.getElementById("sheepDetails");
      detail.innerHTML = `
                <img src="${s.photoUrl || "https://cdn.jsdelivr.net/gh/twitter/twemoji/assets/svg/1f411.svg"}" class="img-thumbnail mb-3" style="max-width:150px">
                <div>گوش 1: ${s.earNumber1}</div>
                <div>گوش 2: ${s.earNumber2 || ""}</div>
                <div>گوش 3: ${s.earNumber3 || ""}</div>
                <div>شماره پلاک: ${s.neckNumber || ""}</div>
                <div>تاریخ تولد: ${DateConverter(s.dateOfBirth)}</div>
                <div>وزن تولد: ${s.birthWeight}</div>
                <div>نژاد: ${s.fatherGen}</div>`;
      document.getElementById("stateReproduction").value = s.reproductionState;
      document.getElementById("stateHealth").value = s.healthState;
      fetch(`${API_BASE}/vaccines`, { headers: headersSheep })
        .then((r) => r.json())
        .then((vlist) => {
          const select = document.getElementById("detailVaccine");
          select.innerHTML = "";
          vlist.forEach((v) => {
            const opt = document.createElement("option");
            opt.value = v.id;
            opt.textContent = v.name;
            select.appendChild(opt);
          });
          const vaccList = document.getElementById("vaccList");
          vaccList.innerHTML = "";
          (s.vaccinations || []).forEach((v, i) => {
            const tr = document.createElement("tr");
            tr.innerHTML = `
              <td>${v.vaccine}</td>
              <td>${v.date.split("T")[0]}</td>
              <td>${v.vaccinator || ""}</td>
              <td>${v.description || ""}</td>
              <td><button class="btn btn-sm btn-danger" onclick="deleteVacc(${i})">حذف</button></td>`;
            vaccList.appendChild(tr);
          });
          const treatList = document.getElementById("treatList");
          treatList.innerHTML = "";
          (s.treatments || []).forEach((t, i) => {
            const tr = document.createElement("tr");
            tr.innerHTML = `
              <td>${t.diseaseDescription}</td>
              <td>${t.treatDescription}</td>
              <td>${t.date.split("T")[0]}</td>
              <td><button class="btn btn-sm btn-danger" onclick="deleteTreat(${i})">حذف</button></td>`;
            treatList.appendChild(tr);
          });
          const lambList = document.getElementById("lambList");
          lambList.innerHTML = "";
          (s.lambings || []).forEach((l, i) => {
            const tr = document.createElement("tr");
            tr.innerHTML = `
              <td>${l.date.split("T")[0]}</td>
              <td>${l.numBorn}</td>
              <td>${l.numMaleBorn}</td>
              <td>${l.numFemaleBorn}</td>
              <td>${l.numDead}</td>
              <td><button class="btn btn-sm btn-danger" onclick="deleteLamb(${i})">حذف</button></td>`;
            lambList.appendChild(tr);
          });
          bootstrap.Modal.getOrCreateInstance(
            document.getElementById("detailModal"),
          ).show();
        });
    });
};

document.getElementById("stateForm").addEventListener("submit", (e) => {
  e.preventDefault();
  if (!currentSheep) return;
  fetch(`${API_BASE}/sheep/${currentSheep.id}`, {
    method: "PUT",
    headers: headersSheep,
    body: JSON.stringify({
      reproductionState: document.getElementById("stateReproduction").value,
      healthState: document.getElementById("stateHealth").value,
    }),
  }).then(() => {
    bootstrap.Modal.getInstance(document.getElementById("detailModal")).hide();
    loadSheep();
  });
});

document.getElementById("treatmentForm").addEventListener("submit", (e) => {
  e.preventDefault();
  if (!currentSheep) return;
  const body = JSON.stringify({
    diseaseDescription: document.getElementById("diseaseDesc").value,
    treatDescription: document.getElementById("treatDesc").value,
    date: toGregorianStr(document.getElementById("treatDate").value),
  });
  fetch(`${API_BASE}/sheep/${currentSheep.id}/treatments`, {
    method: "POST",
    headers: headersSheep,
    body,
  }).then(() => {
    bootstrap.Modal.getInstance(document.getElementById("detailModal")).hide();
    loadSheep();
  });
});

document.getElementById("vaccForm").addEventListener("submit", (e) => {
  e.preventDefault();
  if (!currentSheep) return;
  const body = JSON.stringify({
    vaccine: document.getElementById("detailVaccine").value,
    vaccinator: document.getElementById("detailVaccinator").value,
    description: document.getElementById("detailVDesc").value,
    date: toGregorianStr(document.getElementById("detailVDate").value),
  });
  fetch(`${API_BASE}/sheep/${currentSheep.id}/vaccinations`, {
    method: "POST",
    headers: headersSheep,
    body,
  }).then(() => {
    bootstrap.Modal.getInstance(document.getElementById("detailModal")).hide();
    loadSheep();
  });
});

document.getElementById("lambForm").addEventListener("submit", (e) => {
  e.preventDefault();
  if (!currentSheep) return;
  const males = parseInt(document.getElementById("lambMale").value, 10) || 0;
  const females = parseInt(document.getElementById("lambFemale").value, 10) || 0;
  const numDead = parseInt(document.getElementById("lambDead").value, 10) || 0;
  const body = JSON.stringify({
    date: toGregorianStr(document.getElementById("lambDate").value),
    numBorn: males + females,
    numMaleBorn: males,
    numFemaleBorn: females,
    numDead: numDead,
  });
  fetch(`${API_BASE}/sheep/${currentSheep.id}/lambings`, {
    method: "POST",
    headers: headersSheep,
    body,
  }).then(() => {
    bootstrap.Modal.getInstance(document.getElementById("detailModal")).hide();
    loadSheep();
  });
});

form.addEventListener("submit", (e) => {
  e.preventDefault();
  const body = JSON.stringify({
    gender: document.getElementById("sheepGender").value,
    dateOfBirth: toGregorianStr(document.getElementById("sheepDob").value),
    birthWeight: parseFloat(document.getElementById("birthWeight").value) || 0,
    earNumber1: document.getElementById("ear1").value,
    earNumber2: document.getElementById("ear2").value,
    earNumber3: document.getElementById("ear3").value,
    fatherGen: document.getElementById("gen").value,
  });
  const method = editingId ? "PUT" : "POST";
  const url = editingId
    ? `${API_BASE}/sheep/${editingId}`
    : `${API_BASE}/sheep`;
  fetch(url, { method, headers: headersSheep, body }).then(() => {
    bootstrap.Modal.getInstance(document.getElementById("sheepModal")).hide();
    editingId = null;
    form.reset();
    loadSheep();
  });
});

loadSheep();

document.getElementById("filterForm").addEventListener("submit", (e) => {
  e.preventDefault();
  loadSheep();
});
window.deleteVacc = function (idx) {
  if (!currentSheep) return;
  fetch(`${API_BASE}/sheep/${currentSheep.id}/vaccinations/${idx}`, {
    method: "DELETE",
    headers: headersSheep,
  }).then(() => showSheep(currentSheep.id));
};

window.deleteTreat = function (idx) {
  if (!currentSheep) return;
  fetch(`${API_BASE}/sheep/${currentSheep.id}/treatments/${idx}`, {
    method: "DELETE",
    headers: headersSheep,
  }).then(() => showSheep(currentSheep.id));
};

window.deleteLamb = function (idx) {
  if (!currentSheep) return;
  fetch(`${API_BASE}/sheep/${currentSheep.id}/lambings/${idx}`, {
    method: "DELETE",
    headers: headersSheep,
  }).then(() => showSheep(currentSheep.id));
};
