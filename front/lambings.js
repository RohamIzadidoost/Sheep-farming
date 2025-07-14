const tokenL = localStorage.getItem("token");
if (!tokenL) window.location.href = "index.html";
const headersL = { Authorization: `Bearer ${tokenL}` };
const tbodyL = document.querySelector("#lambTable tbody");

function loadLambings() {
  const params = new URLSearchParams();
  const f = document.getElementById("fromL").value;
  if (f) params.append("from", f);
  const t = document.getElementById("toL").value;
  if (t) params.append("to", t);
  const qs = params.toString() ? `?${params.toString()}` : "";
  fetch(`${API_BASE}/lambings${qs}`, { headers: headersL })
    .then((r) => r.json())
    .then((list) => {
      tbodyL.innerHTML = "";
      list.forEach((l) => {
        const tr = document.createElement("tr");
        tr.innerHTML = `<td>${l.sheepEar || ""}</td><td>${l.date}</td><td>${l.numBorn}</td><td>${(l.sexes || []).join(",")}</td><td>${l.numDead}</td>`;
        tbodyL.appendChild(tr);
      });
    });
}

document.getElementById("lambFilter").addEventListener("submit", (e) => {
  e.preventDefault();
  loadLambings();
});
loadLambings();
