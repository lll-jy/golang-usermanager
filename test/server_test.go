package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_dummy(t *testing.T) {
	t.Run("hello world", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		PlayerServer(response, request)
		got := response.Body.String()
		fmt.Println(got)
	})
}

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
