// Elementos del DOM
const introductionDiv = document.getElementById('introduction');
const surveyDiv = document.getElementById('survey');
const resultsDiv = document.getElementById('results');
const loadingIndicator = document.getElementById('loadingIndicator');
const categoriesSection = document.getElementById('categoriesSection');
const categoriesContent = document.getElementById('categoriesContent');
const mainContainer = document.getElementById('mainContainer');

const startButton = document.getElementById('startButton');
const questionCounter = document.getElementById('questionCounter');
const option1Button = document.getElementById('option1');
const option2Button = document.getElementById('option2');
const affirmationButtons = document.querySelectorAll('.affirmation-button');

const resultMarker = document.getElementById('resultMarker');
const profileName = document.getElementById('profileName');
const profileDescription = document.getElementById('profileDescription');
const restartButton = document.getElementById('restartButton');

// Enlaces de Navegación
const brandLink = document.getElementById('brandLink');
const surveyLink = document.getElementById('surveyLink');
const categoriesLink = document.getElementById('categoriesLink');

// Estado de la aplicación
let questions = [];
let currentQuestionIndex = 0;
let userAnswers = [];
let currentResult = null;
let allCategoriesData = null;

// --- Función auxiliar para crear la tarjeta HTML (MOVIDA AQUÍ ARRIBA) ---
function cardHTML(person, index, kind, categoryProfile) {
  const img = person.imageUrl || "https://via.placeholder.com/300x200?text=Sin+imagen";
  const shortDescription = person.short || person.full || 'Descripción no disponible';
  return `
    <div class="col-md-6 mb-4">
      <div class="card shadow-sm person-card h-100"
           data-kind="${kind}" 
           data-index="${index}"
           data-category="${categoryProfile}"> 
        <img src="${img}" class="card-img-top" alt="Foto de ${person.name || 'Persona'}">
        <div class="card-body d-flex flex-column">
          <h5 class="card-title">${person.name || 'Nombre no disponible'}</h5>
          <p class="card-text truncated mb-3">${shortDescription}</p>
          <button class="btn btn-sm btn-outline-primary mt-auto align-self-start">Ver más</button>
        </div>
      </div>
    </div>`;
}
// --- FIN FUNCIÓN MOVIDA ---

// Función para cargar las preguntas desde el backend
async function fetchQuestions() {
    try {
        const response = await fetch('/api/questions');
        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }
        questions = await response.json();
    } catch (error) {
        console.error("Error al cargar las preguntas:", error);
        alert("No se pudieron cargar las preguntas. Inténtalo de nuevo más tarde.");
    }
}

// Muestra la pregunta actual
function displayQuestion() {
    if (currentQuestionIndex >= questions.length) {
        submitAnswers();
        return;
    }

    const currentQuestion = questions[currentQuestionIndex];
    questionCounter.textContent = `Pregunta ${currentQuestionIndex + 1} de ${questions.length}`;
    option1Button.textContent = currentQuestion.affirmation1.text;
    option2Button.textContent = currentQuestion.affirmation2.text;

    // --- AÑADIDO: Guardar el tipo real en data-type --- 
    option1Button.dataset.type = currentQuestion.affirmation1.type;
    option2Button.dataset.type = currentQuestion.affirmation2.type;
    // --- FIN AÑADIDO ---

    introductionDiv.classList.add('d-none');
    resultsDiv.classList.add('d-none');
    loadingIndicator.classList.add('d-none');
    surveyDiv.classList.remove('d-none');
}

// Maneja la selección de una opción
function handleOptionSelect(event) {
    const chosenButton = event.target.closest('.affirmation-button');
    if (!chosenButton) return;

    // --- MODIFICADO: Obtener tipo en lugar de 'data-choice' --- 
    // const choice = chosenButton.getAttribute('data-choice'); // Ya no se usa
    const chosenType = chosenButton.dataset.type; // Recuperar el tipo 'R','I','S','G'
    if (!chosenType) { // Comprobación de seguridad
        console.error("Error: No se encontró el data-type en el botón pulsado.");
        return;
    }
    // --- FIN MODIFICADO ---

    const currentQuestionId = questions[currentQuestionIndex].id;

    // --- MODIFICADO: Enviar chosenType en lugar de chosen ---
    userAnswers.push({
        questionId: currentQuestionId,
        // chosen: choice, // Ya no se envía
        chosenType: chosenType // Enviar el tipo directamente
    });
    // --- FIN MODIFICADO ---

    currentQuestionIndex++;
    displayQuestion();
}

// Envía las respuestas al backend y muestra los resultados
async function submitAnswers() {
    surveyDiv.classList.add('d-none');
    loadingIndicator.classList.remove('d-none');

    try {
        const response = await fetch('/api/submit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userAnswers),
        });

        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }

        const result = await response.json();
        currentResult = result; // Guardar el resultado globalmente
        displayResults(result);

    } catch (error) {
        console.error("Error al enviar las respuestas:", error);
        alert("Hubo un problema al calcular tus resultados. Inténtalo de nuevo.");
        loadingIndicator.classList.add('d-none');
        showIntroduction();
    }
}

// Muestra los resultados en la UI
function displayResults(result) {
    console.log("Mostrando resultados para:", result);
    // Eje X: Realismo (0%) a Idealismo (100%)
    const markerLeft = result.scoreRI * 100;
    // Eje Y: Soberanismo (top, 0%) a Globalismo (bottom, 100%)
    const markerTop = result.scoreSG * 100; // MODIFICADO

    resultMarker.style.left = `${markerLeft}%`;
    resultMarker.style.top = `${markerTop}%`;

    profileName.textContent = result.profile || 'N/A';
    profileDescription.textContent = result.description || 'Sin descripción.';

    // Mostrar los pensadores y políticos en cuadros separados con tarjeta (card)
    const thinkersContainer = document.getElementById('thinkersContainer');
    const politiciansContainer = document.getElementById('politiciansContainer');

    thinkersContainer.innerHTML = '';
    politiciansContainer.innerHTML = '';

    // Generar tarjetas para pensadores (CORREGIDO: usar result.profile)
    if (result.thinkers && Array.isArray(result.thinkers) && result.thinkers.length > 0) {
        thinkersContainer.innerHTML = result.thinkers
          .map((p, i) => cardHTML(p, i, "thinker", result.profile)).join(""); // Usar profile (minúscula)
    } else {
        thinkersContainer.innerHTML = '<p class="col-12">No se encontraron pensadores cercanos.</p>';
    }

    // Generar tarjetas para políticos (CORREGIDO: usar result.profile)
    if (result.politicians && Array.isArray(result.politicians) && result.politicians.length > 0) {
        politiciansContainer.innerHTML = result.politicians
          .map((p, i) => cardHTML(p, i, "politician", result.profile)).join(""); // Usar profile (minúscula)
    } else {
        politiciansContainer.innerHTML = '<p class="col-12">No se encontraron políticos cercanos.</p>';
    }

    hideAllSectionsExcept(resultsDiv);
}

// Resetea el estado para empezar de nuevo
function resetSurvey() {
    questions = [];
    currentQuestionIndex = 0;
    userAnswers = [];
    resultsDiv.classList.add('d-none');
    surveyDiv.classList.add('d-none');
    loadingIndicator.classList.add('d-none');
    introductionDiv.classList.remove('d-none');
}

// Función para mostrar la pantalla inicial
function showIntroduction() {
    resetSurvey();
    hideAllSectionsExcept(introductionDiv);
    setActiveLink(surveyLink);
}

// Función para marcar el enlace activo en la navbar
function setActiveLink(activeLinkElement) {
    [surveyLink, categoriesLink].forEach(link => {
        if (link === activeLinkElement) {
            link.classList.add('active');
            link.setAttribute('aria-current', 'page');
        } else {
            link.classList.remove('active');
            link.removeAttribute('aria-current');
        }
    });
}

// --- Event Listeners ---
startButton.addEventListener('click', async () => {
    hideAllSectionsExcept(loadingIndicator);
    await fetchQuestions();
    if (questions.length > 0) {
        currentQuestionIndex = 0;
        userAnswers = [];
        hideAllSectionsExcept(surveyDiv);
        displayQuestion();
    } else {
        showIntroduction();
    }
});

// Delegación de eventos para los botones de opción
surveyDiv.addEventListener('click', handleOptionSelect);

restartButton.addEventListener('click', showIntroduction);

// Enlaces de Navegación
brandLink.addEventListener('click', (e) => { e.preventDefault(); showIntroduction(); });
surveyLink.addEventListener('click', (e) => { e.preventDefault(); showIntroduction(); });
categoriesLink.addEventListener('click', (e) => { e.preventDefault(); showCategories(); });

// Listener para mostrar el modal (MODIFICADO: Lógica de búsqueda simplificada)
mainContainer.addEventListener("click", e => { 
  const card = e.target.closest(".person-card");
  if (!card) return; 

  const kind = card.dataset.kind;
  const index = Number(card.dataset.index);
  const categoryProfile = card.dataset.category; // Ahora debería ser consistente

  if (!categoryProfile || categoryProfile === 'undefined') { // Chequeo más robusto
      console.error("Error: Atributo data-category inválido o no encontrado en la tarjeta.", categoryProfile);
      return;
  }

  let sourceData = null;
  let category = null; 

  // Priorizar allCategoriesData si existe y contiene la categoría
  if (allCategoriesData) {
      category = allCategoriesData.find(cat => cat.profile === categoryProfile); 
      if (category) {
          sourceData = category; // Usamos el objeto categoría encontrado
      }
  } 
  
  // Si no se encontró en allCategoriesData O allCategoriesData no existe,
  // intentar con currentResult (solo si su perfil coincide)
  if (!sourceData && currentResult && currentResult.profile === categoryProfile) {
       sourceData = currentResult; // Usamos el resultado actual
       console.log("Usando currentResult como fuente de datos."); // DEBUG
  }

  if (!sourceData) {
      console.error("Error: No se encontraron datos fuente para la categoría:", categoryProfile, "allCategoriesData:", allCategoriesData, "currentResult:", currentResult);
      // Intentar cargar categorías si no existen?
      // if (!allCategoriesData) { fetchCategoriesData(); } // Podría ser una opción, pero añade complejidad
      return; 
  }
  
  let person = null;
  // Buscar persona dentro de la fuente de datos encontrada
  if (kind === "thinker" && sourceData.thinkers && index >= 0 && index < sourceData.thinkers.length) {
      person = sourceData.thinkers[index];
  } else if (kind === "politician" && sourceData.politicians && index >= 0 && index < sourceData.politicians.length) {
      person = sourceData.politicians[index];
  }

  if (!person) {
    console.error("No se pudo encontrar la persona con kind:", kind, "e index:", index, "en categoría", categoryProfile, "Datos fuente:", sourceData);
    return; 
  }

  // --- Rellenar modal (sin cambios) --- 
  const modalElement = document.getElementById('personModal');
  if (!modalElement) return;
  const modalLabel = modalElement.querySelector("#personModalLabel");
  const modalImg = modalElement.querySelector("#personModalImg");
  const modalText = modalElement.querySelector("#personModalText");
  if (modalLabel) modalLabel.textContent = person.name || 'Nombre no disponible';
  if (modalImg) {
      modalImg.src = person.imageUrl || "";
      modalImg.alt = `Foto de ${person.name || 'Persona'}`;
      modalImg.style.display = person.imageUrl ? 'inline-block' : 'none';
  }
  if (modalText) modalText.textContent = person.full || person.short || 'Información no disponible';

  // --- Abrir modal (sin cambios) --- 
  try {
      const modal = bootstrap.Modal.getOrCreateInstance(modalElement);
      modal.show();
  } catch (error) {
      console.error("Error al mostrar el modal de Bootstrap:", error);
  }
});

// --- NUEVAS FUNCIONES PARA LA SECCIÓN DE CATEGORÍAS ---

// Oculta todas las secciones principales excepto la indicada
function hideAllSectionsExcept(sectionToShow) {
    [introductionDiv, surveyDiv, resultsDiv, categoriesSection, loadingIndicator].forEach(div => {
        if (div !== sectionToShow) {
            div.classList.add('d-none');
        } else {
            div.classList.remove('d-none');
        }
    });
}

// Muestra la sección de categorías y carga los datos si es necesario
async function showCategories() {
    hideAllSectionsExcept(categoriesSection);
    setActiveLink(categoriesLink);

    if (!allCategoriesData) {
        categoriesContent.innerHTML = `<div class="text-center">
           <div class="spinner-border" role="status">
             <span class="visually-hidden">Cargando categorías...</span>
           </div>
        </div>`;
        await fetchCategoriesData();
    } else {
        displayCategories(allCategoriesData);
    }
}

// Obtiene los datos de todas las categorías del backend
async function fetchCategoriesData() {
    try {
        const response = await fetch('/api/categories');
        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }
        const data = await response.json();
        allCategoriesData = data;
        displayCategories(allCategoriesData);
    } catch (error) {
        console.error("Error al cargar las categorías:", error);
        categoriesContent.innerHTML = '<p class="text-center text-danger">Error al cargar las categorías. Inténtalo de nuevo más tarde.</p>';
    }
}

// Muestra los datos de todas las categorías en la UI
function displayCategories(categories) {
    console.log("Iniciando displayCategories con:", categories); // DEBUG
    categoriesContent.innerHTML = '';

    if (!categories || !Array.isArray(categories) || categories.length === 0) {
        console.log("No hay categorías válidas para mostrar."); // DEBUG
        categoriesContent.innerHTML = '<p class="text-center">No hay categorías para mostrar.</p>';
        return;
    }

    categories.forEach((category, categoryIndex) => {
        // --- CORREGIDO: Usar claves JSON en minúscula --- 
        console.log(`Procesando categoría ${categoryIndex}:`, category.profile); 
        try {
            const categoryElement = document.createElement('div');
            categoryElement.classList.add('mb-5');

            let thinkersHtml = '<p>No hay pensadores asociados.</p>';
            if (category.thinkers && category.thinkers.length > 0) { // Usar thinkers (minúscula)
                thinkersHtml = category.thinkers
                    .map((p, i) => cardHTML(p, i, 'thinker', category.profile)) // Usar profile (minúscula)
                    .join('');
            }

            let politiciansHtml = '<p>No hay políticos asociados.</p>';
            if (category.politicians && category.politicians.length > 0) { // Usar politicians (minúscula)
                politiciansHtml = category.politicians
                    .map((p, i) => cardHTML(p, i, 'politician', category.profile)) // Usar profile (minúscula)
                    .join('');
            }

            categoryElement.innerHTML = `
                <h3 class="mb-3">${category.profile}</h3> {/* Usar profile */} 
                <p>${category.description}</p> {/* Usar description */} 
                <h4 class="mt-4">Pensadores Cercanos</h4>
                <div class="row">${thinkersHtml}</div>
                <h4 class="mt-4">Políticos Cercanos</h4>
                <div class="row">${politiciansHtml}</div>
                <hr class="my-5">
            `;
            categoriesContent.appendChild(categoryElement);
            console.log(`Categoría ${category.profile} renderizada correctamente.`); // Usar profile (minúscula)
        } catch (error) {
             console.error(`Error al procesar la categoría ${category.profile}:`, error); // Usar profile (minúscula)
             const errorElement = document.createElement('div');
             errorElement.innerHTML = `<p class="text-danger">Error al mostrar la categoría ${category.profile}.</p>`; // Usar profile
             categoriesContent.appendChild(errorElement);
        }
        // --- FIN CORRECCIÓN --- 
    });
}

// Mostrar la introducción al cargar la página
showIntroduction();
