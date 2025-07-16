package http

import (
	"net/http"

	"feed_service/api/http/apierrors"
	"feed_service/api/http/middleware"
	"feed_service/domain"
	"feed_service/usecases"

	"github.com/gin-gonic/gin"
)

// FeedHandler handles HTTP requests for publications and comments
// and delegates to the FeedService business logic.
type FeedHandler struct {
	svc usecases.FeedService
}

// NewFeedHandler constructs a new FeedHandler
func NewFeedHandler(svc usecases.FeedService) *FeedHandler {
	return &FeedHandler{svc: svc}
}

// RegisterRoutes registers feed routes on the Gin engine
func (h *FeedHandler) RegisterRoutes(r *gin.Engine) {
	grp := r.Group("/feed")
	{
		grp.GET("/health", middleware.ErrorHandlerMiddleware(h.Health))
		// Publications
		grp.POST("/publications", middleware.ErrorHandlerMiddleware(h.CreatePublication))
		grp.GET("/publications", middleware.ErrorHandlerMiddleware(h.ListPublications))
		grp.GET("/publications/:id", middleware.ErrorHandlerMiddleware(h.GetPublication))
		grp.PUT("/publications/:id", middleware.ErrorHandlerMiddleware(h.UpdatePublication))
		grp.DELETE("/publications/:id", middleware.ErrorHandlerMiddleware(h.DeletePublication))

		// Comments
		grp.POST("/comments", middleware.ErrorHandlerMiddleware(h.CreateComment))
		grp.GET("/comments", middleware.ErrorHandlerMiddleware(h.ListComments))
		grp.GET("/comments/:id", middleware.ErrorHandlerMiddleware(h.GetComment))
		grp.PUT("/comments/:id", middleware.ErrorHandlerMiddleware(h.UpdateComment))
		grp.DELETE("/comments/:id", middleware.ErrorHandlerMiddleware(h.DeleteComment))
		grp.GET("/user/publications", middleware.ErrorHandlerMiddleware(h.ListUserPublications))
	}
}

// Health endpoint
func (h *FeedHandler) Health(c *gin.Context) error {
	if err := h.svc.Health(c.Request.Context()); err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return nil
}

// CreatePublication handles POST /feed/publications
func (h *FeedHandler) CreatePublication(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return apierrors.NewBadRequest("missing X-User-ID header", nil)
	}

	var req domain.PublicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}
	if err := req.Validate(); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}

	out, err := h.svc.CreatePublication(c.Request.Context(), userID, req)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusCreated, out)
	return nil
}

// GetPublication handles GET /feed/publications/:id
func (h *FeedHandler) GetPublication(c *gin.Context) error {
	id := c.Param("id")
	out, err := h.svc.GetPublication(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound(err.Error())
		}
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, out)
	return nil
}

// ListPublications handles GET /feed/publications
func (h *FeedHandler) ListPublications(c *gin.Context) error {
	list, err := h.svc.ListPublications(c.Request.Context())
	if err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, list)
	return nil
}

func (h *FeedHandler) ListUserPublications(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	// call usecase to get all publications for the user
	posts, err := h.svc.ListPublicationsByUser(c.Request.Context(), userID)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, posts)
	return nil
}

// UpdatePublication handles PUT /feed/publications/:id
// Only the owner of the publication can update it.
func (h *FeedHandler) UpdatePublication(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return apierrors.NewBadRequest("missing X-User-ID header", nil)
	}
	id := c.Param("id")

	// Fetch existing publication to verify ownership
	existing, err := h.svc.GetPublication(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound(err.Error())
		}
		return apierrors.NewInternal(err)
	}
	if existing.UserID != userID {
		return apierrors.NewForbidden("unauthorized to update this publication")
	}

	var req domain.PublicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}
	if err := req.Validate(); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}
	out, err := h.svc.UpdatePublication(c.Request.Context(), id, req)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, out)
	return nil
}

// DeletePublication handles DELETE /feed/publications/:id
// Only the owner of the publication can delete it.
func (h *FeedHandler) DeletePublication(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return apierrors.NewBadRequest("missing X-User-ID header", nil)
	}
	id := c.Param("id")

	// Fetch existing publication to verify ownership
	existing, err := h.svc.GetPublication(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrNotFound {
			return apierrors.NewNotFound(err.Error())
		}
		return apierrors.NewInternal(err)
	}
	if existing.UserID != userID {
		return apierrors.NewForbidden("unauthorized to delete this publication")
	}

	if err := h.svc.DeletePublication(c.Request.Context(), id); err != nil {
		return apierrors.NewInternal(err)
	}
	c.Status(http.StatusNoContent)
	return nil
}

// CreateComment handles POST /feed/comments
func (h *FeedHandler) CreateComment(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return apierrors.NewBadRequest("missing X-User-ID header", nil)
	}

	var req domain.PostCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}
	if err := req.Validate(); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}

	out, err := h.svc.CreateComment(c.Request.Context(), userID, req)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusCreated, out)
	return nil
}

// ListComments handles GET /feed/comments?post_id=...
func (h *FeedHandler) ListComments(c *gin.Context) error {
	postID := c.Query("post_id")
	if postID == "" {
		return apierrors.NewBadRequest("missing post_id query parameter", nil)
	}
	list, err := h.svc.ListComments(c.Request.Context(), postID)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, list)
	return nil
}

// GetComment handles GET /feed/comments/:id
func (h *FeedHandler) GetComment(c *gin.Context) error {
	id := c.Param("id")
	out, err := h.svc.GetComment(c.Request.Context(), id)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, out)
	return nil
}

// UpdateComment handles PUT /feed/comments/:id
// Only the owner of the comment can update it.
func (h *FeedHandler) UpdateComment(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return apierrors.NewBadRequest("missing X-User-ID header", nil)
	}
	id := c.Param("id")

	// Fetch existing comment to verify ownership
	existing, err := h.svc.GetComment(c.Request.Context(), id)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	if existing.UserID != userID {
		return apierrors.NewForbidden("unauthorized to update this comment")
	}

	var req domain.PutCommentRequest
	req.CommentID = id
	if err := c.ShouldBindJSON(&req); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}
	if err := req.Validate(); err != nil {
		return apierrors.NewBadRequest(err.Error(), err)
	}
	out, err := h.svc.UpdateComment(c.Request.Context(), req)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	c.JSON(http.StatusOK, out)
	return nil
}

// DeleteComment handles DELETE /feed/comments/:id
// Only the owner of the comment can delete it.
func (h *FeedHandler) DeleteComment(c *gin.Context) error {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		return apierrors.NewBadRequest("missing X-User-ID header", nil)
	}
	id := c.Param("id")

	// Fetch existing comment to verify ownership
	existing, err := h.svc.GetComment(c.Request.Context(), id)
	if err != nil {
		return apierrors.NewInternal(err)
	}
	if existing.UserID != userID {
		return apierrors.NewForbidden("unauthorized to delete this comment")
	}

	if err := h.svc.DeleteComment(c.Request.Context(), id); err != nil {
		return apierrors.NewInternal(err)
	}
	c.Status(http.StatusNoContent)
	return nil
}
