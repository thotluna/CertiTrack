package certifications

import (
	"certitrack/backend/shared/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cert)
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	certGroup := router.Group("/certifications")
	{
		certGroup.GET("/:id", h.GetCertification)
	}
}
