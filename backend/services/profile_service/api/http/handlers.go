// api/http/handlers.go
package http

import (
	"net/http"

	"profile_service/domain"
	"profile_service/usecases"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	svc usecases.ProfileService
}

func NewProfileHandler(svc usecases.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

func (h *ProfileHandler) RegisterRoutes(r *gin.Engine) {
	grp := r.Group("/profile")
	{
		grp.GET("/health", h.Health)
		grp.POST("", h.Create)
		grp.GET("", h.List)
		grp.GET("/:id", h.GetByID)
		grp.PUT("/:id", h.Update)
		grp.DELETE("/:id", h.Delete)
	}
}

func (h *ProfileHandler) Health(c *gin.Context) {
	if err := h.svc.Health(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db down"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func (h *ProfileHandler) Create(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}

	var req domain.ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := domain.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.CreateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *ProfileHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	out, err := h.svc.GetProfile(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *ProfileHandler) List(c *gin.Context) {
	list, err := h.svc.ListProfiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *ProfileHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req domain.ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := domain.Validate(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.UpdateProfile(c.Request.Context(), id, req)
	if err != nil {
		if err == usecases.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *ProfileHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteProfile(c.Request.Context(), id); err != nil {
		if err == usecases.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
