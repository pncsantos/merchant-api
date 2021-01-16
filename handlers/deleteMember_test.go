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

func TestDeleteMember(t *testing.T) {
	id := AddNewMember()

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
			expectedCode: 502,
		},
		"should return 404": {
			id:           "",
			expectedCode: 404,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/members/"+test.id, nil)
			rec := httptest.NewRecorder()

			r := mux.NewRouter()
			r.HandleFunc("/members/{id}", handlers.DeleteMember(clientMember))
			r.ServeHTTP(rec, req)

			if test.expectedCode == 200 {
				member := models.MemberDelete{}
				_ = json.Unmarshal([]byte(rec.Body.String()), &member)
				assert.Equal(t, test.deletedCount, member.DeletedCount)
			}

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}

	// clean up
	_, _ = clientMember.Delete(id)
}
