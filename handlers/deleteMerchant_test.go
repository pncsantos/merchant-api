package handlers_test

import (
	"encoding/json"
	"go-merchants/handlers"
	"go-merchants/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDeleteMerchant(t *testing.T) {
	id := AddNewMerchant()

	tests := map[string]struct {
		id           string
		expectedCode int
		deletedCount int64
	}{
		"should return 200 and deleted count 1": {
			id:           id,
			expectedCode: 200,
			deletedCount: 1,
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
			req, _ := http.NewRequest("DELETE", "/merchants/"+test.id, nil)
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			r.HandleFunc("/merchants/{id}", handlers.DeleteMerchant(clientMerchant))
			r.ServeHTTP(rec, req)

			if test.expectedCode == 200 {
				merchant := models.MerchantDelete{}
				_ = json.Unmarshal([]byte(rec.Body.String()), &merchant)
				assert.Equal(t, test.deletedCount, merchant.DeletedCount)
			}

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}

	// clean up
	_, _ = clientMerchant.Delete(id)
}
