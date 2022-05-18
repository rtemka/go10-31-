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

func TestMux(t *testing.T) {
	h := testApi.Mux()

	t.Run("root_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com/", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		assert("http status code", http.StatusNoContent, resp.StatusCode, t)
		assert("Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)
	})

	t.Run("posts_good_method_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com/posts", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		assert("http status code", http.StatusOK, resp.StatusCode, t)
		assert("Content-Type", "application/json", resp.Header.Get("Content-Type"), t)
	})

	t.Run("posts_wrong_method_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodTrace, "http://test.com/posts", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		assert("http status code", http.StatusMethodNotAllowed, resp.StatusCode, t)
		assert("Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)

	})

	t.Run("posts_options_request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "http://test.com/posts", nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)

		resp := w.Result()

		assert("http status code", http.StatusOK, resp.StatusCode, t)
		assert("Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)

		if resp.Header.Get("Allow") == "" {
			t.Fatal("expected to receive allowed methods for resource, got nothing")
		}
	})
}

func TestHandlers(t *testing.T) {
	b, err := json.Marshal(memDb.FakePost)
	if err != nil {
		t.Fatalf("due encoding test data [%v]", err)
	}

	t.Run("test_postPostHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "http://test.com/posts", bytes.NewReader(b))
		w := httptest.NewRecorder()

		testApi.postPostHandler(w, req)

		resp := w.Result()

		assert("http status code", http.StatusCreated, resp.StatusCode, t)
		assert("Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)

	})

	t.Run("test_putPostHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "http://test.com/posts", bytes.NewReader(b))
		w := httptest.NewRecorder()

		testApi.putPostHandler(w, req)

		resp := w.Result()

		assert("http status code", http.StatusOK, resp.StatusCode, t)
		assert("Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)
	})

	t.Run("test_deletePostHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "http://test.com/posts", bytes.NewReader(b))
		w := httptest.NewRecorder()

		testApi.deletePostHandler(w, req)

		resp := w.Result()

		assert("http status code", http.StatusOK, resp.StatusCode, t)
		assert("Content-Type", "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), t)
	})

	t.Run("test_getPostsHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "http://test.com/posts", nil)
		w := httptest.NewRecorder()

		testApi.getPostsHandler(w, req)

		resp := w.Result()

		assert("http status code", http.StatusOK, resp.StatusCode, t)
		assert("Content-Type", "application/json", resp.Header.Get("Content-Type"), t)

		expected := new(bytes.Buffer)
		err := json.NewEncoder(expected).Encode(map[string]any{"data": memDb.FakeData})
		if err != nil {
			t.Fatalf("due encoding test data [%v]", err)
		}

		got, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("due reading response body [%v]", err)
		}

		if !bytes.Equal(expected.Bytes(), got) {
			t.Fatalf("expected to receive data [%s], got [%s]", string(expected.Bytes()), string(got))
		}
	})
}

func assert[T comparable](name string, expected, got T, t *testing.T) {
	if expected != got {
		t.Fatalf("expected to receive %s [%v], got [%v]", name, expected, got)
	}
}
