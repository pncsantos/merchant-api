package handlers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"go-merchants/config"
	"go-merchants/database"
	"go-merchants/handlers"
	"go-merchants/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var clientMerchant database.MerchantsInterface
var clientMember database.MembersInterface

func init() {
	conf := config.MongoConfiguration{
		Server:                  "mongodb://localhost:27017",
		Database:                "Mgo",
		MerchantsTestCollection: "MerchantsTest",
		MembersTestCollection:   "MembersTest",
	}
	ctx := context.TODO()

	db := database.ConnectDB(ctx, conf)

	merchantCollections := db.Collection(conf.MerchantsTestCollection)
	clientMerchant = &database.MerchantClient{
		Col: merchantCollections,
		Ctx: ctx,
	}

	membersCollections := db.Collection(conf.MembersTestCollection)
	clientMember = &database.MemberClient{
		Col: membersCollections,
		Ctx: ctx,
	}
}

func TestInsertMerchant(t *testing.T) {
	tests := map[string]struct {
		payload      string
		expectedCode int
		expected     string
	}{
		"should return 200": {
			payload:      `{"name":"merchant a"}`,
			expectedCode: 200,
			expected:     "merchant a",
		},
		"should return 400": {
			payload:      "invalid json string",
			expectedCode: 400,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			req, _ := http.NewRequest("POST", "/merchants", strings.NewReader(test.payload))
			rec := httptest.NewRecorder()

			h := http.HandlerFunc(handlers.InsertMerchant(clientMerchant))
			h.ServeHTTP(rec, req)
			fmt.Println(rec.Body.String())
			if test.expectedCode == 200 {
				merchant := models.Merchant{}
				_ = json.Unmarshal([]byte(rec.Body.String()), &merchant)
				assert.Equal(t, test.expected, merchant.Name)
				assert.NotNil(t, merchant.ID)

				// cleanup
				_, _ = clientMerchant.Delete(merchant.ID.String())
			}

			assert.Equal(t, test.expectedCode, rec.Code)
		})
	}
}
