package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gui-moreira/golang-rest-introduction/user"
)

func TestSaveUser(t *testing.T) {
	repo = user.InmemUserRepo{Users: make(map[int]user.User)}
	router := configureRoutes()

	t.Run("should save user correctly", func(t *testing.T) {
		requestBody := `{
			"name": "Marquinhos",
			"age": 30
		}`

		req, _ := http.NewRequest("POST", "/users", strings.NewReader(requestBody))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		wantStatus := http.StatusOK
		gotStatus := resp.Code

		if wantStatus != gotStatus {
			t.Errorf("Want status: %d, Got status: %d", wantStatus, gotStatus)
		}

		wantBody := `{"id":1}`
		gotBody := resp.Body.String()

		if wantBody != gotBody {
			t.Errorf("Want body: %s, Got body: %s", wantBody, gotBody)
		}
	})

	t.Run("should return 400 when there is validation error", func(t *testing.T) {
		requestBody := `{
			"name": "",
			"age": 30
		}`

		req, _ := http.NewRequest("POST", "/users", strings.NewReader(requestBody))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		wantStatus := http.StatusBadRequest
		gotStatus := resp.Code

		if wantStatus != gotStatus {
			t.Errorf("Want status: %d, Got status: %d", wantStatus, gotStatus)
		}

		wantBody := `{"message":"Invalid user"}`
		gotBody := resp.Body.String()

		if wantBody != gotBody {
			t.Errorf("Want body: %s, Got body: %s", wantBody, gotBody)
		}
	})
}

func TestGetUser(t *testing.T) {
	router := configureRoutes()
	repo = user.InmemUserRepo{Users: map[int]user.User{
		1: user.User{ID: 1, Name: "Bruno", Age: 28},
	}}

	t.Run("should return existing user correctly", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/1", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		wantStatus := http.StatusOK
		gotStatus := resp.Code

		if wantStatus != gotStatus {
			t.Errorf("Want status: %d, Got status: %d", wantStatus, gotStatus)
		}

		wantBody := `{"id":1,"name":"Bruno","age":28}`
		gotBody := resp.Body.String()

		if wantBody != gotBody {
			t.Errorf("Want body: %s, Got body: %s", wantBody, gotBody)
		}
	})

	t.Run("should return 404 when user was not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/2", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		wantStatus := http.StatusNotFound
		gotStatus := resp.Code

		if wantStatus != gotStatus {
			t.Errorf("Want status: %d, Got status: %d", wantStatus, gotStatus)
		}

		wantBody := `{"message":"User not found"}`
		gotBody := resp.Body.String()

		if wantBody != gotBody {
			t.Errorf("Want body: %s, Got body: %s", wantBody, gotBody)
		}
	})

}
