// URL del middleware. El proposito en este proyecto es que sea un scanner local

const MIDDLEWARE_URL = "http://localhost:8080/scan";

// Referencias a los elementos del DOM que vamos a leer/escribir. Se buscan una sola vez al cargar el script, no cada vez que se usan.

const targetInput = document.getElementById("target");
const outputBox = document.getElementById("output");
const scanBtn = document.getElementById("scanBtn");

// Conectar el botón con la función que dispara el escaneo.

scanBtn.addEventListener("click", runScan);

async function runScan() {
  const target = targetInput.value;

  // Junta los módulos marcados: busca todos los checkboxes marcados y extrae su atributo "value" (ej. "port_scan").

  const checkedBoxes = document.querySelectorAll("input[type=checkbox]:checked");
  const modules = Array.from(checkedBoxes).map(checkbox => checkbox.value);

  outputBox.textContent = "Escaneando...";

  try {
    const response = await fetch(MIDDLEWARE_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ target, modules })
    });

    const data = await response.json();

    // JSON.stringify con esos 2 argumentos extra formatea el JSON con indentación, para que se lea bien.

    outputBox.textContent = JSON.stringify(data, null, 2);

  } catch (error) {
    outputBox.textContent = "Error: " + error;
  }
}
