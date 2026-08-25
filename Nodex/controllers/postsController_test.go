package controllers_test

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"Nodex/controllers"
	"Nodex/models"

	"github.com/gin-gonic/gin"
)

// envelope mirrors the success response shape returned by the posts
// controllers, including the pagination extras emitted by FindPosts.
type envelope struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	Total      int64           `json:"total"`
	TotalPages int             `json:"total_pages"`
}

// setupRouter wires the real controllers onto a fresh gin engine, matching
// the routes registered in main.go.
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/api/posts", controllers.FindPosts)
	r.POST("/api/posts", controllers.StorePost)
	r.GET("/api/posts/:id", controllers.FindPostById)
	r.PUT("/api/posts/:id", controllers.UpdatePost)
	r.DELETE("/api/posts/:id", controllers.DeletePost)
	return r
}

// doRequest issues a request against the test router and returns the recorder.
func doRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// resetDB truncates the shared posts table (and restarts its id sequence) so
// every test starts from a clean, deterministic state. The running crudin
// server is not under concurrent load during tests, so truncating the shared
// table here is safe.
func resetDB(t *testing.T) {
	t.Helper()
	if err := models.DB.Exec("TRUNCATE TABLE posts RESTART IDENTITY").Error; err != nil {
		t.Fatalf("truncate posts: %v", err)
	}
}

func TestMain(m *testing.M) {
	// Keep request-log noise out of test output.
	gin.SetMode(gin.TestMode)

	if _, err := models.ConnectDatabase(context.Background()); err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := models.DB.Exec("TRUNCATE TABLE posts RESTART IDENTITY").Error; err != nil {
		log.Fatalf("initial truncate: %v", err)
	}

	code := m.Run()

	// Teardown: leave the shared database clean for any other consumer.
	// Done explicitly (not via defer) because os.Exit does not run defers.
	if err := models.DB.Exec("TRUNCATE TABLE posts RESTART IDENTITY").Error; err != nil {
		log.Fatalf("teardown truncate: %v", err)
	}

	os.Exit(code)
}

// TestFindPostsEmpty asserts an empty table yields an empty data slice with
// zeroed totals.
func TestFindPostsEmpty(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	w := doRequest(r, http.MethodGet, "/api/posts?page=1&limit=20", "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !env.Success {
		t.Errorf("expected success=true, got false (message=%q)", env.Message)
	}
	if env.Message != "Lists Data Posts" {
		t.Errorf("expected message %q, got %q", "Lists Data Posts", env.Message)
	}
	if env.Page != 1 {
		t.Errorf("expected page=1, got %d", env.Page)
	}
	if env.Limit != 20 {
		t.Errorf("expected limit=20, got %d", env.Limit)
	}
	if env.Total != 0 {
		t.Errorf("expected total=0, got %d", env.Total)
	}
	if env.TotalPages != 0 {
		t.Errorf("expected total_pages=0, got %d", env.TotalPages)
	}

	var posts []models.Post
	if err := json.Unmarshal(env.Data, &posts); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if len(posts) != 0 {
		t.Errorf("expected empty data slice, got %d posts", len(posts))
	}
}

// TestFindPostsInvalidPagination asserts out-of-range / non-integer params
// yield a 400 with the validation envelope.
func TestFindPostsInvalidPagination(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	cases := []string{
		"/api/posts?page=0&limit=20",   // page below minimum
		"/api/posts?page=abc&limit=20", // page non-integer
		"/api/posts?page=1&limit=0",    // limit below minimum
		"/api/posts?page=1&limit=101",  // limit above cap
		"/api/posts?page=1&limit=xyz",  // limit non-integer
	}
	for _, path := range cases {
		w := doRequest(r, http.MethodGet, path, "")
		if w.Code != http.StatusBadRequest {
			t.Errorf("%s: expected 400, got %d (%s)", path, w.Code, w.Body.String())
			continue
		}
		var env envelope
		if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
			t.Errorf("%s: unmarshal: %v", path, err)
			continue
		}
		if env.Success {
			t.Errorf("%s: expected success=false, got true", path)
		}
		if env.Message != "invalid page/limit" {
			t.Errorf("%s: expected message %q, got %q", path, "invalid page/limit", env.Message)
		}
	}
}

// TestStorePost validates a successful create returns 201 with the new post.
func TestStorePost(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	w := doRequest(r, http.MethodPost, "/api/posts", `{"title":"Test Title","content":"Test Content"}`)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !env.Success {
		t.Errorf("expected success=true, got false")
	}

	var post models.Post
	if err := json.Unmarshal(env.Data, &post); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if post.Id == 0 {
		t.Errorf("expected non-zero integer id, got %d", post.Id)
	}
	if post.Title != "Test Title" {
		t.Errorf("expected title %q, got %q", "Test Title", post.Title)
	}
}

// TestStorePostValidation ensures invalid input yields a 400 with an errors array.
func TestStorePostValidation(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	w := doRequest(r, http.MethodPost, "/api/posts", `{"title":""}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}

	var errs struct {
		Errors []json.RawMessage `json:"errors"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &errs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(errs.Errors) == 0 {
		t.Errorf("expected non-empty errors array, got %s", w.Body.String())
	}
}

// TestFindPostByIdNotFound asserts a missing id returns 404 with the not-found message.
func TestFindPostByIdNotFound(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	w := doRequest(r, http.MethodGet, "/api/posts/9999", "")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}

	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Success {
		t.Errorf("expected success=false, got true")
	}
	if env.Message != "Post not found" {
		t.Errorf("expected message %q, got %q", "Post not found", env.Message)
	}
}

// TestUpdatePost creates a post, updates it by id, and checks the title changed.
// TRUNCATE restarts the sequence, so the first created id is 1.
func TestUpdatePost(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	create := doRequest(r, http.MethodPost, "/api/posts", `{"title":"Original Title","content":"Original Content"}`)
	if create.Code != http.StatusCreated {
		t.Fatalf("seed create: expected 201, got %d: %s", create.Code, create.Body.String())
	}

	w := doRequest(r, http.MethodPut, "/api/posts/1", `{"title":"Updated Title","content":"Updated Content"}`)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !env.Success {
		t.Errorf("expected success=true, got false")
	}

	var post models.Post
	if err := json.Unmarshal(env.Data, &post); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if post.Title != "Updated Title" {
		t.Errorf("expected title %q, got %q", "Updated Title", post.Title)
	}
}

// TestUpdatePostNotFound asserts updating a missing id returns 404.
func TestUpdatePostNotFound(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	w := doRequest(r, http.MethodPut, "/api/posts/9999", `{"title":"X","content":"Y"}`)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}

	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Success {
		t.Errorf("expected success=false, got true")
	}
	if env.Message != "Post not found" {
		t.Errorf("expected message %q, got %q", "Post not found", env.Message)
	}
}

// TestDeletePost creates a post and deletes it, asserting the success message.
func TestDeletePost(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	create := doRequest(r, http.MethodPost, "/api/posts", `{"title":"To Be Deleted","content":"content"}`)
	if create.Code != http.StatusCreated {
		t.Fatalf("seed create: expected 201, got %d: %s", create.Code, create.Body.String())
	}

	w := doRequest(r, http.MethodDelete, "/api/posts/1", "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !env.Success {
		t.Errorf("expected success=true, got false")
	}
	if env.Message != "Post moved to trash" {
		t.Errorf("expected message %q, got %q", "Post moved to trash", env.Message)
	}
}

// TestDeletePostNotFound asserts deleting a missing id returns 404.
func TestDeletePostNotFound(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	w := doRequest(r, http.MethodDelete, "/api/posts/9999", "")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}

	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Success {
		t.Errorf("expected success=false, got true")
	}
	if env.Message != "Post not found" {
		t.Errorf("expected message %q, got %q", "Post not found", env.Message)
	}
}
