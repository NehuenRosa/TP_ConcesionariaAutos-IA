package handlers

import (
	"net/http"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ChatbotHandler struct {
	chatbotService *services.ChatbotService
}

func NewChatbotHandler(chatbotService *services.ChatbotService) *ChatbotHandler {
	return &ChatbotHandler{chatbotService: chatbotService}
}

func (h *ChatbotHandler) Ask(c *gin.Context) {
	var req struct {
		Question string `json:"question" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer, err := h.chatbotService.Ask(req.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": answer})
}

func (h *ChatbotHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"enabled": h.chatbotService.IsEnabled()})
}
