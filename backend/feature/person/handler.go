package person

import (
	"certitrack/backend/feature/person/dto"
	"certitrack/backend/shared/errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	persons := router.Group("/persons")
	{
		persons.GET("", h.ListPersons)
		persons.GET("/:id", h.GetPerson)
		persons.POST("", h.CreatePerson)
		persons.PUT("/:id", h.UpdatePerson)
		persons.DELETE("/:id", h.DeletePerson)
	}
}

// CreatePerson godoc
// @Summary Create a new person
// @Description Creates a new person with the provided data
// @Tags persons
// @Accept  json
// @Produce  json
// @Param   person  body      dto.CreatePersonRequest  true  "Person data"
// @Success 201 {object} dto.PersonResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/personas [post]
func (h *Handler) CreatePerson(c *gin.Context) {
	var req dto.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	person, err := h.service.CreatePerson(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toPersonResponse(person))
}

// GetPerson godoc
// @Summary Get a person by ID
// @Description Gets the details of a specific person by their ID
// @Tags persons
// @Produce  json
// @Param   id   path      string  true  "Person ID"
// @Success 200 {object} dto.PersonResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/personas/{id} [get]
func (h *Handler) GetPerson(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleError(c, errors.NewErrorResponse(http.StatusBadRequest, "invalid ID format"))
		return
	}

	person, err := h.service.GetPerson(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPersonResponse(person))
}

// ListPersons godoc
// @Summary List persons
// @Description Get a paginated list of persons with optional filters
// @Tags persons
// @Produce  json
// @Param   name    query     string  false  "Filter by name"
// @Param   status  query     string  false  "Filter by status (active, inactive, pending)"
// @Param   page    query     int     false  "Page number" default(1)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/personas [get]
func (h *Handler) ListPersons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize := 10

	filters := make(map[string]interface{})

	if name := strings.TrimSpace(c.Query("name")); name != "" {
		filters["name"] = "%" + name + "%"
	}

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		filters["status"] = status
	}

	total, err := h.service.CountPersons(c.Request.Context(), filters)
	if err != nil {
		h.handleError(c, err)
		return
	}

	offset := (page - 1) * pageSize

	filters["limit"] = pageSize
	filters["offset"] = offset

	persons, err := h.service.ListPersons(c.Request.Context(), filters)
	if err != nil {
		h.handleError(c, err)
		return
	}

	var responseData []*dto.PersonResponse
	for _, p := range persons {
		responseData = append(responseData, toPersonResponse(p))
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	c.JSON(http.StatusOK, PaginationResponse{
		Data:       responseData,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

// UpdatePerson godoc
// @Summary Update a person
// @Description Updates an existing person's data
// @Tags persons
// @Accept  json
// @Produce  json
// @Param   id      path      string  true  "Person ID"
// @Param   person  body      dto.UpdatePersonRequest  true  "Updated person data"
// @Success 200 {object} dto.PersonResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 409 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/personas/{id} [put]
func (h *Handler) UpdatePerson(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleError(c, errors.NewErrorResponse(http.StatusBadRequest, "invalid ID format"))
		return
	}

	var req dto.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	person, err := h.service.UpdatePerson(c.Request.Context(), id, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPersonResponse(person))
}

// DeletePerson godoc
// @Summary Delete a person
// @Description Deletes a person by ID (soft delete)
// @Tags persons
// @Param   id   path      string  true  "Person ID"
// @Success 204
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/personas/{id} [delete]
func (h *Handler) DeletePerson(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleError(c, errors.NewErrorResponse(http.StatusBadRequest, "invalid ID format"))
		return
	}

	if err := h.service.DeletePerson(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ChangePersonStatus godoc
// @Summary Change person status
// @Description Updates a person's status (active, inactive, pending)
// @Tags persons
// @Accept  json
// @Produce  json
// @Param   id      path      string  true  "Person ID"
// @Param   status  body      string  true  "New status"
// @Success 200 {object} dto.PersonResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 404 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router /api/v1/personas/{id}/status [patch]
func (h *Handler) ChangePersonStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.handleError(c, errors.NewErrorResponse(http.StatusBadRequest, "invalid ID format"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active inactive pending"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.handleError(c, errors.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	person, err := h.service.ChangePersonStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPersonResponse(person))
}

func toPersonResponse(p *Person) *dto.PersonResponse {
	return &dto.PersonResponse{
		ID:        p.ID.String(),
		FullName:  p.FirstName + " " + p.LastName,
		Email:     p.Email,
		Phone:     p.Phone,
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errors.ErrorResponse:
		c.JSON(e.Status, e)
	default:
		// Convert Gin validation errors to 400 error
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Validation error",
				"errors":  validationErrs,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer)
	}
}
