package handlers_test

import (
	"encoding/json"
	"go-merchants/handlers"
	"go-merchants/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertMember(t *testing.T) {
	id := AddNewMerchant()

	tests := map[string]struct {
		payload      string
		expectedCode int
		expected     string
	}{
		"should return 200": {
			payload:      `{"name":"member a","email":"testmerchant881@test.com","merchantId":"` + id + `"}`,
			expectedCode: 200,
			expected:     "member a",
		},
		"should return 400": {
			payload:      "invalid json string",
			expectedCode: 400,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			req, _ := http.NewRequest("POST", "/members", strings.NewReader(test.payload))
			rec := httptest.NewRecorder()

			h := http.HandlerFunc(handlers.InsertMember(clientMember))
			h.ServeHTTP(rec, req)

			if test.expectedCode == 200 {
				member := models.Member{}
				_ = json.Unmarshal([]byte(rec.Body.String()), &member)

				assert.Equal(t, test.expected, member.Name)
				assert.NotNil(t, member.ID)

				// cleanup
				_, _ = clientMember.Delete(member.ID.String())
			}

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}
}
