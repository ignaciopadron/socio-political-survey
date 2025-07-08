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
		Affirmation1:  Affirmation{Text: "Las grandes potencias actuarán según sus intereses aunque contradigan el derecho internacional.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "Las instituciones globales (ONU, etc.) y el derecho internacional son fundamentales para la paz.", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	"q4": {
		ID:            "q4",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "Un sistema político o gobierno es legítimo si proporciona prosperidad y seguridad, incluso sin elecciones democráticas.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "Un sistema político es legítimo si respeta procedimientos democráticos, incluso si sus resultados no son óptimos; los medios importan tanto como los fines.", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	"q5": {
		ID:            "q5",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "Las alianzas internacionales solo duran mientras sirvan al propio interés.", Type: "R", Axis: "RI"},
		Affirmation2:  Affirmation{Text: "Las alianzas persisten cuando se basan en la confianza y valores compartidos, incluso ante costes a corto plazo.", Type: "I", Axis: "RI"},
		originalType1: "R", originalType2: "I",
	},
	"q6": {
		ID:            "q6",
		Axis:          "RI",
		Affirmation1:  Affirmation{Text: "El conflicto por el poder es permanente; el progreso moral global es improbable.", Type: "R", Axis: "RI"},
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
		Affirmation1:  Affirmation{Text: "Ningún país u organismo internacional debe intervenir en los asuntos internos de otro Estado sin su consentimiento, aunque se considere que ese gobierno viola derechos humanos.", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "En ciertos casos, es legítimo que la comunidad internacional intervenga en otro país si sus acciones comprometen la estabilidad regional, internacional o violan derechos humanos.", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
	"q9": {
		ID:            "q9",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Un gobierno mundial o la cesión significativa de soberanía a entidades supranacionales pondría en riesgo la autonomía de las naciones y debería evitarse.", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "Deberíamos aspirar a instituciones globales más fuertes, incluso a alguna forma de autoridad multilateral, para enfrentar desafíos que ningún país puede resolver solo.", Type: "G", Axis: "SG"},
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
		Affirmation1:  Affirmation{Text: "Cada país debería poder proteger su economía e industria, aunque eso implique restringir acuerdos internacionales o limitar el libre comercio si es necesario.", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "A cada país le conviene promover el libre comercio y la integración económica global porque le benefician a largo plazo.", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
	"q12": {
		ID:            "q12",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Un Estado soberano debe controlar estrictamente sus fronteras y decidir quién entra, sin presiones externas.", Type: "S", Axis: "SG"},
		Affirmation2:  Affirmation{Text: "La inmigración enriquece a los países, deberíamos facilitar la libre circulación de personas.", Type: "G", Axis: "SG"},
		originalType1: "S", originalType2: "G",
	},
	"q13": {
		ID:            "q13",
		Axis:          "SG",
		Affirmation1:  Affirmation{Text: "Muchos acuerdos internacionales (clima, comercio, salud, etc) limitan la capacidad de un país para actuar en su propio beneficio.", Type: "S", Axis: "SG"},
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
		description = "Este perfil enfatiza una visión pragmática y nacional de la política. El realista-soberanista valora ante todo la soberanía del Estado y la búsqueda del poder e interés nacional dentro de un sistema internacional anárquico. Considera que el orden mundial más estable es el basado en Estados fuertes, soberanos y en equilibrio de poder (sistema Westfaliano), desconfiando de esquemas globalistas o ideales universalistas que puedan limitar la independencia nacional. En este cuadrante se prioriza la seguridad, el orden interno y la autonomía del país, asumiendo que las relaciones internacionales son un juego de suma cero donde cada nación vela por sí misma."
		thinkers = []Person{
			{Name: "Nicolás Maquiavelo",
				Short:    `(1469-1527) fue un diplomático florentino del Renacimiento y autor de El Príncipe y Los discursos sobre la primera década de Tito Livio. Con estas obras inauguró el realismo político moderno: estudia la política tal como es —lucha por el poder y la seguridad del Estado— y no como "debería" ser según ideales morales.`,
				Full:     `Nicolás Maquiavelo (1469-1527) fue un diplomático florentino del Renacimiento y autor de El Príncipe y Los discursos sobre la primera década de Tito Livio. Con estas obras inauguró el realismo político moderno: estudia la política tal como es —lucha por el poder y la seguridad del Estado— y no como "debería" ser según ideales morales. Defiende que el gobernante debe usar cualquier medio —incluso la fuerza o el engaño— si con ello protege la estabilidad y los intereses del Estado. Separa ética privada y razón de Estado, núcleo del realismo político. Escribe para príncipes que necesitan conservar la autonomía y el poder interno frente a potencias externas (el Papado, el Sacro Imperio) y rivales internos. La prioridad absoluta es la soberanía y seguridad del propio Estado; rechaza toda tutela externa sobre las decisiones del gobernante, anticipando la lógica westfaliana de Estados soberanos.`,
				ImageURL: "/static/img/maquiavelo.jpg"},
			{Name: "John Mearsheimer",
				Short:    `Es un influyente politólogo estadounidense y destacado académico de relaciones internacionales, conocido principalmente por desarrollar la teoría del 'Realismo Ofensivo'.`,
				Full:     `Es un influyente politólogo estadounidense y destacado académico de relaciones internacionales, conocido principalmente por desarrollar la teoría del 'Realismo Ofensivo'. Sostiene que las grandes potencias buscan maximizar su poder y aspiran a la hegemonía para garantizar su supervivencia en el sistema internacional anárquico. Se sitúa como Realista-Soberanista porque su teoría del 'realismo ofensivo' considera a los Estados soberanos como los actores principales, obligados a buscar poder para sobrevivir en un mundo anárquico y competitivo. Por ello, desconfía profundamente de la capacidad de las instituciones globales para alterar esta lógica de poder y prioriza la seguridad y el interés nacional por encima de consideraciones morales o de cooperación idealista, defendiendo así la primacía de la acción estatal autónoma.`,
				ImageURL: "/static/img/mearsheimer.jpg"},
		}
		politicians = []Person{
			{Name: "Xi Jinping",
				Short:    `Presidente de la República Popular China y Secretario General del Partido Comunista. Encaja como Realista porque prioriza el interés nacional, la seguridad y el poder de China por encima de consideraciones morales universales en las relaciones internacionales.`,
				Full:     `Presidente de la República Popular China y Secretario General del Partido Comunista. Encaja como Realista porque prioriza el interés nacional, la seguridad y el poder de China por encima de consideraciones morales universales en las relaciones internacionales. Es Soberanista por su férrea defensa de la autonomía nacional, el rechazo a la injerencia externa en asuntos internos y la promoción de un modelo de desarrollo propio sin imposiciones foráneas. Considera la cooperación global principalmente como una herramienta pragmática para avanzar estos objetivos nacionales.`,
				ImageURL: "/static/img/xijinping.jpg"},
			{Name: "Vladimir Putin",
				Short:    `Actual Presidente de la Federación de Rusia, figura política dominante en el país desde principios del siglo XXI. Su largo mandato ha estado marcado por la restauración del poder estatal ruso.`,
				Full:     `Actual Presidente de la Federación de Rusia, figura política dominante en el país desde principios del siglo XXI. Su largo mandato ha estado marcado por la restauración del poder estatal ruso. Su política prioriza el interés nacional, la seguridad y el poder de Rusia (Realismo), viendo las relaciones internacionales principalmente como una competición estratégica. Al mismo tiempo, defiende férreamente la autoridad del Estado ruso, promueve una identidad nacional fuerte y la soberanía nacional frente a cualquier injerencia externa o limitación por parte de instituciones globales (Soberanismo).`,
				ImageURL: "/static/img/putin.jpg"},
		}
	case "Realista-Globalista":
		description = `Quienes se ubican en el cuadrante realista-globalista comparten la visión de que la política mundial es principalmente una competencia estratégica y que los Estados actúan guiados por interés propio, pero al mismo tiempo reconocen y operan dentro de la interdependencia global. Este enfoque hace hincapié en la naturaleza competitiva del sistema internacional y asume que los estados operan en un entorno sin justicia inherente, donde las normas éticas pueden quedar supeditadas al poder. Un realista-globalista típicamente apoya la cooperación internacional solo de forma pragmática, por ejemplo, mediante alianzas o instituciones, siempre y cuando esto beneficie el equilibrio de poder o los intereses de su propio país. También suele abogar por un liderazgo fuerte de las potencias para mantener la estabilidad global, en lugar de confiar en principios idealistas.`
		thinkers = []Person{
			{Name: "Zbigniew Brzezinski",
				Short:    `Fue un influyente diplomático y politólogo polaco-estadounidense, conocido principalmente por ser Consejero de Seguridad Nacional del presidente Jimmy Carter. Su pensamiento estratégico se centró en la geopolítica y tuvo un gran impacto en la política exterior de Estados Unidos durante y después de la Guerra Fría.`,
				Full:     `Fue un influyente diplomático y politólogo polaco-estadounidense, conocido principalmente por ser Consejero de Seguridad Nacional del presidente Jimmy Carter. Su pensamiento estratégico se centró en la geopolítica y tuvo un gran impacto en la política exterior de Estados Unidos durante y después de la Guerra Fría. Se le considera Realista-Globalista porque su análisis priorizaba descarnadamente el poder, la estrategia y el interés nacional estadounidense (Realismo), tal como expuso en 'El Gran Tablero Mundial'. Sin embargo, su campo de acción y visión eran intrínsecamente globales, abogando por gestionar activamente las relaciones internacionales y usar las estructuras de poder a nivel mundial (Globalismo) para asegurar la primacía de EE. UU., en lugar de enfocarse en ideales universales o en un soberanismo defensivo.`,
				ImageURL: "/static/img/brzezinski.jpg"},
			{Name: "Nicholas Spykman",
				Short:    `Fue un influyente geoestratega y profesor neerlandés-estadounidense, considerado uno de los padres del realismo geopolítico y célebre por su teoría del 'Rimland' sobre la importancia estratégica de las zonas costeras de Eurasia.`,
				Full:     `Fue un influyente geoestratega y profesor neerlandés-estadounidense, considerado uno de los padres del realismo geopolítico y célebre por su teoría del 'Rimland' sobre la importancia estratégica de las zonas costeras de Eurasia. Se sitúa como Realista-Globalista porque, si bien su análisis se basa crudamente en el poder, la geografía y el interés nacional (Realismo), defendía que la seguridad de un Estado (como EE.UU.) requería una activa intervención y la formación de alianzas a escala mundial. Argumentaba vehementemente en contra del aislacionismo, insistiendo en que la proyección de poder y el control del equilibrio en zonas clave lejanas como el 'Rimland' eran esenciales para la supervivencia y primacía del Estado (enfoque Globalista estratégico).`,
				ImageURL: "/static/img/spykman.jpg"},
		}
		politicians = []Person{
			{Name: "Deng Xiaoping",
				Short:    `Fue el líder supremo de China tras Mao Zedong, considerado el arquitecto de la política de 'Reforma y Apertura' iniciada en 1978. Utilizó estratégicamente la integración en la economía mundial (globalismo instrumental) como la herramienta más eficaz para alcanzar la modernización y la prosperidad. `,
				Full:     `Fue el líder supremo de China tras Mao Zedong, considerado el arquitecto de la política de 'Reforma y Apertura' iniciada en 1978. Utilizó estratégicamente la integración en la economía mundial (globalismo instrumental) como la herramienta más eficaz para alcanzar la modernización y la prosperidad. Su famoso lema 'No importa si el gato es blanco o negro, mientras cace ratones' encapsula su pragmatismo. Dejó de lado la rigidez ideológica maoísta para centrarse en la política de 'Reforma y Apertura' fue una decisión estratégica calculada para aumentar el poder nacional integral de China, no basada en ideales universales de cooperación, sino en la necesidad práctica de modernizar el país. Su enfoque en 'ocultar la fuerza y esperar el momento' también es una clara estrategia realista de acumulación de poder. Tiene una forma pragmática de Globalismo instrumental, al servicio de fines Realistas y Soberanistas a largo plazo.`,
				ImageURL: "/static/img/xiaoping.jpg"},
			{Name: "Henry Kissinger",
				Short:    `Estadista y diplomático, fue un influyente diplomático y politólogo germano-estadounidense, que ejerció como Secretario de Estado y Consejero de Seguridad Nacional de EE.UU., moldeando decisivamente la política exterior estadounidense durante la Guerra Fría.`,
				Full:     `Estadista y diplomático, fue un influyente diplomático y politólogo germano-estadounidense, que ejerció como Secretario de Estado y Consejero de Seguridad Nacional de EE.UU., moldeando decisivamente la política exterior estadounidense durante la Guerra Fría. Es una figura central asociada a la Realpolitik y la diplomacia pragmática a escala global. Encaja como Globalista (en un sentido instrumental y estratégico) porque operó a escala mundial, utilizando la diplomacia compleja (como la apertura a China o la distensión con la URSS) y las negociaciones globales como herramientas clave para gestionar activamente el sistema internacional y asegurar ventajas estratégicas para Estados Unidos. Su objetivo no era un ideal globalista per se, sino gestionar un orden global desde una perspectiva de poder e interés nacional.`,
				ImageURL: "/static/img/kissinger.jpg"},
		}
	case "Idealista-Soberanista":
		description = `Este perfil combina la defensa de la soberanía nacional con la adhesión a principios e ideales claros. El idealista-soberanista cree en la autodeterminación de los pueblos y en la importancia de que cada nación sea libre para perseguir sus propios ideales. Suele oponerse a cualquier dominación imperial o extranjera sobre una nación, argumentando que la legitimidad política nace de los valores y la voluntad del propio pueblo. Por ejemplo, se enfatiza que ninguna nación merece la libertad si no es capaz de conquistarla por sí misma, reflejando un fuerte compromiso con la independencia nacional y la justicia. A diferencia del realista puro, este perfil sí otorga un peso moral a la política, pero concentrado en el ámbito nacional: la soberanía se valora como un medio para garantizar libertad, dignidad y progreso de la nación conforme a sus propios ideales (ya sean democracia, igualdad, tradición cultural, etc.).`
		thinkers = []Person{
			{Name: "Giuseppe Mazzini",
				Short:    `Pensador del siglo XIX. Político, periodista y activista italiano, una figura clave del Risorgimento (la unificación italiana). Fundador de la organización 'Joven Italia'. Mazzini era un republicano ferviente y creía profundamente en ideales como la libertad, la igualdad, el progreso y la fraternidad entre los pueblos, a menudo con un tinte casi religioso ('Dios y el Pueblo').`,
				Full:     `Pensador del siglo XIX. Político, periodista y activista italiano, una figura clave del Risorgimento (la unificación italiana). Fundador de la organización 'Joven Italia'. Mazzini era un republicano ferviente y creía profundamente en ideales como la libertad, la igualdad, el progreso y la fraternidad entre los pueblos, a menudo con un tinte casi religioso ('Dios y el Pueblo'). Creía que cada nación tenía una misión moral específica que cumplir para el progreso de la humanidad. Su visión no era meramente pragmática, sino basada en principios éticos y políticos elevados. Su objetivo principal era liberar a Italia del dominio extranjero (austríaco, papal, borbónico) y unificarla como una república independiente y soberana. Creía firmemente en el principio de autodeterminación nacional: cada pueblo tenía derecho a gobernarse a sí mismo y a formar su propio Estado-nación libre de injerencias externas. Su idealismo estaba intrínsecamente ligado a la consecución de la soberanía nacional italiana para que esta pudiera cumplir su 'misión'.`,
				ImageURL: "/static/img/mazzini.jpg"},
			{Name: "Mahatma Gandhi",
				Short:    `Fue el líder preeminente del movimiento de independencia de la India contra el dominio británico. Es célebre mundialmente por su filosofía de la desobediencia civil no violenta (Satyagraha), que inspiró movimientos por los derechos civiles en todo el mundo.`,
				Full:     `Fue el líder preeminente del movimiento de independencia de la India contra el dominio británico. Es célebre mundialmente por su filosofía de la desobediencia civil no violenta (Satyagraha), que inspiró movimientos por los derechos civiles en todo el mundo. Gandhi encaja como Idealista porque fundamentó su lucha en principios éticos y espirituales profundos, como la verdad y la no violencia (Ahimsa), creyendo firmemente en la superioridad moral de estos métodos para lograr un cambio político y social justo. A su vez, es Soberanista porque su objetivo central e irrenunciable fue alcanzar la independencia completa (Purna Swaraj) de la India del control imperial británico, defendiendo el derecho inalienable de la nación india a la autodeterminación y al autogobierno integral (Swaraj).`,
				ImageURL: "/static/img/gandhi.jpg"},
		}
		politicians = []Person{
			{Name: "Charles de Gaulle",
				Short:    `Militar y estadista francés, líder de la Francia Libre durante la Segunda Guerra Mundial y arquitecto de la Quinta República Francesa, de la que fue el primer presidente. De Gaulle tenía una idea muy particular y elevada de Francia, una visión casi mística de su grandeza, su historia y su papel en el mundo.`,
				Full:     `Militar y estadista francés, líder de la Francia Libre durante la Segunda Guerra Mundial y arquitecto de la Quinta República Francesa, de la que fue el primer presidente. De Gaulle tenía una idea muy particular y elevada de Francia, una visión casi mística de su grandeza, su historia y su papel en el mundo. Creía en la singularidad de la civilización francesa y en la necesidad de preservar sus valores republicanos y su independencia cultural. Su política estaba guiada por este ideal de la 'grandeza' y la misión histórica de Francia. Es quizás uno de los soberanistas más emblemáticos del siglo XX. Defendió a ultranza la independencia nacional de Francia frente a las superpotencias (EE.UU. y la URSS), desarrollando un programa nuclear propio, retirando a Francia del mando militar integrado de la OTAN y vetando la entrada del Reino Unido en la Comunidad Económica Europea por considerarlo un 'caballo de Troya' de los intereses estadounidenses. Su ideal de Francia solo podía realizarse a través de una soberanía política, militar y económica plena.`,
				ImageURL: "/static/img/degaulle.jpg"},
			{Name: "Simón Bolívar",
				Short:    `Fue un militar y político hispanoamericano (1783-1830), conocido como "El Libertador" por liderar la independencia de gran parte de Sudamérica del dominio español.`,
				Full:     `Bolívar encaja como Idealista-Soberanista porque su lucha emancipadora estaba profundamente inspirada en ideales ilustrados de libertad, igualdad y autogobierno republicano (Idealismo). Estaba convencido de que los pueblos debían autodeterminarse y construir su destino. Rechazó tanto el colonialismo como el dominio extranjero posterior. La independencia era el vehículo necesario para realizar su visión ideal de naciones libres.`,
				ImageURL: "/static/img/bolivar.png"},
		}
	case "Idealista-Globalista":
		description = `Este cuadrante representa una postura que apuesta por principios universales y cooperación a escala global. El idealista-globalista cree que los Estados deben trascender sus diferencias en pos de objetivos comunes de la humanidad, y tiende a apoyar la creación de instituciones internacionales fuertes e incluso estructuras de gobernanza global. Históricamente, esta visión se inspira en ideas como las de Immanuel Kant, quien proponía lograr una "paz perpetua" mediante una federación de repúblicas y la cooperación internacional basada en la razón y la moral.  En este perfil se priorizan valores como los derechos humanos, el derecho internacional, el desarrollo sostenible y la resolución pacífica de conflictos a través del diálogo multinacional. Un idealista-globalista está dispuesto a ceder parte de la soberanía nacional en favor de instituciones o acuerdos globales si con ello se alcanzan bienes mayores (por ejemplo, la paz mundial, la lucha contra el cambio climático o la erradicación de la pobreza).`
		thinkers = []Person{
			{Name: "Immanuel Kant",
				Short:    `Fue un filósofo alemán del siglo XVIII, considerado uno de los pensadores más influyentes de la Ilustración y de la filosofía occidental moderna. `,
				Full:     `Fue un filósofo alemán del siglo XVIII, considerado uno de los pensadores más influyentes de la Ilustración y de la filosofía occidental moderna. Se sitúa en la categoría Idealista-Globalista porque desarrolló el idealismo trascendental, una teoría que sostiene que el conocimiento humano se basa en estructuras mentales innatas que dan forma a nuestra experiencia del mundo . Además, en su ensayo La paz perpetua propuso la creación de una federación de estados republicanos como medio para alcanzar una paz duradera a través de la cooperación internacional y el respeto mutuo entre las naciones.`,
				ImageURL: "/static/img/kant.webp"},
			{Name: "John Rawls",
				Short:    `Fue un filósofo político estadounidense, reconocido por su influyente obra Teoría de la justicia (1971), donde propuso el concepto de 'justicia como equidad'. `,
				Full:     `Fue un filósofo político estadounidense, reconocido por su influyente obra Teoría de la justicia (1971), donde propuso el concepto de 'justicia como equidad'. Se sitúa en la categoría Idealista-Globalista porque su teoría se basa en principios éticos universales aplicables a todas las sociedades, enfatizando la equidad y los derechos fundamentales. Además, en su obra El Derecho de los Pueblos (1999), abogó por una comunidad internacional justa, donde los pueblos cooperan bajo normas comunes de justicia y respeto mutuo, promoviendo una visión global de la justicia.`,
				ImageURL: "/static/img/rawls.webp"},
		}
		politicians = []Person{
			{Name: "George Soros",
				Short:    `Es un inversor y filántropo húngaro-estadounidense conocido por financiar causas progresistas y fundar la Open Society Foundations, dedicada a la promoción de la democracia, los derechos humanos y la gobernanza global. `,
				Full:     `Es un inversor y filántropo húngaro-estadounidense conocido por financiar causas progresistas y fundar la Open Society Foundations, dedicada a la promoción de la democracia, los derechos humanos y la gobernanza global. Su inclusión en esta categoría se debe principalmente a la misión de sus Open Society Foundations, que dedican miles de millones a promover globalmente la democracia liberal, los derechos humanos y la sociedad civil, reflejando un idealismo basado en valores universales. Además, es un firme defensor de la cooperación internacional, apoya instituciones supranacionales como la Unión Europea y critica abiertamente los nacionalismos, alineándose claramente con una visión globalista frente a la soberanista.`,
				ImageURL: "/static/img/soros.jpg"},
			{Name: "Barack Obama",
				Short:    `Fue el 44.º presidente de los Estados Unidos, reconocido por su enfoque en la diplomacia multilateral y la promoción de valores democráticos a nivel global.`,
				Full:     `Fue el 44.º presidente de los Estados Unidos, reconocido por su enfoque en la diplomacia multilateral y la promoción de valores democráticos a nivel global. Se sitúa en la categoría Idealista-Globalista debido a su compromiso con el multilateralismo y la cooperación internacional. Durante su mandato, promovió acuerdos como el Acuerdo de París sobre cambio climático y defendió el fortalecimiento de instituciones globales como la ONU, reflejando su creencia en que los desafíos globales requieren soluciones colectivas basadas en principios éticos compartidos.`,
				ImageURL: "/static/img/obama.jpg"},
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

// --- NUEVA FUNCIÓN: Calcular Resultado ---
func calculateResult(userChoices []UserChoice) (float64, float64, string) {
	var scoreR, scoreI, scoreS, scoreG int
	var totalRI, totalSG int

	for _, choice := range userChoices {
		chosenType := choice.ChosenType

		// No necesitamos la pregunta original aquí, solo el tipo elegido
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
		default:
			// Podríamos loguear aquí si quisiéramos, pero no afecta el cálculo
			// fmt.Printf("Advertencia: Se recibió tipo inválido ('%s') para la pregunta %s\n", chosenType, choice.QuestionID)
		}
	}

	// Normalizar puntuaciones
	var finalScoreRI float64
	if totalRI > 0 {
		// Realismo (0) <-> Idealismo (1) -> Más I, más cerca de 1
		finalScoreRI = float64(scoreI) / float64(totalRI)
	} // Si totalRI es 0, finalScoreRI permanece 0.0

	var finalScoreSG float64
	if totalSG > 0 {
		// Soberanismo (0) <-> Globalismo (1) -> Más G, más cerca de 1
		finalScoreSG = float64(scoreG) / float64(totalSG)
	} // Si totalSG es 0, finalScoreSG permanece 0.0

	// Determinar perfil
	profile := ""
	threshold := 0.5
	// Usamos una pequeña épsilon para manejar imprecisiones de punto flotante cerca del umbral
	epsilon := 1e-9
	if finalScoreRI < threshold-epsilon && finalScoreSG < threshold-epsilon {
		profile = "Realista-Soberanista"
	} else if finalScoreRI < threshold-epsilon && finalScoreSG >= threshold-epsilon {
		profile = "Realista-Globalista"
	} else if finalScoreRI >= threshold-epsilon && finalScoreSG < threshold-epsilon {
		profile = "Idealista-Soberanista"
	} else { // Cubre finalScoreRI >= threshold && finalScoreSG >= threshold
		profile = "Idealista-Globalista"
	}

	return finalScoreRI, finalScoreSG, profile
}

// --- FIN NUEVA FUNCIÓN ---

// Handler para recibir las respuestas y calcular el resultado (MODIFICADO)
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

	// --- Llamar a la nueva función para calcular ---
	finalScoreRI, finalScoreSG, profile := calculateResult(userChoices)
	// --- FIN Llamada ---

	// Obtener los datos descriptivos para el perfil calculado
	categoryResultData := getCategoryData(profile)

	// Crear el objeto resultado final
	result := Result{
		ScoreRI:     finalScoreRI,
		ScoreSG:     finalScoreSG,
		Profile:     profile,                        // Usamos el nombre calculado
		Description: categoryResultData.Description, // Obtenido de getCategoryData
		Thinkers:    categoryResultData.Thinkers,    // Obtenido de getCategoryData
		Politicians: categoryResultData.Politicians, // Obtenido de getCategoryData
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
	http.HandleFunc("/api/questions", questionsHandler)
	http.HandleFunc("/api/submit", submitHandler)
	http.HandleFunc("/api/categories", categoriesHandler)

	// Servir archivos estáticos (CSS, JS, imágenes)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Servir el HTML principal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Asegurarse de que solo se sirva para la ruta exacta "/"
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// Medir tiempo de ServeFile
		start := time.Now()
		http.ServeFile(w, r, "templates/index.html")
		duration := time.Since(start)
		fmt.Printf("DEBUG: ServeFile('/') tomó %v\n", duration)
	})

	port := "8080"
	fmt.Printf("Servidor iniciado en http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
