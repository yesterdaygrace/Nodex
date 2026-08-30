package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Nodex/models"
)

// mustUnmarshal envelopes the body or fatals.
func mustEnvelope(t *testing.T, w *httptest.ResponseRecorder) envelope {
	t.Helper()
	var env envelope
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v (body=%s)", err, w.Body.String())
	}
	return env
}

// TestCreateDraftAndFilter verifies the optional status field, defaulting to
// "published" when absent, and the ?status= filter on the live listing.
func TestCreateDraftAndFilter(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	// Explicit draft.
	if w := doRequest(r, http.MethodPost, "/api/posts", `{"title":"Draft","content":"d","status":"draft"}`); w.Code != http.StatusCreated {
		t.Fatalf("create draft: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	// Absent status -> defaults to published.
	if w := doRequest(r, http.MethodPost, "/api/posts", `{"title":"Pub","content":"p"}`); w.Code != http.StatusCreated {
		t.Fatalf("create published: expected 201, got %d: %s", w.Code, w.Body.String())
	}
	// Published explicitly.
	if w := doRequest(r, http.MethodPost, "/api/posts", `{"title":"Pub2","content":"p","status":"published"}`); w.Code != http.StatusCreated {
		t.Fatalf("create published explicit: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// ?status=draft -> only the draft.
	wd := doRequest(r, http.MethodGet, "/api/posts?status=draft", "")
	envd := mustEnvelope(t, wd)
	if envd.Total != 1 {
		t.Errorf("draft filter total = %d, want 1", envd.Total)
	}
	var dposts []models.Post
	if err := json.Unmarshal(envd.Data, &dposts); err != nil {
		t.Fatalf("unmarshal draft data: %v", err)
	}
	if len(dposts) != 1 || dposts[0].Status != "draft" {
		t.Errorf("draft filter: got %+v, want exactly 1 draft", dposts)
	}

	// ?status=published -> only published (2 of them).
	wp := doRequest(r, http.MethodGet, "/api/posts?status=published", "")
	envp := mustEnvelope(t, wp)
	if envp.Total != 2 {
		t.Errorf("published filter total = %d, want 2", envp.Total)
	}
	var pposts []models.Post
	if err := json.Unmarshal(envp.Data, &pposts); err != nil {
		t.Fatalf("unmarshal published data: %v", err)
	}
	for _, p := range pposts {
		if p.Status != "published" {
			t.Errorf("published filter: got status %q, want published", p.Status)
		}
	}

	// No filter -> all non-trashed posts (3).
	wa := doRequest(r, http.MethodGet, "/api/posts", "")
	enva := mustEnvelope(t, wa)
	if enva.Total != 3 {
		t.Errorf("no filter total = %d, want 3", enva.Total)
	}
}

// TestUpdateStatus verifies a post's status can be flipped between draft and
// published, and that omitting status on update does NOT overwrite the stored
// value.
func TestUpdateStatus(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	if w := doRequest(r, http.MethodPost, "/api/posts", `{"title":"t","content":"c"}`); w.Code != http.StatusCreated {
		t.Fatalf("seed: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	// Flip to draft.
	w := doRequest(r, http.MethodPut, "/api/posts/1", `{"title":"t","content":"c","status":"draft"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("update to draft: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	env := mustEnvelope(t, w)
	var post models.Post
	if err := json.Unmarshal(env.Data, &post); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if post.Status != "draft" {
		t.Errorf("after update: status = %q, want draft", post.Status)
	}
	if post.Title != "t" {
		t.Errorf("after update: title = %q, want t", post.Title)
	}

	// Filter reflects the draft.
	wd := doRequest(r, http.MethodGet, "/api/posts?status=draft", "")
	envd := mustEnvelope(t, wd)
	if envd.Total != 1 {
		t.Errorf("draft filter total = %d, want 1", envd.Total)
	}

	// Flip back to published.
	w2 := doRequest(r, http.MethodPut, "/api/posts/1", `{"title":"t","content":"c","status":"published"}`)
	if w2.Code != http.StatusOK {
		t.Fatalf("update to published: expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
	env2 := mustEnvelope(t, w2)
	var post2 models.Post
	if err := json.Unmarshal(env2.Data, &post2); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if post2.Status != "published" {
		t.Errorf("after update: status = %q, want published", post2.Status)
	}

	// Omit status on update -> must NOT overwrite existing value (stays published).
	w3 := doRequest(r, http.MethodPut, "/api/posts/1", `{"title":"title-2","content":"c2"}`)
	if w3.Code != http.StatusOK {
		t.Fatalf("update without status: expected 200, got %d: %s", w3.Code, w3.Body.String())
	}
	env3 := mustEnvelope(t, w3)
	var post3 models.Post
	if err := json.Unmarshal(env3.Data, &post3); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if post3.Status != "published" {
		t.Errorf("status without body field = %q, want published (not overwritten)", post3.Status)
	}
	if post3.Title != "title-2" {
		t.Errorf("title = %q, want title-2", post3.Title)
	}

	// Invalid status on update is rejected (400), body untouched.
	wb := doRequest(r, http.MethodPut, "/api/posts/1", `{"title":"t","content":"c","status":"bogus"}`)
	if wb.Code != http.StatusBadRequest {
		t.Fatalf("update invalid status: expected 400, got %d: %s", wb.Code, wb.Body.String())
	}
	envb := mustEnvelope(t, wb)
	if envb.Success {
		t.Errorf("expected success=false for invalid update status")
	}
	if envb.Message != "invalid status" {
		t.Errorf("message = %q, want %q", envb.Message, "invalid status")
	}
}

// TestInvalidStatus asserts an unknown status on create is rejected with 400
// "invalid status" and that no row is written.
func TestInvalidStatus(t *testing.T) {
	resetDB(t)
	r := setupRouter()

	w := doRequest(r, http.MethodPost, "/api/posts", `{"title":"t","content":"c","status":"bogus"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
	env := mustEnvelope(t, w)
	if env.Success {
		t.Errorf("expected success=false, got true")
	}
	if env.Message != "invalid status" {
		t.Errorf("message = %q, want %q", env.Message, "invalid status")
	}

	// No record should have been created by the rejected request.
	wa := doRequest(r, http.MethodGet, "/api/posts", "")
	enva := mustEnvelope(t, wa)
	if enva.Total != 0 {
		t.Errorf("expected 0 posts after rejected create, got %d", enva.Total)
	}
}
