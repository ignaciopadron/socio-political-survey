package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Affirmation (sin cambios)
type Affirmation struct {
	Text string `json:"text"`
	Type string `json:"-"`
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

// UserChoice (sin cambios)
type UserChoice struct {
	QuestionID string `json:"questionId"`
	Chosen     string `json:"chosen"`
}

// Result representa la puntuación final y el perfil.
// *** MODIFICADO: Thinker y Politician ahora son slices de strings ***
type Result struct {
	ScoreRI     float64  `json:"scoreRI"`     // Puntuación Eje Realismo(0) <-> Idealismo(1)
	ScoreSG     float64  `json:"scoreSG"`     // Puntuación Eje Soberanismo(0) <-> Globalismo(1)
	Profile     string   `json:"profile"`     // Nombre del perfil (Ej: "Realista-Soberanista")
	Description string   `json:"description"` // Descripción del perfil
	Thinkers    []string `json:"thinkers"`    // Pensadores asociados (plural y slice)
	Politicians []string `json:"politicians"` // Políticos asociados (plural y slice)
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
		Affirmation2:  Affirmation{Text: "Facilitar la libre circulación de personas y la inmigración enriquece a los países, deberíamos facilitar el tránsito de flujos migratorios", Type: "G", Axis: "SG"},
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

// Handler para recibir las respuestas y calcular el resultado
// *** MODIFICADO: Usa slices para thinkers y politicians ***
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

	// Calcular puntuaciones (lógica sin cambios)
	var scoreR, scoreI, scoreS, scoreG int
	var totalRI, totalSG int

	for _, choice := range userChoices {
		originalQuestion, ok := questionStore[choice.QuestionID]
		if !ok {
			fmt.Printf("Advertencia: Se recibió respuesta para ID de pregunta desconocido: %s\n", choice.QuestionID)
			continue
		}

		var chosenType string
		if choice.Chosen == "affirmation1" {
			// Necesitamos buscar la afirmación original que está ahora en affirmation1
			if originalQuestion.Affirmation1.Text == questionStore[choice.QuestionID].Affirmation1.Text {
				chosenType = questionStore[choice.QuestionID].originalType1
			} else {
				chosenType = questionStore[choice.QuestionID].originalType2
			}
		} else if choice.Chosen == "affirmation2" {
			// Necesitamos buscar la afirmación original que está ahora en affirmation2
			if originalQuestion.Affirmation2.Text == questionStore[choice.QuestionID].Affirmation2.Text {
				chosenType = questionStore[choice.QuestionID].originalType2
			} else {
				chosenType = questionStore[choice.QuestionID].originalType1
			}
		} else {
			fmt.Printf("Advertencia: Elección inválida ('%s') para la pregunta %s\n", choice.Chosen, choice.QuestionID)
			continue
		}

		switch chosenType {
		case "R":
			scoreR++
			totalRI++
		case "I":
			scoreI++
			totalRI++
		case "S":
			scoreS++
			totalSG++
		case "G":
			scoreG++
			totalSG++
		}
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
	description := ""
	// *** MODIFICADO: Declarar slices para thinkers y politicians ***
	var thinkers []string
	var politicians []string

	threshold := 0.5

	if finalScoreRI < threshold && finalScoreSG < threshold { // Realista-Soberanista
		profile = "Realista-Soberanista"
		description = "Este perfil enfatiza una visión pragmática y nacional de la política. El realista-soberanista valora ante todo la soberanía del Estado y la búsqueda del poder e interés nacional dentro de un sistema internacional anárquico. Considera que el orden mundial más estable es el basado en Estados fuertes, soberanos y en equilibrio de poder (sistema Westfaliano), desconfiando de esquemas globalistas o ideales universalistas que puedan limitar la independencia nacional. En este cuadrante se prioriza la seguridad, el orden interno y la autonomía del país, asumiendo que las relaciones internacionales son un juego de suma cero donde cada nación vela por sí misma."
		// *** MODIFICADO: Asignar slice con dos elementos ***
		thinkers = []string{
			"Nicolás Maquiavelo (1469-1527) fue un diplomático florentino del Renacimiento y autor de El Príncipe y Los discursos sobre la primera década de Tito Livio. Con estas obras inauguró el realismo político moderno: estudia la política tal como es —lucha por el poder y la seguridad del Estado— y no como “debería” ser según ideales morales. Defiende que el gobernante debe usar cualquier medio —incluso la fuerza o el engaño— si con ello protege la estabilidad y los intereses del Estado. Separa ética privada y razón de Estado, núcleo del realismo político. Escribe para príncipes que necesitan conservar la autonomía y el poder interno frente a potencias externas (el Papado, el Sacro Imperio) y rivales internos. La prioridad absoluta es la soberanía y seguridad del propio Estado; rechaza toda tutela externa sobre las decisiones del gobernante, anticipando la lógica westfaliana de Estados soberanos.",
			"John Mearsheimer. Es un influyente politólogo estadounidense y destacado académico de relaciones internacionales, conocido principalmente por desarrollar la teoría del 'Realismo Ofensivo'. Sostiene que las grandes potencias buscan maximizar su poder y aspiran a la hegemonía para garantizar su supervivencia en el sistema internacional anárquico. Se sitúa como Realista-Soberanista porque su teoría del 'realismo ofensivo' considera a los Estados soberanos como los actores principales, obligados a buscar poder para sobrevivir en un mundo anárquico y competitivo. Por ello, desconfía profundamente de la capacidad de las instituciones globales para alterar esta lógica de poder y prioriza la seguridad y el interés nacional por encima de consideraciones morales o de cooperación idealista, defendiendo así la primacía de la acción estatal autónoma.",
		}
		politicians = []string{
			"Xi Jinping. Presidente de la República Popular China y Secretario General del Partido Comunista. Encaja como Realista porque prioriza el interés nacional, la seguridad y el poder de China por encima de consideraciones morales universales en las relaciones internacionales. Es Soberanista por su férrea defensa de la autonomía nacional, el rechazo a la injerencia externa en asuntos internos y la promoción de un modelo de desarrollo propio sin imposiciones foráneas. Considera la cooperación global principalmente como una herramienta pragmática para avanzar estos objetivos nacionales.",
			"Vladimir Putin. Actual Presidente de la Federación de Rusia, figura política dominante en el país desde principios del siglo XXI. Su largo mandato ha estado marcado por la restauración del poder estatal ruso. Su política prioriza el interés nacional, la seguridad y el poder de Rusia (Realismo), viendo las relaciones internacionales principalmente como una competición estratégica. Al mismo tiempo, defiende férreamente la autoridad del Estado ruso, promueve una identidad nacional fuerte y la soberanía nacional frente a cualquier injerencia externa o limitación por parte de instituciones globales (Soberanismo).",
		}
	} else if finalScoreRI < threshold && finalScoreSG >= threshold { // Realista-Globalista
		profile = "Realista-Globalista"
		description = "Quienes se ubican en el cuadrante realista-globalista comparten la visión de que la política mundial es principalmente una competencia estratégica y que los Estados actúan guiados por interés propio, pero al mismo tiempo reconocen y operan dentro de la interdependencia global. Este enfoque hace hincapié en la naturaleza competitiva del sistema internacional y asume que los estados operan en un entorno sin justicia inherente, donde las normas éticas pueden quedar supeditadas al poder. Un realista-globalista típicamente apoya la cooperación internacional solo de forma pragmática, por ejemplo, mediante alianzas o instituciones, siempre y cuando esto beneficie el equilibrio de poder o los intereses de su propio país. También suele abogar por un liderazgo fuerte de las potencias para mantener la estabilidad global, en lugar de confiar en principios idealistas."
		// *** MODIFICADO: Asignar slice con dos elementos ***
		thinkers = []string{
			"Zbigniew Brzezinski. Fue un influyente diplomático y politólogo polaco-estadounidense, conocido principalmente por ser Consejero de Seguridad Nacional del presidente Jimmy Carter. Su pensamiento estratégico se centró en la geopolítica y tuvo un gran impacto en la política exterior de Estados Unidos durante y después de la Guerra Fría. Se le considera Realista-Globalista porque su análisis priorizaba descarnadamente el poder, la estrategia y el interés nacional estadounidense (Realismo), tal como expuso en 'El Gran Tablero Mundial'. Sin embargo, su campo de acción y visión eran intrínsecamente globales, abogando por gestionar activamente las relaciones internacionales y usar las estructuras de poder a nivel mundial (Globalismo) para asegurar la primacía de EE. UU., en lugar de enfocarse en ideales universales o en un soberanismo defensivo.",
			"Nicholas Spykman. Fue un influyente geoestratega y profesor neerlandés-estadounidense, considerado uno de los padres del realismo geopolítico y célebre por su teoría del 'Rimland' sobre la importancia estratégica de las zonas costeras de Eurasia. Se sitúa como Realista-Globalista porque, si bien su análisis se basa crudamente en el poder, la geografía y el interés nacional (Realismo), defendía que la seguridad de un Estado (como EE.UU.) requería una activa intervención y la formación de alianzas a escala mundial. Argumentaba vehementemente en contra del aislacionismo, insistiendo en que la proyección de poder y el control del equilibrio en zonas clave lejanas como el 'Rimland' eran esenciales para la supervivencia y primacía del Estado (enfoque Globalista estratégico).",
		}
		politicians = []string{
			"Deng XiaoPing. Fue el líder supremo de China tras Mao Zedong, considerado el arquitecto de la política de 'Reforma y Apertura' iniciada en 1978. Utilizó estratégicamente la integración en la economía mundial (globalismo instrumental) como la herramienta más eficaz para alcanzar la modernización y la prosperidad. Su famoso lema 'No importa si el gato es blanco o negro, mientras cace ratones' encapsula su pragmatismo. Dejó de lado la rigidez ideológica maoísta para centrarse en la política de 'Reforma y Apertura' fue una decisión estratégica calculada para aumentar el poder nacional integral de China, no basada en ideales universales de cooperación, sino en la necesidad práctica de modernizar el país. Su enfoque en 'ocultar la fuerza y esperar el momento' también es una clara estrategia realista de acumulación de poder. Tiene una forma pragmática de Globalismo instrumental, al servicio de fines Realistas y Soberanistas a largo plazo.",
			"Henry Kissinger. Estadista y diplomático, fue un influyente diplomático y politólogo germano-estadounidense, que ejerció como Secretario de Estado y Consejero de Seguridad Nacional de EE.UU., moldeando decisivamente la política exterior estadounidense durante la Guerra Fría. Es una figura central asociada a la Realpolitik y la diplomacia pragmática a escala global. Encaja como Globalista (en un sentido instrumental y estratégico) porque operó a escala mundial, utilizando la diplomacia compleja (como la apertura a China o la distensión con la URSS) y las negociaciones globales como herramientas clave para gestionar activamente el sistema internacional y asegurar ventajas estratégicas para Estados Unidos. Su objetivo no era un ideal globalista per se, sino gestionar un orden global desde una perspectiva de poder e interés nacional.",
		}
	} else if finalScoreRI >= threshold && finalScoreSG < threshold { // Idealista-Soberanista
		profile = "Idealista-Soberanista"
		description = "Este perfil combina la defensa de la soberanía nacional con la adhesión a principios e ideales claros. El idealista-soberanista cree en la autodeterminación de los pueblos y en la importancia de que cada nación sea libre para perseguir sus propios ideales. Suele oponerse a cualquier dominación imperial o extranjera sobre una nación, argumentando que la legitimidad política nace de los valores y la voluntad del propio pueblo. Por ejemplo, se enfatiza que ninguna nación merece la libertad si no es capaz de conquistarla por sí misma, reflejando un fuerte compromiso con la independencia nacional y la justicia. A diferencia del realista puro, este perfil sí otorga un peso moral a la política, pero concentrado en el ámbito nacional: la soberanía se valora como un medio para garantizar libertad, dignidad y progreso de la nación conforme a sus propios ideales (ya sean democracia, igualdad, tradición cultural, etc.)."
		// *** MODIFICADO: Asignar slice con dos elementos ***
		thinkers = []string{
			"Giuseppe Mazzini. Pensador del siglo XIX. Político, periodista y activista italiano, una figura clave del Risorgimento (la unificación italiana). Fundador de la organización 'Joven Italia'. Mazzini era un republicano ferviente y creía profundamente en ideales como la libertad, la igualdad, el progreso y la fraternidad entre los pueblos, a menudo con un tinte casi religioso ('Dios y el Pueblo'). Creía que cada nación tenía una misión moral específica que cumplir para el progreso de la humanidad. Su visión no era meramente pragmática, sino basada en principios éticos y políticos elevados. Su objetivo principal era liberar a Italia del dominio extranjero (austríaco, papal, borbónico) y unificarla como una república independiente y soberana. Creía firmemente en el principio de autodeterminación nacional: cada pueblo tenía derecho a gobernarse a sí mismo y a formar su propio Estado-nación libre de injerencias externas. Su idealismo estaba intrínsecamente ligado a la consecución de la soberanía nacional italiana para que esta pudiera cumplir su 'misión'.",
			"Mahatma Gandhi. Fue el líder preeminente del movimiento de independencia de la India contra el dominio británico. Es célebre mundialmente por su filosofía de la desobediencia civil no violenta (Satyagraha), que inspiró movimientos por los derechos civiles en todo el mundo. Gandhi encaja como Idealista porque fundamentó su lucha en principios éticos y espirituales profundos, como la verdad y la no violencia (Ahimsa), creyendo firmemente en la superioridad moral de estos métodos para lograr un cambio político y social justo. A su vez, es Soberanista porque su objetivo central e irrenunciable fue alcanzar la independencia completa (Purna Swaraj) de la India del control imperial británico, defendiendo el derecho inalienable de la nación india a la autodeterminación y al autogobierno integral (Swaraj).",
		}
		politicians = []string{
			"Charles de Gaulle. Militar y estadista francés, líder de la Francia Libre durante la Segunda Guerra Mundial y arquitecto de la Quinta República Francesa, de la que fue el primer presidente. De Gaulle tenía una idea muy particular y elevada de Francia, una visión casi mística de su grandeza, su historia y su papel en el mundo. Creía en la singularidad de la civilización francesa y en la necesidad de preservar sus valores republicanos y su independencia cultural. Su política estaba guiada por este ideal de la 'grandeza' y la misión histórica de Francia. Es quizás uno de los soberanistas más emblemáticos del siglo XX. Defendió a ultranza la independencia nacional de Francia frente a las superpotencias (EE.UU. y la URSS), desarrollando un programa nuclear propio, retirando a Francia del mando militar integrado de la OTAN y vetando la entrada del Reino Unido en la Comunidad Económica Europea por considerarlo un 'caballo de Troya' de los intereses estadounidenses. Su ideal de Francia solo podía realizarse a través de una soberanía política, militar y económica plena.",
			"Simón Bolívar. Bolívar encaja como Idealista-Soberanista porque su lucha emancipadora estaba profundamente inspirada en ideales ilustrados de libertad, igualdad y autogobierno republicano (Idealismo). Estaba convencido de que los pueblos debían autodeterminarse y construir su destino. Rechazó tanto el colonialismo como el dominio extranjero posterior. La independencia era el vehículo necesario para realizar su visión ideal de naciones libres.",
		}
	} else { // Idealista-Globalista
		profile = "Idealista-Globalista"
		description = "Este cuadrante representa una postura que apuesta por principios universales y cooperación a escala global. El idealista-globalista cree que los Estados deben trascender sus diferencias en pos de objetivos comunes de la humanidad, y tiende a apoyar la creación de instituciones internacionales fuertes y incluso estructuras de gobernanza global. Históricamente, esta visión se inspira en ideas como las de Immanuel Kant, quien proponía lograr una “paz perpetua” mediante una federación de repúblicas y la cooperación internacional basada en la razón y la moral.  En este perfil se priorizan valores como los derechos humanos, el derecho internacional, el desarrollo sostenible y la resolución pacífica de conflictos a través del diálogo multinacional. Un idealista-globalista está dispuesto a ceder parte de la soberanía nacional en favor de instituciones o acuerdos globales si con ello se alcanzan bienes mayores (por ejemplo, la paz mundial, la lucha contra el cambio climático o la erradicación de la pobreza)."
		// *** MODIFICADO: Asignar slice con dos elementos ***
		thinkers = []string{
			"Immanuel Kant. fue un filósofo alemán del siglo XVIII, considerado uno de los pensadores más influyentes de la Ilustración y de la filosofía occidental moderna. Se sitúa en la categoría Idealista-Globalista porque desarrolló el idealismo trascendental, una teoría que sostiene que el conocimiento humano se basa en estructuras mentales innatas que dan forma a nuestra experiencia del mundo . Además, en su ensayo La paz perpetua propuso la creación de una federación de estados republicanos como medio para alcanzar una paz duradera a través de la cooperación internacional y el respeto mutuo entre las naciones.",
			"John Rawls. Fue un filósofo político estadounidense, reconocido por su influyente obra Teoría de la justicia (1971), donde propuso el concepto de 'justicia como equidad'. Se sitúa en la categoría Idealista-Globalista porque su teoría se basa en principios éticos universales aplicables a todas las sociedades, enfatizando la equidad y los derechos fundamentales. Además, en su obra El Derecho de los Pueblos (1999), abogó por una comunidad internacional justa, donde los pueblos cooperan bajo normas comunes de justicia y respeto mutuo, promoviendo una visión global de la justicia.",
		}
		politicians = []string{
			"George Soros. Es un inversor y filántropo húngaro-estadounidense conocido por financiar causas progresistas y fundar la Open Society Foundations, dedicada a la promoción de la democracia, los derechos humanos y la gobernanza global. Su inclusión en esta categoría se debe principalmente a la misión de sus Open Society Foundations, que dedican miles de millones a promover globalmente la democracia liberal, los derechos humanos y la sociedad civil, reflejando un idealismo basado en valores universales. Además, es un firme defensor de la cooperación internacional, apoya instituciones supranacionales como la Unión Europea y critica abiertamente los nacionalismos, alineándose claramente con una visión globalista frente a la soberanista.",
			"Barack Obama fue el 44.º presidente de los Estados Unidos, reconocido por su enfoque en la diplomacia multilateral y la promoción de valores democráticos a nivel global. Se sitúa en la categoría Idealista-Globalista debido a su compromiso con el multilateralismo y la cooperación internacional. Durante su mandato, promovió acuerdos como el Acuerdo de París sobre cambio climático y defendió el fortalecimiento de instituciones globales como la ONU, reflejando su creencia en que los desafíos globales requieren soluciones colectivas basadas en principios éticos compartidos.",
		}
	}

	// Crear el objeto resultado
	// *** MODIFICADO: Usar los campos plurales y las slices ***
	result := Result{
		ScoreRI:     finalScoreRI,
		ScoreSG:     finalScoreSG,
		Profile:     profile,
		Description: description,
		Thinkers:    thinkers,    // Asigna la slice
		Politicians: politicians, // Asigna la slice
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// main (sin cambios)
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

	port := "8080"
	fmt.Printf("Servidor iniciado en http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
