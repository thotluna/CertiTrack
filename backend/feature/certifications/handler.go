// feature/certifications/handler.go
package certifications

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"certitrack/backend/shared/errors"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// GetCertification maneja la obtención de una certificación por ID
func (h *Handler) GetCertification(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		c.Error(errors.NewErrorResponse(
			http.StatusBadRequest, 
			"Invalid certification ID format",
		))
		return
	}

	var cert Certification
	if err := h.db.First(&cert, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(errors.ErrNotFound)
			return
		}
		// Error de base de datos
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cert)
}

// RegisterRoutes registra las rutas para el manejo de certificaciones
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	certGroup := router.Group("/certifications")
	{
		certGroup.GET("/:id", h.GetCertification)
		// Agregar más rutas aquí según sea necesario
	}
}
