package api

import (
	memDb "GoNews/pkg/storage/memdb"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testApi *Api

func TestMain(m *testing.M) {

	// создаем API сервера
	l := log.New(os.Stderr, "[GoNews test server]\t->\t", log.LstdFlags|log.Lmsgprefix)
	testApi = New(memDb.New(), l)

	os.Exit(m.Run())
}

func TestApi_ServeHTTP(t *testing.T) {
	h := testApi.Mux()

	t.Run("root_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		assert("http status code", http.StatusNoContent, resp.StatusCode, t)
		assert("Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)
	})

	t.Run("good_method_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com/posts", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		assert("http status code", http.StatusOK, resp.StatusCode, t)
		assert("Content-Type", "application/json", resp.Header.Get("Content-Type"), t)
	})

	t.Run("wrong_method_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodTrace, "http://test.com/posts", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		assert("http status code", http.StatusMethodNotAllowed, resp.StatusCode, t)
		assert("Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)

	})

	t.Run("options_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "http://test.com/posts", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		assert("http status code", http.StatusOK, resp.StatusCode, t)
		assert("Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)

		if resp.Header.Get("Allow") == "" {
			t.Fatal("want to receive allowed methods for resource, got nothing")
		}
	})
}

func TestApi_getPostsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://test.com/posts", nil)
	w := httptest.NewRecorder()

	testApi.getPostsHandler(w, req)

	resp := w.Result()

	assert("api.getPostsHandler() http status code", http.StatusOK, resp.StatusCode, t)
	assert("api.getPostsHandler() Content-Type", "application/json", resp.Header.Get("Content-Type"), t)

	want := new(bytes.Buffer)
	err := json.NewEncoder(want).Encode(map[string]any{"data": memDb.FakeData})
	if err != nil {
		t.Fatalf("api.getPostsHandler() due encoding test data = %v", err)
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("api.getPostsHandler() due reading response body = %v", err)
	}

	if !bytes.Equal(want.Bytes(), got) {
		t.Fatalf("api.getPostsHandler() = %s, want %s", string(got), string(want.Bytes()))
	}
}

func TestApi_postPostHandler(t *testing.T) {
	b, err := json.Marshal(memDb.FakePost)
	if err != nil {
		t.Fatalf("api.postPostHandler() due encoding test data %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "http://test.com/posts", bytes.NewReader(b))
	w := httptest.NewRecorder()

	testApi.postPostHandler(w, req)

	resp := w.Result()

	assert("api.postPostHandler() http status code", http.StatusCreated, resp.StatusCode, t)
	assert("api.postPostHandler() Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)
}

func TestApi_putPostHandler(t *testing.T) {
	b, err := json.Marshal(memDb.FakePost)
	if err != nil {
		t.Fatalf("api.putPostHandler() due encoding test data %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "http://test.com/posts", bytes.NewReader(b))
	w := httptest.NewRecorder()

	testApi.putPostHandler(w, req)

	resp := w.Result()

	assert("api.putPostHandler() http status code", http.StatusOK, resp.StatusCode, t)
	assert("api.putPostHandler() Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)
}

func TestApi_deletePostHandler(t *testing.T) {
	b, err := json.Marshal(memDb.FakePost)
	if err != nil {
		t.Fatalf("api.deletePostHandler() due encoding test data %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "http://test.com/posts", bytes.NewReader(b))
	w := httptest.NewRecorder()

	testApi.deletePostHandler(w, req)

	resp := w.Result()

	assert("api.deletePostHandler() http status code", http.StatusOK, resp.StatusCode, t)
	assert("api.deletePostHandler() Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)
}

func assert[T comparable](name string, want, got T, t *testing.T) {
	if got != want {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}
