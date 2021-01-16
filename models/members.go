package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Member represents the schema for the "Members" collection
type Member struct {
	MerchantID  primitive.ObjectID `json:"merchantId" bson:"merchantId,omitempty"`
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name,omitempty"`
	Email       string             `json:"email" bson:"email,omitempty"`
	CreatedDate time.Time          `json:"createdDate" bson:"createdDate,omitempty"`
}

// MemberUpdate represents the schema when member is updated
type MemberUpdate struct {
	ModifiedCount int64 `json:"modifiedCount"`
	Result        Member
}

// MemberDelete represents the schema when member is deleted
type MemberDelete struct {
	DeletedCount int64 `json:"deletedCount"`
}

// MemberList represents the schema for paginated list and overall count of all members
type MemberList struct {
	TotalCount int64    `json:"totalCount"`
	Members    []Member `json:"members"`
}
