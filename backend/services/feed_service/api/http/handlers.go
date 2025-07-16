package http

import (
	"net/http"

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
		grp.GET("/health", h.Health)
		// Publications
		grp.POST("/publications", h.CreatePublication)
		grp.GET("/publications", h.ListPublications)
		grp.GET("/publications/:id", h.GetPublication)
		grp.PUT("/publications/:id", h.UpdatePublication)
		grp.DELETE("/publications/:id", h.DeletePublication)

		// Comments
		grp.POST("/comments", h.CreateComment)
		grp.GET("/comments", h.ListComments)
		grp.GET("/comments/:id", h.GetComment)
		grp.PUT("/comments/:id", h.UpdateComment)
		grp.DELETE("/comments/:id", h.DeleteComment)
		grp.GET("/user/publications", h.ListUserPublications)
	}
}

// Health endpoint
func (h *FeedHandler) Health(c *gin.Context) {
	if err := h.svc.Health(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db down"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// CreatePublication handles POST /feed/publications
func (h *FeedHandler) CreatePublication(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}

	var req domain.PublicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.CreatePublication(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

// GetPublication handles GET /feed/publications/:id
func (h *FeedHandler) GetPublication(c *gin.Context) {
	id := c.Param("id")
	out, err := h.svc.GetPublication(c.Request.Context(), id)
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

// ListPublications handles GET /feed/publications
func (h *FeedHandler) ListPublications(c *gin.Context) {
	list, err := h.svc.ListPublications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *FeedHandler) ListUserPublications(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	// call usecase to get all publications for the user
	posts, err := h.svc.ListPublicationsByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// UpdatePublication handles PUT /feed/publications/:id
// Only the owner of the publication can update it.
func (h *FeedHandler) UpdatePublication(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}
	id := c.Param("id")

	// Fetch existing publication to verify ownership
	existing, err := h.svc.GetPublication(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if existing.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to update this publication"})
		return
	}

	var req domain.PublicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.svc.UpdatePublication(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

// DeletePublication handles DELETE /feed/publications/:id
// Only the owner of the publication can delete it.
func (h *FeedHandler) DeletePublication(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}
	id := c.Param("id")

	// Fetch existing publication to verify ownership
	existing, err := h.svc.GetPublication(c.Request.Context(), id)
	if err != nil {
		if err == usecases.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if existing.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to delete this publication"})
		return
	}

	if err := h.svc.DeletePublication(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// CreateComment handles POST /feed/comments
func (h *FeedHandler) CreateComment(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}

	var req domain.PostCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.CreateComment(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

// ListComments handles GET /feed/comments?post_id=...
func (h *FeedHandler) ListComments(c *gin.Context) {
	postID := c.Query("post_id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing post_id query parameter"})
		return
	}
	list, err := h.svc.ListComments(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// GetComment handles GET /feed/comments/:id
func (h *FeedHandler) GetComment(c *gin.Context) {
	id := c.Param("id")
	out, err := h.svc.GetComment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

// UpdateComment handles PUT /feed/comments/:id
// Only the owner of the comment can update it.
func (h *FeedHandler) UpdateComment(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}
	id := c.Param("id")

	// Fetch existing comment to verify ownership
	existing, err := h.svc.GetComment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to update this comment"})
		return
	}

	var req domain.PutCommentRequest
	req.CommentID = id
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.svc.UpdateComment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

// DeleteComment handles DELETE /feed/comments/:id
// Only the owner of the comment can delete it.
func (h *FeedHandler) DeleteComment(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing X-User-ID header"})
		return
	}
	id := c.Param("id")

	// Fetch existing comment to verify ownership
	existing, err := h.svc.GetComment(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to delete this comment"})
		return
	}

	if err := h.svc.DeleteComment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
