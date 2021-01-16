package handlers_test

import (
	"encoding/json"
	"fmt"
	"go-merchants/handlers"
	"go-merchants/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetMember(t *testing.T) {
	id := AddNewMember()

	tests := map[string]struct {
		id           string
		expectedCode int
		expected     string
	}{
		"should return 200": {
			id:           id,
			expectedCode: 200,
			expected:     "Member X",
		},
		"should return 400": {
			id:           "abc",
			expectedCode: 400,
		},
		"should return 404": {
			id:           "",
			expectedCode: 404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/members/"+test.id, nil)
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			r.HandleFunc("/members/{id}", handlers.GetMember(clientMember))
			r.ServeHTTP(rec, req)

			if test.expectedCode == 200 {
				member := models.Member{}
				fmt.Println("rec.Body.String()", rec.Body.String())
				_ = json.Unmarshal([]byte(rec.Body.String()), &member)
				assert.Equal(t, test.expected, member.Name)
			}

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}

	// clean up
	_, _ = clientMember.Delete(id)
}
