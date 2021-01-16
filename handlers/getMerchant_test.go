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

func TestGetMerchant(t *testing.T) {
	id := AddNewMerchant()

	tests := map[string]struct {
		id           string
		expectedCode int
		expected     string
	}{
		"should return 200": {
			id:           id,
			expectedCode: 200,
			expected:     "Merchant X",
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
			req, _ := http.NewRequest("GET", "/merchants/"+test.id, nil)
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			r.HandleFunc("/merchants/{id}", handlers.GetMerchant(clientMerchant, clientMember))
			r.ServeHTTP(rec, req)

			if test.expectedCode == 200 {
				merchant := models.Merchant{}
				fmt.Println("rec.Body.String()", rec.Body.String())
				_ = json.Unmarshal([]byte(rec.Body.String()), &merchant)
				assert.Equal(t, test.expected, merchant.Name)
			}

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}

	// clean up
	_, _ = clientMerchant.Delete(id)
}
