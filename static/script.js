// Elementos del DOM
const introductionDiv = document.getElementById('introduction');
const surveyDiv = document.getElementById('survey');
const resultsDiv = document.getElementById('results');
const loadingIndicator = document.getElementById('loadingIndicator');

const startButton = document.getElementById('startButton');
const questionCounter = document.getElementById('questionCounter');
const option1Button = document.getElementById('option1');
const option2Button = document.getElementById('option2');
const affirmationButtons = document.querySelectorAll('.affirmation-button');

const resultMarker = document.getElementById('resultMarker');
const profileName = document.getElementById('profileName');
const profileDescription = document.getElementById('profileDescription');
const restartButton = document.getElementById('restartButton');

// Estado de la aplicación
let questions = [];
let currentQuestionIndex = 0;
let userAnswers = [];

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

    introductionDiv.classList.add('d-none');
    resultsDiv.classList.add('d-none');
    loadingIndicator.classList.add('d-none');
    surveyDiv.classList.remove('d-none');
}

// Maneja la selección de una opción
function handleOptionSelect(event) {
    const chosenButton = event.target.closest('.affirmation-button');
    if (!chosenButton) return;

    const choice = chosenButton.getAttribute('data-choice');
    const currentQuestionId = questions[currentQuestionIndex].id;

    userAnswers.push({
        questionId: currentQuestionId,
        chosen: choice
    });

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

    result.thinkers.forEach((thinker, index) => {
        const card = document.createElement('div');
        card.classList.add('col-md-6', 'mb-3');
        card.innerHTML = `
            <div class="card">
                <img src="https://via.placeholder.com/150" class="card-img-top" alt="Imagen de pensador">
                <div class="card-body">
                    <p class="card-text truncated-text">${thinker}</p>
                </div>
            </div>
        `;
        // Agregar el listener para expandir/retraer el párrafo
        card.querySelector('.truncated-text').addEventListener('click', function() {
            this.classList.toggle('expanded');
        });
        thinkersContainer.appendChild(card);
    });

    result.politicians.forEach((politician, index) => {
        const card = document.createElement('div');
        card.classList.add('col-md-6', 'mb-3');
        card.innerHTML = `
            <div class="card">
                <img src="https://via.placeholder.com/150" class="card-img-top" alt="Imagen de político">
                <div class="card-body">
                    <p class="card-text truncated-text">${politician}</p>
                </div>
            </div>
        `;
        card.querySelector('.truncated-text').addEventListener('click', function() {
            this.classList.toggle('expanded');
        });
        politiciansContainer.appendChild(card);
    });

    loadingIndicator.classList.add('d-none');
    resultsDiv.classList.remove('d-none');
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
    introductionDiv.classList.remove('d-none');
    surveyDiv.classList.add('d-none');
    resultsDiv.classList.add('d-none');
    loadingIndicator.classList.add('d-none');
}

// --- Event Listeners ---
startButton.addEventListener('click', async () => {
    introductionDiv.classList.add('d-none');
    loadingIndicator.classList.remove('d-none');
    await fetchQuestions();
    loadingIndicator.classList.add('d-none');
    if (questions.length > 0) {
        currentQuestionIndex = 0;
        userAnswers = [];
        displayQuestion();
    } else {
        showIntroduction();
    }
});

// Delegación de eventos para los botones de opción
surveyDiv.addEventListener('click', handleOptionSelect);

restartButton.addEventListener('click', resetSurvey);

// Mostrar la introducción al cargar la página
showIntroduction();
