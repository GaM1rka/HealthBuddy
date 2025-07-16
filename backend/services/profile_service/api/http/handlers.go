package http

import (
	"net/http"
	"profile_service/api/http/apierrors"
	"profile_service/api/http/middleware"
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
		grp.GET("/health", middleware.ErrorHandlerMiddleware(h.Health))
		grp.POST("", middleware.ErrorHandlerMiddleware(h.Create))
		grp.GET("", middleware.ErrorHandlerMiddleware(h.GetByID))
		grp.PUT("", middleware.ErrorHandlerMiddleware(h.Update))
		grp.DELETE("", middleware.ErrorHandlerMiddleware(h.Delete))
	}
}

// Health endpoint
func (h *ProfileHandler) Health(c *gin.Context) error {
	if err := h.svc.Health(c.Request.Context()); err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return nil
}

// Create handles POST /profile
func (h *ProfileHandler) Create(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return apierrors.NewBadRequest("missing X-User-ID header", nil)
	}

	var req domain.ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}
	if err := domain.Validate(&req); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}

	out, err := h.svc.CreateProfile(c.Request.Context(), userID, req)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusCreated, out)
	return nil
}

// GetByID handles GET /profile
func (h *ProfileHandler) GetByID(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	out, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound(err.Error())
		}
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, out)
	return nil
}

// Update handles PUT /profile
// Only the owner of the profile can update it.
func (h *ProfileHandler) Update(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return apierrors.NewBadRequest("missing X-User-ID header", nil)
	}

	// Fetch existing profile to verify ownership
	existing, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound(err.Error())
		}
		return apierrors.NewInternal(err)
	}
	if existing.UserID != userID {
		return apierrors.NewForbidden("unauthorized to update this profile")
	}

	var req domain.ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}
	if err := domain.Validate(&req); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}

	out, err := h.svc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound(err.Error())
		}
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, out)
	return nil
}

// Delete handles DELETE /profile
// Only the owner of the profile can delete it.
func (h *ProfileHandler) Delete(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return apierrors.NewBadRequest("missing X-User-ID header", nil)
	}

	// Fetch existing profile to verify ownership
	existing, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound(err.Error())
		}
		return apierrors.NewInternal(err)
	}
	if existing.UserID != userID {
		return apierrors.NewForbidden("unauthorized to delete this profile")
	}

	if err := h.svc.DeleteProfile(c.Request.Context(), userID); err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound(err.Error())
		}
		return apierrors.NewInternal(err)
	}
	c.Status(http.StatusNoContent)
	return nil
}
