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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddNewMember() string {
	// merchant ID
	merchantID := AddNewMerchant()
	_id, _ := primitive.ObjectIDFromHex(merchantID)

	member := models.Member{
		Name:       "Member X",
		Email:      "test8888@test.com",
		MerchantID: _id,
	}

	res, _ := clientMember.Insert(member)
	return res.ID.Hex()
}

func TestUpdateMembers(t *testing.T) {
	id := AddNewMember()

	tests := map[string]struct {
		id            string
		payload       string
		expectedCode  int
		modifiedCount int64
	}{
		"should return 200 and modified count 1": {
			id:            id,
			payload:       `{"name": "Member Z"}`,
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
			req, _ := http.NewRequest("PATCH", "/members/"+test.id, strings.NewReader(test.payload))
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			r.HandleFunc("/members/{id}", handlers.UpdateMember(clientMember))
			r.ServeHTTP(rec, req)

			if test.expectedCode == 200 {
				member := models.MemberUpdate{}
				_ = json.Unmarshal([]byte(rec.Body.String()), &member)
				assert.Equal(t, test.modifiedCount, member.ModifiedCount)
			}

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}

	// clean up
	_, _ = clientMember.Delete(id)
}
