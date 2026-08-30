package controllers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"Nodex/controllers"
	"Nodex/models"

	"github.com/gin-gonic/gin"
)

// trashRouter registers the full set of posts routes — including restore and
// permanent-delete — so the soft-delete/trash behavior can be exercised end to end.
func trashRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/api/posts", controllers.FindPosts)
	r.POST("/api/posts", controllers.StorePost)
	// /trashed before /:id so it isn't swallowed as the id param.
	r.GET("/api/posts/trashed", controllers.FindTrashedPosts)
	r.GET("/api/posts/:id", controllers.FindPostById)
	r.PUT("/api/posts/:id", controllers.UpdatePost)
	r.DELETE("/api/posts/:id", controllers.DeletePost)
	r.PUT("/api/posts/:id/restore", controllers.RestorePost)
	r.DELETE("/api/posts/:id/permanent", controllers.DeletePermanentPost)
	return r
}

// TestSoftDeleteHidesPost verifies DELETE performs a soft-delete: the post
// disappears from direct fetch (GORM's default scope excludes trashed rows)
// and the response reports the trash envelope.
func TestSoftDeleteHidesPost(t *testing.T) {
	resetDB(t)
	r := trashRouter()

	create := doRequest(r, http.MethodPost, "/api/posts", `{"title":"soft","content":"c"}`)
	if create.Code != http.StatusCreated {
		t.Fatalf("seed create: expected 201, got %d: %s", create.Code, create.Body.String())
	}

	del := doRequest(r, http.MethodDelete, "/api/posts/1", "")
	if del.Code != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d: %s", del.Code, del.Body.String())
	}
	var env envelope
	if err := json.Unmarshal(del.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Message != "Post moved to trash" {
		t.Errorf("delete message = %q, want %q", env.Message, "Post moved to trash")
	}

	// Soft-deleted → hidden from the default scope.
	if got := doRequest(r, http.MethodGet, "/api/posts/1", ""); got.Code != http.StatusNotFound {
		t.Errorf("fetch trashed: expected 404, got %d", got.Code)
	}
}

// TestRestorePost verifies a trashed post can be restored and then fetched.
func TestRestorePost(t *testing.T) {
	resetDB(t)
	r := trashRouter()

	doRequest(r, http.MethodPost, "/api/posts", `{"title":"r","content":"c"}`)
	doRequest(r, http.MethodDelete, "/api/posts/1", "")

	rest := doRequest(r, http.MethodPut, "/api/posts/1/restore", "")
	if rest.Code != http.StatusOK {
		t.Fatalf("restore: expected 200, got %d: %s", rest.Code, rest.Body.String())
	}
	var env envelope
	if err := json.Unmarshal(rest.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Message != "Post restored" {
		t.Errorf("restore message = %q, want %q", env.Message, "Post restored")
	}

	if got := doRequest(r, http.MethodGet, "/api/posts/1", ""); got.Code != http.StatusOK {
		t.Errorf("post-restore fetch: expected 200, got %d", got.Code)
	}
}

// TestDeletePermanently removes a trashed post and verifies it is gone,
// even via an unscoped lookup.
func TestDeletePermanently(t *testing.T) {
	resetDB(t)
	r := trashRouter()

	doRequest(r, http.MethodPost, "/api/posts", `{"title":"p","content":"c"}`)
	doRequest(r, http.MethodDelete, "/api/posts/1", "")

	perm := doRequest(r, http.MethodDelete, "/api/posts/1/permanent", "")
	if perm.Code != http.StatusOK {
		t.Fatalf("permanent delete: expected 200, got %d: %s", perm.Code, perm.Body.String())
	}
	var env envelope
	if err := json.Unmarshal(perm.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Message != "Post permanently deleted" {
		t.Errorf("perm message = %q, want %q", env.Message, "Post permanently deleted")
	}

	if got := doRequest(r, http.MethodGet, "/api/posts/1", ""); got.Code != http.StatusNotFound {
		t.Errorf("post-permanent fetch: expected 404, got %d", got.Code)
	}

	// And hard-gone from the database, even unscoped.
	var p models.Post
	if err := models.DB.Unscoped().Where("id = ?", 1).First(&p).Error; err == nil {
		t.Error("expected post to be hard-deleted from DB, but found it")
	}
}

// TestDeletePermanentOnLiveReturns404 ensures permanent-delete only targets
// already-trashed posts: a live post yields 404 (never a hard-delete-by-mistake).
func TestDeletePermanentOnLiveReturns404(t *testing.T) {
	resetDB(t)
	r := trashRouter()

	doRequest(r, http.MethodPost, "/api/posts", `{"title":"live","content":"c"}`)

	perm := doRequest(r, http.MethodDelete, "/api/posts/1/permanent", "")
	if perm.Code != http.StatusNotFound {
		t.Errorf("permanent-delete on live post: expected 404, got %d: %s", perm.Code, perm.Body.String())
	}
}

// TestTrashedListing verifies FindTrashedPosts returns only soft-deleted rows,
// that live (non-trashed) listings exclude them, and that pagination works on
// the trashed set.
func TestTrashedListing(t *testing.T) {
	resetDB(t)
	r := trashRouter()

	// Seed two posts and soft-delete both (ids 1 and 2).
	if w := doRequest(r, http.MethodPost, "/api/posts", `{"title":"a","content":"c"}`); w.Code != http.StatusCreated {
		t.Fatalf("seed create 1: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	if w := doRequest(r, http.MethodPost, "/api/posts", `{"title":"b","content":"c"}`); w.Code != http.StatusCreated {
		t.Fatalf("seed create 2: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	if w := doRequest(r, http.MethodDelete, "/api/posts/1", ""); w.Code != http.StatusOK {
		t.Fatalf("soft delete 1: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if w := doRequest(r, http.MethodDelete, "/api/posts/2", ""); w.Code != http.StatusOK {
		t.Fatalf("soft delete 2: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Only trashed posts are returned.
	w := doRequest(r, http.MethodGet, "/api/posts/trashed", "")
	if w.Code != http.StatusOK {
		t.Fatalf("trashed list: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Message != "Trashed Posts" {
		t.Errorf("message = %q, want %q", env.Message, "Trashed Posts")
	}
	if env.Total != 2 {
		t.Errorf("total = %d, want 2", env.Total)
	}
	if env.Page != 1 || env.Limit != 20 {
		t.Errorf("pagination = page %d limit %d, want 1/20", env.Page, env.Limit)
	}
	if env.TotalPages != 1 {
		t.Errorf("total_pages = %d, want 1", env.TotalPages)
	}
	var trashed []models.Post
	if err := json.Unmarshal(env.Data, &trashed); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if len(trashed) != 2 {
		t.Fatalf("expected 2 trashed posts, got %d", len(trashed))
	}

	// The normal (non-trashed) listing must exclude the soft-deleted rows.
	live := doRequest(r, http.MethodGet, "/api/posts", "")
	var liveEnv envelope
	if err := json.Unmarshal(live.Body.Bytes(), &liveEnv); err != nil {
		t.Fatalf("unmarshal live: %v", err)
	}
	if liveEnv.Total != 0 {
		t.Errorf("non-trashed total = %d, want 0", liveEnv.Total)
	}

	// Pagination on the trashed set: 1 per page across 2 total rows.
	paged := doRequest(r, http.MethodGet, "/api/posts/trashed?page=1&limit=1", "")
	if paged.Code != http.StatusOK {
		t.Fatalf("paged trashed: expected 200, got %d: %s", paged.Code, paged.Body.String())
	}
	var pagedEnv envelope
	if err := json.Unmarshal(paged.Body.Bytes(), &pagedEnv); err != nil {
		t.Fatalf("unmarshal paged: %v", err)
	}
	if pagedEnv.Total != 2 {
		t.Errorf("paged total = %d, want 2", pagedEnv.Total)
	}
	if pagedEnv.TotalPages != 2 {
		t.Errorf("paged total_pages = %d, want 2", pagedEnv.TotalPages)
	}
	var pagedPosts []models.Post
	if err := json.Unmarshal(pagedEnv.Data, &pagedPosts); err != nil {
		t.Fatalf("unmarshal paged data: %v", err)
	}
	if len(pagedPosts) != 1 {
		t.Errorf("page 1 of 2: expected 1 post, got %d", len(pagedPosts))
	}
}

// TestTrashedStatusFilter verifies the optional ?status= filter applies within
// the trashed scope.
func TestTrashedStatusFilter(t *testing.T) {
	resetDB(t)
	r := trashRouter()

	doRequest(r, http.MethodPost, "/api/posts", `{"title":"p","content":"c","status":"published"}`)
	doRequest(r, http.MethodPost, "/api/posts", `{"title":"d","content":"c","status":"draft"}`)
	doRequest(r, http.MethodDelete, "/api/posts/1", "")
	doRequest(r, http.MethodDelete, "/api/posts/2", "")

	wd := doRequest(r, http.MethodGet, "/api/posts/trashed?status=draft", "")
	var envd envelope
	if err := json.Unmarshal(wd.Body.Bytes(), &envd); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if wd.Code != http.StatusOK || envd.Total != 1 {
		t.Errorf("draft trashed: code %d total %d, want 200/1", wd.Code, envd.Total)
	}

	wp := doRequest(r, http.MethodGet, "/api/posts/trashed?status=published", "")
	var envp envelope
	if err := json.Unmarshal(wp.Body.Bytes(), &envp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if wp.Code != http.StatusOK || envp.Total != 1 {
		t.Errorf("published trashed: code %d total %d, want 200/1", wp.Code, envp.Total)
	}
}
