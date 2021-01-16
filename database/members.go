package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-merchants/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MembersInterface for members API calls
type MembersInterface interface {
	Insert(models.Member) (models.Member, error)
	Update(string, interface{}) (models.MemberUpdate, error)
	Delete(string) (models.MemberDelete, error)
	Get(string) (models.Member, error)
	Search(string, int64, int64) (models.MemberList, error)
	GetMembersCountByMerchantID(primitive.ObjectID) (int64, error)
}

// MemberClient used access pointers to collection
type MemberClient struct {
	Ctx context.Context
	Col *mongo.Collection
}

// Insert new member in the database and return its details if it was successfully added
func (c *MemberClient) Insert(docs models.Member) (models.Member, error) {
	member := models.Member{}

	// check if user exist
	condition := bson.M{"email": docs.Email}
	singleRes := c.Col.FindOne(c.Ctx, condition)

	if singleRes.Err() == nil {
		singleRes.Decode(&member)
		return member, errors.New("Email exist")
	}

	docs.CreatedDate = time.Now()

	res, err := c.Col.InsertOne(c.Ctx, docs)
	if err != nil {
		fmt.Println("err", err)
		return member, err
	}

	// convert ID string to object ID
	id := res.InsertedID.(primitive.ObjectID).Hex()

	return c.Get(id)
}

// Update members data by ID if there are necessary changes with the details
func (c *MemberClient) Update(id string, update interface{}) (models.MemberUpdate, error) {
	result := models.MemberUpdate{
		ModifiedCount: 0,
	}

	// convert id string to mongo object id
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	// get existing records
	member, err := c.Get(id)
	if err != nil {
		return result, err
	}

	// map member changes
	updateMap := update.(map[string]interface{})

	// check if member exist and if user changed it's email
	if updateMap["email"] != member.Email {
		condition := bson.M{"email": updateMap["email"]}
		err := c.Col.FindOne(c.Ctx, condition).Decode(&result.Result)
		if err == nil {
			return result, errors.New("Email exist")
		}
	}

	var exist map[string]interface{}
	//encode basic data types to JSON strings
	b, err := json.Marshal(member)
	if err != nil {
		return result, err
	}
	//decoding
	json.Unmarshal(b, &exist)

	// map existing update records
	for k := range updateMap {
		// update vs compare existing record
		if updateMap[k] == exist[k] {
			delete(updateMap, k)
		}
	}

	// return same value if there's no changes
	if len(updateMap) == 0 {
		return result, nil
	}

	// update database when are changes with the details
	res, err := c.Col.UpdateOne(c.Ctx, bson.M{"_id": _id}, bson.M{"$set": updateMap})
	if err != nil {
		return result, err
	}

	// get updated member details
	newMember, err := c.Get(id)
	if err != nil {
		return result, err
	}

	result.ModifiedCount = res.ModifiedCount
	result.Result = newMember

	return result, nil
}

// Delete members data by ID if they exist in the database
func (c *MemberClient) Delete(id string) (models.MemberDelete, error) {
	result := models.MemberDelete{
		DeletedCount: 0,
	}

	// check if merchant details
	_, err := c.Get(id)
	if err != nil {
		return result, errors.New("No records found")
	}

	// convert ID string to object ID
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	res, err := c.Col.DeleteOne(c.Ctx, bson.M{"_id": _id})
	if err != nil {
		return result, err
	}
	result.DeletedCount = res.DeletedCount
	return result, nil
}

// Get members details by ID in the database
func (c *MemberClient) Get(id string) (models.Member, error) {
	member := models.Member{}

	// convert ID string to object ID
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return member, err
	}

	// get member by id
	err = c.Col.FindOne(c.Ctx, bson.M{"_id": _id}).Decode(&member)
	if err != nil {
		return member, errors.New("No records found")
	}

	return member, nil
}

// Search list of members with paginated results in the database
func (c *MemberClient) Search(id string, page int64, limit int64) (models.MemberList, error) {
	memberList := models.MemberList{
		TotalCount: 0,
		Members:    nil,
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return memberList, err
	}

	// set option filters
	opts := options.Find().SetLimit(limit).SetSkip(page * limit)

	merchantID := bson.M{"merchantId": _id}
	cursor, err := c.Col.Find(c.Ctx, merchantID, opts)
	if err != nil {
		return memberList, err
	}

	defer cursor.Close(c.Ctx)

	members := []models.Member{}
	for cursor.Next(c.Ctx) {
		row := models.Member{}
		err := cursor.Decode(&row)
		if err != nil {
			fmt.Println("err", err)
		}
		members = append(members, row)
	}

	// get Total number of members in the collection
	count, err := c.GetMembersCountByMerchantID(_id)
	if err != nil {
		return memberList, err
	}

	memberList.Members = members
	memberList.TotalCount = count

	return memberList, nil
}

// GetMembersCountByMerchantID returns the total number of members within the selected merchant
func (c *MemberClient) GetMembersCountByMerchantID(id primitive.ObjectID) (int64, error) {
	merchantID := bson.M{"merchantId": id}
	count, err := c.Col.CountDocuments(c.Ctx, merchantID)
	if err != nil {
		return count, err
	}
	return count, nil
}
