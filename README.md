# Brújula Sociopolítica 🧭
https://radar.ignaciopadron.es/

> Una aplicación web interactiva para descubrir tu orientación sociopolítica en los ejes Realismo/Idealismo y Soberanismo/Globalismo.

Esta aplicación presenta una encuesta de 14 preguntas diseñadas para evaluar tus inclinaciones políticas y filosóficas. Al finalizar, te mostrará tu posición en un gráfico de brújula y te asociará con uno de los cuatro perfiles definidos, junto con pensadores y figuras políticas relacionadas.

## 📜 Tabla de Contenidos

- [✨ Características](#-características)
- [🛠️ Arquitectura](#️-arquitectura)
- [📦 Requisitos](#-requisitos)
- [🔧 Instalación](#-instalación)
- [🚀 Uso](#-uso)
- [🔌 API Endpoints](#-api-endpoints)
- [🤝 Contribuciones](#-contribuciones)

## ✨ Características

-   **Encuesta Interactiva:** 14 pares de afirmaciones para evaluar tu posición.
-   **Aleatorización:** El orden de las preguntas y la posición de las afirmaciones (izquierda/derecha) se mezclan en cada sesión.
-   **Cálculo de Resultados:** Puntuaciones numéricas en los ejes Realismo-Idealismo (RI) y Soberanismo-Globalismo (SG).
-   **Perfiles Detallados:** Asignación a uno de los cuatro perfiles:
    -   Realista-Soberanista
    -   Realista-Globalista
    -   Idealista-Soberanista
    -   Idealista-Globalista
-   **Visualización Gráfica:** Un gráfico tipo brújula muestra tu posición exacta.
-   **Figuras Relacionadas:** Se muestran pensadores y políticos asociados al perfil resultante.
-   **Exploración de Categorías:** Una sección dedicada para aprender sobre los cuatro perfiles, sus descripciones y figuras asociadas.
-   **Interfaz Responsiva:** Diseño adaptable gracias a Bootstrap.

## 🛠️ Arquitectura

La aplicación utiliza una arquitectura simple:

-   **Backend:** Go (usando la librería estándar `net/http` para el servidor web y `encoding/json` para la API).
-   **Frontend:** HTML5, CSS3 (con Bootstrap 5.3) y JavaScript (vanilla JS) para la interfaz de usuario y la interacción.
-   **Servidor Web:** El backend de Go sirve tanto los archivos estáticos (HTML, CSS, JS, imágenes) como la API REST.

## 📦 Requisitos

-   Go 1.16 o superior.
-   Un navegador web moderno (Chrome, Firefox, Edge, Safari).

## 🔧 Instalación

1.  Clona el repositorio:
    ```bash
    git clone <URL_DEL_REPOSITORIO>
    cd socio-political-survey
    ```
2.  Ejecuta la aplicación:
    ```bash
    go run cmd/main.go
    ```
    Esto compilará y ejecutará el servidor backend.

## 🚀 Uso

Una vez que el servidor esté en ejecución (verás el mensaje "Servidor iniciado en http://localhost:8080"), abre tu navegador web y visita:

<http://localhost:8080>

Puedes navegar entre la encuesta y la sección de categorías usando los enlaces de la barra de navegación superior.

## 🔌 API Endpoints

El backend expone los siguientes endpoints REST:

-   `GET /api/questions`
    -   **Descripción:** Devuelve la lista de 14 preguntas de la encuesta en un orden aleatorio y con las afirmaciones también ordenadas aleatoriamente dentro de cada par.
    -   **Respuesta:** `[]QuestionPair` (ver `cmd/main.go` para la estructura).
-   `POST /api/submit`
    -   **Descripción:** Recibe las respuestas del usuario y calcula el resultado final.
    -   **Cuerpo (Body) de la Petición:** `[]UserChoice` (Ej: `[{"questionId":"q5","chosenType":"R"}, {"questionId":"q12","chosenType":"G"}, ...]`)
    -   **Respuesta:** `Result` (incluye puntuaciones, perfil, descripción, pensadores y políticos asociados).
-   `GET /api/categories`
    -   **Descripción:** Devuelve la información detallada (descripción, pensadores, políticos) de las cuatro categorías/perfiles posibles.
    -   **Respuesta:** `[]CategoryData`.

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Si encuentras un error o tienes una sugerencia, por favor abre un *issue* en el repositorio.



---

¡Esperamos que disfrutes descubriendo tu lugar en la Brújula Sociopolítica! 🧭
