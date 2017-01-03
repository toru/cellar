package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/husobee/vestigo"
	"github.com/tmaesaka/cellar/config"
)

var (
	cfg     = config.NewApiConfig()
	repoId  = "project_x"
	testDir = "/tmp/cellar_test"
)

func TestIndexRepositoryHandler(t *testing.T) {
	handler := IndexRepositoryHandler(cfg)
	req, _ := http.NewRequest("GET", "/repositories", nil)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Exepected status code 200; got %d", recorder.Code)
	}

	if recorder.Body.String() != "all repositories" {
		t.Error("Unexpected response body")
	}
}

func TestShowRepositoryHandler(t *testing.T) {
	router := vestigo.NewRouter()
	router.Get("/repositories/:id", ShowRepositoryHandler(cfg))

	path := "/repositories/" + repoId
	req, _ := http.NewRequest("GET", path, nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Exepected status code 200; got %d", recorder.Code)
	}

	if recorder.Body.String() != "showing "+repoId {
		t.Error("Unexpected response body")
	}
}

func TestCreateRepositoryHandler(t *testing.T) {
	cfg.DataDir = testDir
	router := vestigo.NewRouter()
	router.Post("/repositories", CreateRepositoryHandler(cfg))

	t.Run("missing name param", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/repositories", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusBadRequest {
			t.Errorf("Exepected status code 400; got %d", recorder.Code)
		}

		err := decodeError(recorder.Body.String())

		if err.ErrorType != ErrorInvalidRequest {
			t.Errorf("Expected invalid_request_error; got %s", err.ErrorType)
		}

		if err.Message != "name parameter required" {
			t.Errorf("Unexpected response; got %s", err.Message)
		}
	})

	t.Run("name=<repoId>", func(t *testing.T) {
		params := url.Values{}
		params.Set("name", repoId)
		encodedParams := strings.NewReader(params.Encode())

		req, _ := http.NewRequest("POST", "/repositories", encodedParams)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Exepected status code 200; got %d", recorder.Code)
		}

		if recorder.Body.String() != "success" {
			t.Errorf("Expected success; got %s", recorder.Body.String())
		}
	})

	if err := os.RemoveAll(cfg.DataDir); err != nil {
		t.Error(err)
	}
}

func TestUpdateRepositoryHandler(t *testing.T) {
	router := vestigo.NewRouter()
	router.Put("/repositories/:id", UpdateRepositoryHandler(cfg))

	path := "/repositories/" + repoId
	req, _ := http.NewRequest("PUT", path, nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Exepected status code 200; got %d", recorder.Code)
	}

	if recorder.Body.String() != "updating "+repoId {
		t.Error("Unexpected response body")
	}
}

func TestDestroyRepositoryHandler(t *testing.T) {
	router := vestigo.NewRouter()
	router.Delete("/repositories/:id", DestroyRepositoryHandler(cfg))

	path := "/repositories/" + repoId
	req, _ := http.NewRequest("DELETE", path, nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Exepected status code 200; got %d", recorder.Code)
	}

	if recorder.Body.String() != "destroying "+repoId {
		t.Error("Unexpected response body")
	}
}