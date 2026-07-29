package handlers

import (
	"net/http"
	"strconv"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ConsultationHandler struct {
	consultationService *services.ConsultationService
}

func NewConsultationHandler(consultationService *services.ConsultationService) *ConsultationHandler {
	return &ConsultationHandler{consultationService: consultationService}
}

func (h *ConsultationHandler) Create(c *gin.Context) {
	var req models.CreateConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID := c.GetUint("user_id")
	consultation, err := h.consultationService.Create(clientID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, consultation)
}

func (h *ConsultationHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	consultation, err := h.consultationService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "consulta no encontrada"})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if role == string(models.RoleSeller) || role == string(models.RoleAdmin) {
		h.consultationService.MarkAsRead(uint(id))
	} else if consultation.ClientID == userID {
		h.consultationService.MarkAsReadForClient(uint(id))
	}

	c.JSON(http.StatusOK, consultation)
}

func (h *ConsultationHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req struct {
		Status models.ConsultationStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sellerID := c.GetUint("user_id")
	if err := h.consultationService.UpdateStatus(uint(id), req.Status, sellerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "estado actualizado"})
}

func (h *ConsultationHandler) ListMy(c *gin.Context) {
	clientID := c.GetUint("user_id")
	consultations, err := h.consultationService.ListByClient(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": consultations})
}

func (h *ConsultationHandler) ListAll(c *gin.Context) {
	consultations, err := h.consultationService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": consultations})
}

func (h *ConsultationHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	if err := h.consultationService.Delete(uint(id), userID, role); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "consulta eliminada"})
}

func (h *ConsultationHandler) GetPendingCount(c *gin.Context) {
	count, err := h.consultationService.GetPendingCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *ConsultationHandler) GetNotificationCounts(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetUint("user_id")
	counts, err := h.consultationService.GetNotificationCounts(role, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, counts)
}

func (h *ConsultationHandler) AddResponse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	consultation, err := h.consultationService.AddResponse(uint(id), userID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, consultation)
}
