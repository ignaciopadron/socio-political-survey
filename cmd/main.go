package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Person struct {
	Name     string `json:"name"`     // «Nicolás Maquiavelo»
	Short    string `json:"short"`    // 1‑2 líneas
	Full     string `json:"full"`     // texto largo (biografía)
	ImageURL string `json:"imageUrl"` // puede quedar vacío ""
}

// Affirmation (MODIFICADO: Quitar json:"-" de Type)
type Affirmation struct {
	Text string `json:"text"`
	Type string `json:"type"`
	Axis string `json:"-"`
}

// QuestionPair (sin cambios)
type QuestionPair struct {
	ID            string      `json:"id"`
	Axis          string      `json:"axis"`
	Affirmation1  Affirmation `json:"affirmation1"`
	Affirmation2  Affirmation `json:"affirmation2"`
	originalType1 string      `json:"-"`
	originalType2 string      `json:"-"`
}

// UserChoice (MODIFICADO)
type UserChoice struct {
	QuestionID string `json:"questionId"`
	ChosenType string `json:"chosenType"` // Añadir este campo
}

// Result representa la puntuación final y el perfil.
// *** MODIFICADO: Thinker y Politician ahora son slices de strings ***
type Result struct {
	ScoreRI     float64  `json:"scoreRI"`     // Puntuación Eje Realismo(0) <-> Idealismo(1)
	ScoreSG     float64  `json:"scoreSG"`     // Puntuación Eje Soberanismo(0) <-> Globalismo(1)
	Profile     string   `json:"profile"`     // Nombre del perfil (Ej: "Realista-Soberanista")
	Description string   `json:"description"` // Descripción del perfil
	Thinkers    []Person `json:"thinkers"`    // Pensadores asociados (plural y slice)
	Politicians []Person `json:"politicians"` // Políticos asociados (plural y slice)
}

// Nueva struct para la respuesta de /api/categories
type CategoryData struct {
	Profile     string   `json:"profile"`
	Description string   `json:"description"`
	Thinkers    []Person `json:"thinkers"`
	Politicians []Person `json:"politicians"`
}

// questionStore (sin cambios, asumiendo que ya tiene las 14 preguntas)
var questionStore = map[string]QuestionPair{
	"q1": {
		ID:            "q1",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "Las consideraciones morales no deben anteponerse al interés nacional en la política exterior.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "La política exterior debe guiarse por principios morales, incluso si a veces eso va en contra del interés nacional.", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	"q2": {
		ID:            "q2",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "El mundo está regido por la competencia entre naciones; los conflictos son inevitables.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "La cooperación y la confianza entre las naciones pueden prevenir conflictos.", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	"q3": {
		ID:            "q3",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "Las leyes y organismos internacionales importan poco si contradicen los intereses de las grandes potencias.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "Las instituciones globales (ONU, etc.) y el derecho internacional son fundamentales para la paz.", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	"q4": {
		ID:            "q4",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "Un sistema político es legítimo si proporciona bienestar y seguridad al pueblo, aunque no sea democrático; en política, el fin puede justificar los medios.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "Un sistema político es legítimo si respeta procedimientos democráticos, incluso si sus resultados no son óptimos; los medios importan tanto como los fines", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	"q5": {
		ID:            "q5",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "Las alianzas internacionales solo duran mientras sirvan al propio interés.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "Las alianzas deben basarse en confianza y valores compartidos, manteniéndose firmes.", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	"q6": {
		ID:            "q6",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "La política mundial siempre será una lucha de poder; es ilusorio pensar que habrá progreso moral.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "La humanidad puede avanzar hacia un orden internacional más justo y pacífico.", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	"q7": {
		ID:            "q7",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "A veces es necesario usar la fuerza militar de forma preventiva para proteger intereses nacionales.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "El uso de la fuerza solo se justifica como último recurso y con legitimidad internacional.", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	// --- Eje Soberanismo vs Globalismo (S/G) ---
	"q8": {
		ID:            "q8",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Ningún país u organismo internacional debe intervenir en los asuntos internos de otro Estado sin su consentimiento, incluso si considera que su gobierno es autoritario o vulnera principios democráticos y derechos fundamentales", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "En ciertos casos, es legítimo que la comunidad internacional intervenga en otro país si sus acciones comprometen la estabilidad regional, internacional o violan derechos humanos.", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
	"q9": {
		ID:            "q9",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Un gobierno mundial o la cesión significativa de soberanía a entidades supranacionales pondría en riesgo la autonomía de las naciones y debería evitarse.", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "Deberíamos aspirar a instituciones globales más fuertes, incluso a alguna forma de autoridad mundial, para enfrentar desafíos que ningún país puede resolver solo.", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
	"q10": {
		ID:            "q10",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Es importante preservar las tradiciones y la identidad nacional propias frente a influencias externas globales.", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "Deberíamos fomentar una identidad más cosmopolita, abiertos a adoptar valores culturales universales y aprender de otras sociedades.", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
	"q11": {
		ID:            "q11",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Cada país debería poder proteger su economía e industria, aunque eso implique salir de acuerdos internacionales o limitar el libre comercio si es necesario", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "Se debe promover el libre comercio y la integración económica global porque benefician a largo plazo.", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
	"q12": {
		ID:            "q12",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Un Estado soberano debe controlar estrictamente sus fronteras y decidir quién entra, sin presiones externas.", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "Facilitar la libre circulación de personas y la inmigración enriquece a los países, deberíamos facilitar lalibre circulación de personas", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
	"q13": {
		ID:            "q13",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Muchos acuerdos internacionales (clima, comercio, salud, etc) limitan injustamente la capacidad de un país de actuar en su propio beneficio.", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "Desafíos globales como el cambio climático o las pandemias exigen respuestas coordinadas a nivel mundial, aunque eso limite en parte la autonomía nacional.", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
	"q14": {
		ID:            "q14",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Organismos internacionales como la ONU o la UE no deben imponer decisiones sobre un gobierno nacional electo.", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "Para lograr un orden global justo, las instituciones internacionales legítimas deberían tener mayor autoridad para hacer cumplir acuerdos y resolver problemas comunes.", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
}

// Mutex y Rand (sin cambios)
var mu sync.Mutex
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// getRandomizedQuestions (sin cambios)
func getRandomizedQuestions() []QuestionPair {
	mu.Lock()
	defer mu.Unlock()

	questions := make([]QuestionPair, 0, len(questionStore))
	for _, q := range questionStore {
		questions = append(questions, q)
	}

	seededRand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	randomizedCopy := make([]QuestionPair, len(questions))
	for i, q := range questions {
		newQ := q
		if seededRand.Intn(2) == 0 {
			newQ.Affirmation1, newQ.Affirmation2 = q.Affirmation2, q.Affirmation1
		}
		randomizedCopy[i] = newQ
	}

	return randomizedCopy
}

// questionsHandler (sin cambios)
func questionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	randomizedQuestions := getRandomizedQuestions()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Añadir CORS si es necesario para desarrollo
	json.NewEncoder(w).Encode(randomizedQuestions)
}

// --- NUEVA FUNCIÓN: Obtener datos de una categoría ---
func getCategoryData(profileName string) CategoryData {
	var description string
	var thinkers []Person
	var politicians []Person

	// Reutilizamos la lógica de asignación del submitHandler
	// (Esta sección podría refactorizarse aún más si los datos fueran externos)
	switch profileName {
	case "Realista-Soberanista":
		description = "Este perfil enfatiza una visión pragmática y nacional de la política..."
		thinkers = []Person{
			{Name: "Nicolás Maquiavelo", Short: `(1469-1527)...`, Full: `Nicolás Maquiavelo (1469-1527)...`, ImageURL: "/static/img/maquiavelo.jpg"},
			{Name: "John Mearsheimer", Short: `Es un influyente...`, Full: `John Mearsheimer...`, ImageURL: "/static/img/mearsheimer.jpg"},
		}
		politicians = []Person{
			{Name: "Xi Jinping", Short: `Líder de la RPC...`, Full: `Xi Jinping (1953-)...`, ImageURL: "/static/img/xijinping.jpg"},
			{Name: "Vladimir Putin", Short: `Presidente de Rusia...`, Full: `Vladimir Putin (1952-)...`, ImageURL: "/static/img/putin.jpg"},
		}
	case "Realista-Globalista":
		description = `Quienes se ubican en el cuadrante realista-globalista...`
		thinkers = []Person{
			{Name: "Zbigniew Brzezinski", Short: `Diplomático y politólogo...`, Full: `Zbigniew Brzezinski...`, ImageURL: "/static/img/brzezinski.jpg"},
			{Name: "Nicholas Spykman", Short: `Geoestratega neerlandés...`, Full: `Nicholas Spykman...`, ImageURL: "/static/img/spykman.jpg"},
		}
		politicians = []Person{
			{Name: "Deng Xiaoping", Short: `Líder supremo de China...`, Full: `Deng XiaoPing...`, ImageURL: "/static/img/xiaoping.jpg"},
			{Name: "Henry Kissinger", Short: `Diplomático y Secretario...`, Full: `Henry Kissinger...`, ImageURL: "/static/img/kissinger.jpg"},
		}
	case "Idealista-Soberanista":
		description = `Este perfil combina la defensa de la soberanía nacional...`
		thinkers = []Person{
			{Name: "Giuseppe Mazzini", Short: `Político y activista...`, Full: `Giuseppe Mazzini...`, ImageURL: "/static/img/mazzini.jpg"},
			{Name: "Mahatma Gandhi", Short: `Líder del movimiento...`, Full: `Mahatma Gandhi...`, ImageURL: "/static/img/gandhi.jpg"},
		}
		politicians = []Person{
			{Name: "Charles de Gaulle", Short: `Militar y estadista...`, Full: `Charles de Gaulle...`, ImageURL: "/static/img/degaulle.jpg"},
			{Name: "Simón Bolívar", Short: `Líder militar y político...`, Full: `Simón Bolívar...`, ImageURL: "/static/img/bolivar.png"},
		}
	case "Idealista-Globalista":
		description = `Este cuadrante representa una postura que apuesta por principios universales...`
		thinkers = []Person{
			{Name: "Immanuel Kant", Short: `Filósofo alemán...`, Full: `Immanuel Kant...`, ImageURL: "/static/img/kant.webp"},
			{Name: "John Rawls", Short: `Filósofo político...`, Full: `John Rawls...`, ImageURL: "/static/img/rawls.webp"},
		}
		politicians = []Person{
			{Name: "George Soros", Short: `Inversor y filántropo...`, Full: `George Soros...`, ImageURL: "/static/img/soros.jpg"},
			{Name: "Barack Obama", Short: `44.º presidente...`, Full: `Barack Obama...`, ImageURL: "/static/img/obama.jpg"},
		}
		// Añadir default o manejo de error si profileName no es válido
	}

	return CategoryData{
		Profile:     profileName,
		Description: description,
		Thinkers:    thinkers,
		Politicians: politicians,
	}
}

// --- FIN NUEVA FUNCIÓN ---

// Handler para recibir las respuestas y calcular el resultado
func submitHandler(w http.ResponseWriter, r *http.Request) {
	// Permitir CORS si el frontend está en un origen diferente
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Manejar preflight request de CORS
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var userChoices []UserChoice
	err := json.NewDecoder(r.Body).Decode(&userChoices)
	if err != nil {
		http.Error(w, "Error al decodificar las respuestas: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Calcular puntuaciones
	var scoreR, scoreI, scoreS, scoreG int
	var totalRI, totalSG int

	for _, choice := range userChoices {
		chosenType := choice.ChosenType

		// Obtener pregunta original solo para saber el eje (si es necesario para error)
		originalQuestion, ok := questionStore[choice.QuestionID]
		if !ok {
			fmt.Printf("Advertencia: Se recibió respuesta para ID de pregunta desconocido: %s\n", choice.QuestionID)
			continue
		}

		// --- LÓGICA DE PUNTUACIÓN Y TOTALES CORREGIDA ---
		switch chosenType {
		case "R":
			scoreR++
			totalRI++ // Incrementar total aquí
		case "I":
			scoreI++
			totalRI++ // Incrementar total aquí
		case "S":
			scoreS++
			totalSG++ // Incrementar total aquí
		case "G":
			scoreG++
			totalSG++ // Incrementar total aquí
		default:
			// Ya no necesitamos revertir totales, simplemente registramos el error
			fmt.Printf("Advertencia: Se recibió tipo inválido ('%s') para la pregunta %s (Eje: %s)\n", chosenType, choice.QuestionID, originalQuestion.Axis)
		}
		// --- FIN CORRECCIÓN ---
	}

	// Normalizar puntuaciones (lógica sin cambios)
	var finalScoreRI float64
	if totalRI > 0 {
		// Realismo (0) <-> Idealismo (1) -> Más I, más cerca de 1
		finalScoreRI = float64(scoreI) / float64(totalRI)
	}

	var finalScoreSG float64
	if totalSG > 0 {
		// Soberanismo (0) <-> Globalismo (1) -> Más G, más cerca de 1
		finalScoreSG = float64(scoreG) / float64(totalSG)
	}

	// Determinar perfil
	profile := ""
	threshold := 0.5
	if finalScoreRI < threshold && finalScoreSG < threshold {
		profile = "Realista-Soberanista"
	} else if finalScoreRI < threshold && finalScoreSG >= threshold {
		profile = "Realista-Globalista"
	} else if finalScoreRI >= threshold && finalScoreSG < threshold {
		profile = "Idealista-Soberanista"
	} else {
		profile = "Idealista-Globalista"
	}

	// Obtener los datos para el perfil calculado usando la nueva función
	categoryResultData := getCategoryData(profile)

	// Crear el objeto resultado final
	result := Result{
		ScoreRI:     finalScoreRI,
		ScoreSG:     finalScoreSG,
		Profile:     profile,                        // Usamos el nombre calculado
		Description: categoryResultData.Description, // Obtenido de la función
		Thinkers:    categoryResultData.Thinkers,    // Obtenido de la función
		Politicians: categoryResultData.Politicians, // Obtenido de la función
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// --- NUEVO HANDLER para /api/categories ---
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	profiles := []string{
		"Realista-Soberanista",
		"Realista-Globalista",
		"Idealista-Soberanista",
		"Idealista-Globalista",
	}

	allData := make([]CategoryData, 0, len(profiles))
	for _, profileName := range profiles {
		allData = append(allData, getCategoryData(profileName))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // CORS si es necesario
	json.NewEncoder(w).Encode(allData)                 // Devolver directamente el slice
}

// --- FIN NUEVO HANDLER ---

// main (MODIFICADO: añadir nueva ruta)
func main() {
	// Servir archivos estáticos desde el directorio ../static bajo la ruta /static/
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Servir index.html para la ruta raíz
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Asegurarse de que solo se sirva para la ruta exacta "/"
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "templates/index.html")
	})

	http.HandleFunc("/api/questions", questionsHandler)
	http.HandleFunc("/api/submit", submitHandler)
	http.HandleFunc("/api/categories", categoriesHandler) // <-- AÑADIR NUEVA RUTA

	port := "8080"
	fmt.Printf("Servidor iniciado en http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
