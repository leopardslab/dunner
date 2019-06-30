package initialize

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListTemplatesSuccess(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "recipes:\n  - foo\n  - bar", http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	recipeListURL = server.URL

	err := ListRecipes()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestListTemplatesSuccessNoRecipes(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusOK)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	recipeListURL = server.URL

	err := ListRecipes()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestListTemplatesWhen404(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusNotFound)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	recipeListURL = server.URL

	err := ListRecipes()

	expectedErr := fmt.Sprintf("Failed to download list of recipes from dunner cookbook: Error downloading file %s: 404 Not Found", server.URL)
	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected %s, got %s", expectedErr, err)
	}
}

