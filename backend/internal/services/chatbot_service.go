package services

import (
	"context"
	"fmt"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/config"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/repositories"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type ChatbotService struct {
	vehicleRepo *repositories.VehicleRepository
	llm         llms.LLM
	enabled     bool
}

func NewChatbotService(vehicleRepo *repositories.VehicleRepository, cfg *config.Config) *ChatbotService {
	enabled := cfg.OpenAIAPIKey != ""
	var llmInstance llms.LLM
	if enabled {
		var err error
		llmInstance, err = openai.New(openai.WithToken(cfg.OpenAIAPIKey))
		if err != nil {
			enabled = false
		}
	}

	return &ChatbotService{
		vehicleRepo: vehicleRepo,
		llm:         llmInstance,
		enabled:     enabled,
	}
}

func (s *ChatbotService) IsEnabled() bool {
	return s.enabled
}

func (s *ChatbotService) Ask(question string) (string, error) {
	if !s.enabled {
		return "El asistente no esta disponible en este momento.", nil
	}

	vehicles, err := s.vehicleRepo.ListAll()
	if err != nil {
		return "", err
	}

	var inventoryContext string
	for _, v := range vehicles {
		inventoryContext += fmt.Sprintf("- %s %s %d | %s | %s | %.0f km | $%.2f | %s\n",
			v.Brand, v.Model, v.Year, v.Fuel, v.Transmission, v.Mileage, v.Price, v.Status)
	}

	prompt := fmt.Sprintf(`Eres un asistente de una concesionaria de autos. Tu tarea es responder preguntas sobre el inventario disponible.

Inventario actual:
%s

Pregunta del usuario: %s

Responde de manera amable y util. Si la pregunta no esta relacionada con autos o el inventario, redirige amablemente al tema de la concesionaria.`, inventoryContext, question)

	ctx := context.Background()
	response, err := llms.GenerateFromSinglePrompt(ctx, s.llm, prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}
