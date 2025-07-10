# Sociopolitical Compass 🧭 (English)
https://radar.ignaciopadron.es/

## 🌐 Select your language / Selecciona tu idioma

**🇺🇸 [English](README_EN.md)** | **🇪🇸 [Español](README.md)**

---

> An interactive web application to discover your sociopolitical orientation on the Realism/Idealism and Sovereignty/Globalism axes.

This application presents a 14-question survey designed to evaluate your political and philosophical inclinations. Upon completion, it will show your position on a compass chart and associate you with one of the four defined profiles, along with related thinkers and political figures.

## 📜 Table of Contents

- [✨ Features](#-features)
- [🛠️ Architecture](#️-architecture)
- [📦 Requirements](#-requirements)
- [🔧 Installation](#-installation)
- [🚀 Usage](#-usage)
- [🔌 API Endpoints](#-api-endpoints)
- [🤝 Contributions](#-contributions)

## ✨ Features

-   **Interactive Survey:** 14 pairs of statements to evaluate your position.
-   **Randomization:** The order of questions and the position of statements are shuffled in each session.
-   **Result Calculation:** Numerical scores on the Realism-Idealism (RI) and Sovereignty-Globalism (SG) axes.
-   **Detailed Profiles:** Assignment to one of the four profiles:
    -   Realist-Sovereignist
    -   Realist-Globalist
    -   Idealist-Sovereignist
    -   Idealist-Globalist
-   **Graphical Visualization:** A compass-type chart shows your exact position.
-   **Related Figures:** Thinkers and politicians associated with the resulting profile are displayed.
-   **Category Exploration:** A dedicated section to learn about the four profiles, their descriptions and associated figures.
-   **Responsive Interface:** Adaptive design thanks to Bootstrap.

## 🛠️ Architecture

The application uses a simple architecture:

-   **Backend:** Go (using the standard library `net/http` for the web server and `encoding/json` for the API).
-   **Frontend:** HTML5, CSS3 (with Bootstrap 5.3) and JavaScript (vanilla JS) for the user interface and interaction.
-   **Web Server:** The Go backend serves both static files (HTML, CSS, JS, images) and the REST API.

## 📦 Requirements

-   Go 1.16 or higher.
-   A modern web browser (Chrome, Firefox, Edge, Safari).

## 🔧 Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/ignaciopadron/socio-political-survey
    cd socio-political-survey
    ```
2.  Run the application:
    ```bash
    go run cmd/main.go
    ```
    This will compile and run the backend server.

## 🚀 Usage

Once the server is running (you'll see the message "Server started at http://localhost:8080"), open your web browser and visit:

<http://localhost:8080>

You can navigate between the survey and the categories section using the links in the top navigation bar.

## 🔌 API Endpoints

The backend exposes the following REST endpoints:

-   `GET /api/questions`
    -   **Description:** Returns the list of 14 survey questions in random order and with statements also randomly ordered within each pair.
    -   **Response:** `[]QuestionPair` (see `cmd/main.go` for the structure).
-   `POST /api/submit`
    -   **Description:** Receives user responses and calculates the final result.
    -   **Request Body:** `[]UserChoice` (E.g.: `[{"questionId":"q5","chosenType":"R"}, {"questionId":"q12","chosenType":"G"}, ...]`)
    -   **Response:** `Result` (includes scores, profile, description, associated thinkers and politicians).
-   `GET /api/categories`
    -   **Description:** Returns detailed information (description, thinkers, politicians) of the four possible categories/profiles.
    -   **Response:** `[]CategoryData`.

## 🤝 Contributions

Contributions are welcome. If you find a bug or have a suggestion, please open an *issue* in the repository.

---

We hope you enjoy discovering your place on the Sociopolitical Compass! 🧭

---

**📖 También disponible en español / Also available in Spanish:** [README.md](README.md) 