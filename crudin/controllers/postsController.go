package controllers

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"crudin/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// ValidatePostInput is the request DTO for creating/updating a Post.
// Status is optional: "published" | "draft" | "archived" (§34-35 Archive).
// Folder/Tags are optional Nodex organization fields (§11-12).
type ValidatePostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Status  string `json:"status"`
	Folder  string `json:"folder"`
	Tags    string `json:"tags"`
}

// isValidPostStatus reports whether s is an accepted post status.
func isValidPostStatus(s string) bool {
	return s == "published" || s == "draft" || s == "archived"
}

// ErrorMsg describes a single field-level validation failure.
type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// GetErrorMsg maps a validator error tag to a user-friendly message.
func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	}
	return "Unknown error"
}

// jsonError sends a failure response using the same envelope as successes,
// so the frontend always reads success/message/data consistently.
func jsonError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
		"data":    nil,
	})
}

// FindPosts lists posts with simple page/limit pagination and stable
// newest-first ordering. trashed posts are excluded automatically by GORM's
// soft-delete default scope (Post has a DeletedAt field).
func FindPosts(c *gin.Context) {
	// TODO: switch to cursor-based pagination (seek by created_at+id from the
	// last row) for write-heavy workloads where offset row-shifting matters.
	page, limit, ok := parsePagination(c)
	if !ok {
		jsonError(c, http.StatusBadRequest, "invalid page/limit")
		return
	}

	var posts []models.Post
	query := models.DB.Model(&models.Post{}).Order("created_at DESC, id DESC")

	// Status filter: All Notes (no ?status) shows active only (exclude
	// archived per Nodex §34), while ?status=archived/published/draft
	// returns that single status. Trashed rows are already excluded.
	if status := c.Query("status"); status != "" {
		if !isValidPostStatus(status) {
			jsonError(c, http.StatusBadRequest, "invalid status")
			return
		}
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status != ?", "archived")
	}

	// Nodex §11-12: optional folder/tags filters (Folder exact, Tags LIKE).
	if folder := c.Query("folder"); folder != "" {
		query = query.Where("folder = ?", folder)
	}
	if tags := c.Query("tags"); tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Lists Data Posts",
		"data":        posts,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	})
}

// FindTrashedPosts lists only soft-deleted (trashed) posts. Because GORM's
// default scope hides DeletedAt rows, the query is run Unscoped and gated
// on deleted_at IS NOT NULL so live posts never leak in. It supports the
// same optional ?status= filter and page/limit pagination (capped at 100,
// newest trashed-first ordering by deleted_at DESC, id DESC).
func FindTrashedPosts(c *gin.Context) {
	page, limit, ok := parsePagination(c)
	if !ok {
		jsonError(c, http.StatusBadRequest, "invalid page/limit")
		return
	}

	var posts []models.Post
	query := models.DB.Unscoped().
		Model(&models.Post{}).
		Where("deleted_at IS NOT NULL").
		Order("deleted_at DESC, id DESC")

	// Optional status filter alongside the trash scope.
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if folder := c.Query("folder"); folder != "" {
		query = query.Where("folder = ?", folder)
	}
	if tags := c.Query("tags"); tags != "" {
		query = query.Where("tags LIKE ?", "%"+tags+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&posts).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Trashed Posts",
		"data":        posts,
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
	})
}

// parsePagination reads and validates the ?page and ?limit query params,
// applying defaults (page=1, limit=20) and bounds (page>=1, 1<=limit<=100).
// It returns ok=false when a value is out-of-range or non-integer.
func parsePagination(c *gin.Context) (page, limit int, ok bool) {
	pageStr := c.Query("page")
	if pageStr == "" {
		page = 1
	} else {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			return 0, 0, false
		}
		page = p
	}

	limitStr := c.Query("limit")
	if limitStr == "" {
		limit = 20
	} else {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 || l > 100 {
			return 0, 0, false
		}
		limit = l
	}

	return page, limit, true
}

// StorePost creates a new post from the validated request body.
func StorePost(c *gin.Context) {
	// Validate input.
	var input ValidatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}

	// Resolve status: absent means "published"; otherwise it must be a
	// known value, otherwise the request is rejected with 400.
	status := input.Status
	if status == "" {
		status = "published"
	}
	if !isValidPostStatus(status) {
		jsonError(c, http.StatusBadRequest, "invalid status")
		return
	}

	// Create the post, defaulting absent status to "published".
	post := models.Post{
		Title:   input.Title,
		Content: input.Content,
		Status:  status,
		Folder:  input.Folder,
		Tags:    input.Tags,
	}
	if err := models.DB.Create(&post).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to create post")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Post Created Successfully",
		"data":    post,
	})
}

// FindPostById returns a single post by id, distinguishing a missing
// record (404) from a real database error (500). trashed posts are hidden
// by GORM's default scope and read as "not found".
func FindPostById(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonError(c, http.StatusNotFound, "Post not found")
		} else {
			jsonError(c, http.StatusInternalServerError, "Failed to fetch post")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail Data Post By ID : " + c.Param("id"),
		"data":    post,
	})
}

// UpdatePost modifies an existing post.
func UpdatePost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonError(c, http.StatusNotFound, "Post not found")
		} else {
			jsonError(c, http.StatusInternalServerError, "Failed to fetch post")
		}
		return
	}

	// Validate input.
	var input ValidatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}

	// Validate optional status. When provided it must be a known value;
	// when absent it is left untouched (GORM's struct Updates also skips
	// the zero-value string, so an empty status never overwrites the stored value).
	if input.Status != "" && !isValidPostStatus(input.Status) {
		jsonError(c, http.StatusBadRequest, "invalid status")
		return
	}

	// Apply updates. Status is only written when explicitly provided.
	if err := models.DB.Model(&post).Updates(input).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to update post")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post Updated Successfully",
		"data":    post,
	})
}

// DeletePost soft-deletes a post, moving it to Trash. GORM issues an
// UPDATE setting deleted_at because Post carries a DeletedAt field, so the
// row is preserved and excluded from normal listings.
func DeletePost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonError(c, http.StatusNotFound, "Post not found")
		} else {
			jsonError(c, http.StatusInternalServerError, "Failed to fetch post")
		}
		return
	}

	if err := models.DB.Delete(&post).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post moved to trash",
	})
}

// RestorePost un-soft-deletes a trashed post. Unscoped() lets the lookup see
// soft-deleted rows; the restore clears deleted_at back to NULL.
func RestorePost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Unscoped().Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonError(c, http.StatusNotFound, "Post not found")
		} else {
			jsonError(c, http.StatusInternalServerError, "Failed to restore post")
		}
		return
	}

	if err := models.DB.Unscoped().Model(&post).Update("deleted_at", nil).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to restore post")
		return
	}

	post.DeletedAt = gorm.DeletedAt{} // reflect restored state in the response payload
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post restored",
		"data":    post,
	})
}

// DeletePermanentPost permanently removes an already-trashed post. It gates
// on deleted_at IS NOT NULL so a live (non-trashed) post yields 404 rather
// than being silently hard-deleted.
func DeletePermanentPost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", c.Param("id")).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonError(c, http.StatusNotFound, "Post not found")
		} else {
			jsonError(c, http.StatusInternalServerError, "Failed to permanently delete post")
		}
		return
	}

	if err := models.DB.Unscoped().Delete(&post).Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to permanently delete post")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post permanently deleted",
	})
}

// ArchivePost moves a post to Archive by setting status=archived.
func ArchivePost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonError(c, http.StatusNotFound, "Post not found")
		} else {
			jsonError(c, http.StatusInternalServerError, "Failed to fetch post")
		}
		return
	}
	if err := models.DB.Model(&post).Update("status", "archived").Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to archive post")
		return
	}
	post.Status = "archived"
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post archived",
		"data":    post,
	})
}

// UnarchivePost restores an archived post to published.
func UnarchivePost(c *gin.Context) {
	var post models.Post
	if err := models.DB.Where("id = ?", c.Param("id")).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			jsonError(c, http.StatusNotFound, "Post not found")
		} else {
			jsonError(c, http.StatusInternalServerError, "Failed to fetch post")
		}
		return
	}
	if err := models.DB.Model(&post).Update("status", "published").Error; err != nil {
		jsonError(c, http.StatusInternalServerError, "Failed to unarchive post")
		return
	}
	post.Status = "published"
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Post unarchived",
		"data":    post,
	})
}
