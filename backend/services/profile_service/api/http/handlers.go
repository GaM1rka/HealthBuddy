package http

import (
	"net/http"

	"profile_service/domain"
	"profile_service/usecases"

	"github.com/gin-gonic/gin"
)

// ProfileHandler handles HTTP requests for profiles
// and delegates to the ProfileService business logic.
type ProfileHandler struct {
	svc usecases.ProfileService
}

// NewProfileHandler constructs a new ProfileHandler
func NewProfileHandler(svc usecases.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// RegisterRoutes registers profile routes on the Gin engine
func (h *ProfileHandler) RegisterRoutes(r *gin.Engine) {
	grp := r.Group("/profile")
	{
		grp.GET("/health", h.Health)
		grp.POST("", h.Create)
		/* grp.GET("", h.List) */
		grp.GET("", h.GetByID)
		grp.PUT("", h.Update)
		grp.DELETE("", h.Delete)
	}
}

// Health endpoint
func (h *ProfileHandler) Health(c *gin.Context) {
	if err := h.svc.Health(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db down"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// Create handles POST /profile
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

// GetByID handles GET /profile
func (h *ProfileHandler) GetByID(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	out, err := h.svc.GetProfile(c.Request.Context(), userID)
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

// List handles GET /profile
/* func (h *ProfileHandler) List(c *gin.Context) {
	list, err := h.svc.ListProfiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
} */

// Update handles PUT /profile
// Only the owner of the profile can update it.
func (h *ProfileHandler) Update(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}
	// Fetch existing profile to verify ownership
	existing, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if err == usecases.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if existing.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to update this profile"})
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

	out, err := h.svc.UpdateProfile(c.Request.Context(), userID, req)
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

// Delete handles DELETE /profile/:id
// Only the owner of the profile can delete it.
func (h *ProfileHandler) Delete(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}
	// Fetch existing profile to verify ownership
	existing, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if err == usecases.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if existing.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to delete this profile"})
		return
	}

	if err := h.svc.DeleteProfile(c.Request.Context(), userID); err != nil {
		if err == usecases.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
