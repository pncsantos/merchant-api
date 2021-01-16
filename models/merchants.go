package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Merchant represents the schema for the "Merchants" collection
type Merchant struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	MembersCount int64              `json:"membersCount" bson:"membersCount"`
	CreatedDate  time.Time          `json:"createdDate" bson:"createdDate,omitempty"`
	LogoFileName string             `json:"logoFileName" bson:"logoFileName,omitempty"`
	LogoBase64   string             `json:"logoBase64" bson:"logoBase64,omitempty"`
}

// MerchantUpdate represents the schema when merchant is updated
type MerchantUpdate struct {
	ModifiedCount int64 `json:"modifiedCount"`
	Result        Merchant
}

// MerchantDelete represents the schema when merchant is deleted
type MerchantDelete struct {
	DeletedCount int64 `json:"deletedCount"`
}

// MerchantList represents the schema for paginated list and overall count of all merchants
type MerchantList struct {
	TotalCount int64      `json:"totalCount"`
	Merchants  []Merchant `json:"merchants"`
}
