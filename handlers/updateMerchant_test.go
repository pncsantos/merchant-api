package handlers_test

import (
	"encoding/json"
	"go-merchants/handlers"
	"go-merchants/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func AddNewMerchant() string {
	merchant := models.Merchant{
		Name: "Merchant X",
	}

	res, _ := clientMerchant.Insert(merchant)
	return res.ID.Hex()
}

func TestUpdateMerchant(t *testing.T) {
	id := AddNewMerchant()

	tests := map[string]struct {
		id            string
		payload       string
		expectedCode  int
		modifiedCount int64
	}{
		"should return 200 and modified count 1": {
			id:            id,
			payload:       `{"name": "Merchant Z"}`,
			expectedCode:  200,
			modifiedCount: 1,
		},
		"should return 400": {
			id:           id,
			payload:      "invalid json string",
			expectedCode: 400,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("PATCH", "/merchants/"+test.id, strings.NewReader(test.payload))
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			r.HandleFunc("/merchants/{id}", handlers.UpdateMerchant(clientMerchant))
			r.ServeHTTP(rec, req)

			if test.expectedCode == 200 {
				merchant := models.MerchantUpdate{}
				_ = json.Unmarshal([]byte(rec.Body.String()), &merchant)
				assert.Equal(t, test.modifiedCount, merchant.ModifiedCount)
			}

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}

	// clean up
	_, _ = clientMerchant.Delete(id)
}
