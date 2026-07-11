package handlers

import (
	"net/http"
	"strconv"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type TestDriveHandler struct {
	testDriveService *services.TestDriveService
}

func NewTestDriveHandler(testDriveService *services.TestDriveService) *TestDriveHandler {
	return &TestDriveHandler{testDriveService: testDriveService}
}

func (h *TestDriveHandler) Create(c *gin.Context) {
	var req models.CreateTestDriveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID := c.GetUint("user_id")
	td, err := h.testDriveService.Create(clientID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, td)
}

func (h *TestDriveHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	td, err := h.testDriveService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "turno no encontrado"})
		return
	}

	c.JSON(http.StatusOK, td)
}

func (h *TestDriveHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var req struct {
		Status models.TestDriveStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.testDriveService.UpdateStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "estado actualizado"})
}

func (h *TestDriveHandler) ListMy(c *gin.Context) {
	clientID := c.GetUint("user_id")
	list, err := h.testDriveService.ListByClient(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *TestDriveHandler) ListAll(c *gin.Context) {
	list, err := h.testDriveService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}
