package main

import (
	"math"
	"testing"
)

// Función auxiliar para comparar floats con tolerancia
func floatsAlmostEqual(a, b float64, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestCalculateResult(t *testing.T) {
	epsilon := 1e-9 // Tolerancia para comparación de floats

	testCases := []struct {
		name            string
		inputChoices    []UserChoice
		expectedScoreRI float64
		expectedScoreSG float64
		expectedProfile string
	}{
		{
			name: "Todo Realista-Soberanista",
			inputChoices: []UserChoice{
				{QuestionID: "q1", ChosenType: "R"}, {QuestionID: "q2", ChosenType: "R"}, {QuestionID: "q3", ChosenType: "R"}, {QuestionID: "q4", ChosenType: "R"}, {QuestionID: "q5", ChosenType: "R"}, {QuestionID: "q6", ChosenType: "R"}, {QuestionID: "q7", ChosenType: "R"}, // 7 R
				{QuestionID: "q8", ChosenType: "S"}, {QuestionID: "q9", ChosenType: "S"}, {QuestionID: "q10", ChosenType: "S"}, {QuestionID: "q11", ChosenType: "S"}, {QuestionID: "q12", ChosenType: "S"}, {QuestionID: "q13", ChosenType: "S"}, {QuestionID: "q14", ChosenType: "S"}, // 7 S
			},
			expectedScoreRI: 0.0,
			expectedScoreSG: 0.0,
			expectedProfile: "Realista-Soberanista",
		},
		{
			name: "Todo Idealista-Globalista",
			inputChoices: []UserChoice{
				{QuestionID: "q1", ChosenType: "I"}, {QuestionID: "q2", ChosenType: "I"}, {QuestionID: "q3", ChosenType: "I"}, {QuestionID: "q4", ChosenType: "I"}, {QuestionID: "q5", ChosenType: "I"}, {QuestionID: "q6", ChosenType: "I"}, {QuestionID: "q7", ChosenType: "I"}, // 7 I
				{QuestionID: "q8", ChosenType: "G"}, {QuestionID: "q9", ChosenType: "G"}, {QuestionID: "q10", ChosenType: "G"}, {QuestionID: "q11", ChosenType: "G"}, {QuestionID: "q12", ChosenType: "G"}, {QuestionID: "q13", ChosenType: "G"}, {QuestionID: "q14", ChosenType: "G"}, // 7 G
			},
			expectedScoreRI: 1.0,
			expectedScoreSG: 1.0,
			expectedProfile: "Idealista-Globalista",
		},
		{
			name: "Todo Realista-Globalista",
			inputChoices: []UserChoice{
				{QuestionID: "q1", ChosenType: "R"}, {QuestionID: "q2", ChosenType: "R"}, {QuestionID: "q3", ChosenType: "R"}, {QuestionID: "q4", ChosenType: "R"}, {QuestionID: "q5", ChosenType: "R"}, {QuestionID: "q6", ChosenType: "R"}, {QuestionID: "q7", ChosenType: "R"}, // 7 R
				{QuestionID: "q8", ChosenType: "G"}, {QuestionID: "q9", ChosenType: "G"}, {QuestionID: "q10", ChosenType: "G"}, {QuestionID: "q11", ChosenType: "G"}, {QuestionID: "q12", ChosenType: "G"}, {QuestionID: "q13", ChosenType: "G"}, {QuestionID: "q14", ChosenType: "G"}, // 7 G
			},
			expectedScoreRI: 0.0,
			expectedScoreSG: 1.0,
			expectedProfile: "Realista-Globalista",
		},
		{
			name: "Todo Idealista-Soberanista",
			inputChoices: []UserChoice{
				{QuestionID: "q1", ChosenType: "I"}, {QuestionID: "q2", ChosenType: "I"}, {QuestionID: "q3", ChosenType: "I"}, {QuestionID: "q4", ChosenType: "I"}, {QuestionID: "q5", ChosenType: "I"}, {QuestionID: "q6", ChosenType: "I"}, {QuestionID: "q7", ChosenType: "I"}, // 7 I
				{QuestionID: "q8", ChosenType: "S"}, {QuestionID: "q9", ChosenType: "S"}, {QuestionID: "q10", ChosenType: "S"}, {QuestionID: "q11", ChosenType: "S"}, {QuestionID: "q12", ChosenType: "S"}, {QuestionID: "q13", ChosenType: "S"}, {QuestionID: "q14", ChosenType: "S"}, // 7 S
			},
			expectedScoreRI: 1.0,
			expectedScoreSG: 0.0,
			expectedProfile: "Idealista-Soberanista",
		},
		{
			name: "Mitad R/I, Mitad S/G (Ligeramente RS)",
			inputChoices: []UserChoice{
				{QuestionID: "q1", ChosenType: "R"}, {QuestionID: "q2", ChosenType: "R"}, {QuestionID: "q3", ChosenType: "R"}, {QuestionID: "q4", ChosenType: "R"}, // 4 R
				{QuestionID: "q5", ChosenType: "I"}, {QuestionID: "q6", ChosenType: "I"}, {QuestionID: "q7", ChosenType: "I"}, // 3 I
				{QuestionID: "q8", ChosenType: "S"}, {QuestionID: "q9", ChosenType: "S"}, {QuestionID: "q10", ChosenType: "S"}, {QuestionID: "q11", ChosenType: "S"}, // 4 S
				{QuestionID: "q12", ChosenType: "G"}, {QuestionID: "q13", ChosenType: "G"}, {QuestionID: "q14", ChosenType: "G"}, // 3 G
			},
			expectedScoreRI: 3.0 / 7.0,
			expectedScoreSG: 3.0 / 7.0,
			expectedProfile: "Realista-Soberanista", // Ambos scores < 0.5
		},
		{
			name: "Mitad R/I, Mitad S/G (Ligeramente IG)",
			inputChoices: []UserChoice{
				{QuestionID: "q1", ChosenType: "R"}, {QuestionID: "q2", ChosenType: "R"}, {QuestionID: "q3", ChosenType: "R"}, // 3 R
				{QuestionID: "q4", ChosenType: "I"}, {QuestionID: "q5", ChosenType: "I"}, {QuestionID: "q6", ChosenType: "I"}, {QuestionID: "q7", ChosenType: "I"}, // 4 I
				{QuestionID: "q8", ChosenType: "S"}, {QuestionID: "q9", ChosenType: "S"}, {QuestionID: "q10", ChosenType: "S"}, // 3 S
				{QuestionID: "q11", ChosenType: "G"}, {QuestionID: "q12", ChosenType: "G"}, {QuestionID: "q13", ChosenType: "G"}, {QuestionID: "q14", ChosenType: "G"}, // 4 G
			},
			expectedScoreRI: 4.0 / 7.0,
			expectedScoreSG: 4.0 / 7.0,
			expectedProfile: "Idealista-Globalista", // Ambos scores >= 0.5
		},
		{
			name:            "Caso Vacío",
			inputChoices:    []UserChoice{},
			expectedScoreRI: 0.0,
			expectedScoreSG: 0.0,
			expectedProfile: "Realista-Soberanista", // Lógica actual lleva a esto
		},
		{
			name: "Solo respuestas RI (Todas R)",
			inputChoices: []UserChoice{
				{QuestionID: "q1", ChosenType: "R"}, {QuestionID: "q2", ChosenType: "R"}, {QuestionID: "q3", ChosenType: "R"}, {QuestionID: "q4", ChosenType: "R"}, {QuestionID: "q5", ChosenType: "R"}, {QuestionID: "q6", ChosenType: "R"}, {QuestionID: "q7", ChosenType: "R"},
			},
			expectedScoreRI: 0.0,
			expectedScoreSG: 0.0, // totalSG es 0
			expectedProfile: "Realista-Soberanista",
		},
		{
			name: "Solo respuestas SG (Todas G)",
			inputChoices: []UserChoice{
				{QuestionID: "q8", ChosenType: "G"}, {QuestionID: "q9", ChosenType: "G"}, {QuestionID: "q10", ChosenType: "G"}, {QuestionID: "q11", ChosenType: "G"}, {QuestionID: "q12", ChosenType: "G"}, {QuestionID: "q13", ChosenType: "G"}, {QuestionID: "q14", ChosenType: "G"},
			},
			expectedScoreRI: 0.0, // totalRI es 0
			expectedScoreSG: 1.0,
			expectedProfile: "Realista-Globalista", // RI=0 < 0.5, SG=1.0 >= 0.5
		},
		{
			name: "Con respuesta inválida (debe ignorarse)",
			inputChoices: []UserChoice{
				{QuestionID: "q1", ChosenType: "R"}, {QuestionID: "q2", ChosenType: "R"}, {QuestionID: "q3", ChosenType: "R"}, {QuestionID: "q4", ChosenType: "R"}, {QuestionID: "q5", ChosenType: "R"}, {QuestionID: "q6", ChosenType: "R"}, {QuestionID: "q7", ChosenType: "R"}, // 7 R
				{QuestionID: "q8", ChosenType: "S"}, {QuestionID: "q9", ChosenType: "S"}, {QuestionID: "q10", ChosenType: "S"}, {QuestionID: "q11", ChosenType: "S"}, {QuestionID: "q12", ChosenType: "S"}, {QuestionID: "q13", ChosenType: "S"}, {QuestionID: "q14", ChosenType: "S"}, // 7 S
				{QuestionID: "q_invalid", ChosenType: "X"}, // Tipo inválido
			},
			expectedScoreRI: 0.0,
			expectedScoreSG: 0.0,
			expectedProfile: "Realista-Soberanista", // La X se ignora
		},
		{
			name: "Puntuación exacta en el límite 0.5 (Caso IG)",
			inputChoices: []UserChoice{
				// 2 I / 4 RI total = 0.5
				{QuestionID: "q1", ChosenType: "R"}, {QuestionID: "q2", ChosenType: "R"},
				{QuestionID: "q3", ChosenType: "I"}, {QuestionID: "q4", ChosenType: "I"},
				// 2 G / 4 SG total = 0.5
				{QuestionID: "q8", ChosenType: "S"}, {QuestionID: "q9", ChosenType: "S"},
				{QuestionID: "q10", ChosenType: "G"}, {QuestionID: "q11", ChosenType: "G"},
			},
			expectedScoreRI: 0.5,
			expectedScoreSG: 0.5,
			expectedProfile: "Idealista-Globalista", // >= 0.5 para ambos
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualScoreRI, actualScoreSG, actualProfile := calculateResult(tc.inputChoices)

			if !floatsAlmostEqual(actualScoreRI, tc.expectedScoreRI, epsilon) {
				t.Errorf("ScoreRI: esperado %f, obtenido %f", tc.expectedScoreRI, actualScoreRI)
			}

			if !floatsAlmostEqual(actualScoreSG, tc.expectedScoreSG, epsilon) {
				t.Errorf("ScoreSG: esperado %f, obtenido %f", tc.expectedScoreSG, actualScoreSG)
			}

			if actualProfile != tc.expectedProfile {
				t.Errorf("Profile: esperado %s, obtenido %s", tc.expectedProfile, actualProfile)
			}
		})
	}
}
