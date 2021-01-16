package database

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go-merchants/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MerchantsInterface for api calls
type MerchantsInterface interface {
	Insert(models.Merchant) (models.Merchant, error)
	Update(string, interface{}) (models.MerchantUpdate, error)
	Delete(string) (models.MerchantDelete, error)
	Get(string) (models.Merchant, error)
	Search(int64, int64) (models.MerchantList, error)
	UploadFile(string, []byte) error
	GenerateBase64Image(string) (string, error)
}

// MerchantClient used access pointers to collection
type MerchantClient struct {
	Ctx context.Context
	Col *mongo.Collection
}

// Insert new merchant in the database and return its details if it was successfully added
func (c *MerchantClient) Insert(docs models.Merchant) (models.Merchant, error) {
	merchant := models.Merchant{}

	docs.CreatedDate = time.Now()
	res, err := c.Col.InsertOne(c.Ctx, docs)
	if err != nil {
		return merchant, err
	}

	// convert ID string to object ID
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return c.Get(id)
}

// Update merchants data by ID if there are necessary changes
func (c *MerchantClient) Update(id string, update interface{}) (models.MerchantUpdate, error) {
	result := models.MerchantUpdate{
		ModifiedCount: 0,
	}

	// convert ID string to object ID
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	// get existing records
	merchant, err := c.Get(id)
	if err != nil {
		return result, errors.New("No records found")
	}

	var exist map[string]interface{}
	//encode basic data types to JSON strings
	b, err := json.Marshal(merchant)
	if err != nil {
		return result, err
	}
	//decoding
	json.Unmarshal(b, &exist)

	// map existing update records
	change := update.(map[string]interface{})
	for k := range change {
		// compare updates vs existing record
		if change[k] == exist[k] {
			delete(change, k)
		}
	}

	// return same value if there's no changes
	if len(change) == 0 {
		return result, nil
	}

	// update database when are changes with the details
	res, err := c.Col.UpdateOne(c.Ctx, bson.M{"_id": _id}, bson.M{"$set": change})
	if err != nil {
		return result, err
	}

	// get updated merchant details
	newMerchant, err := c.Get(id)
	if err != nil {
		return result, errors.New("No records found")
	}

	result.ModifiedCount = res.ModifiedCount
	result.Result = newMerchant

	return result, nil
}

// Delete merchants data by ID if they exist in the database
func (c *MerchantClient) Delete(id string) (models.MerchantDelete, error) {
	result := models.MerchantDelete{
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

// Get merchants details by ID in the database
func (c *MerchantClient) Get(id string) (models.Merchant, error) {
	merchant := models.Merchant{}

	// convert ID string to object ID
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return merchant, err
	}

	// get merchant by id
	err = c.Col.FindOne(c.Ctx, bson.M{"_id": _id}).Decode(&merchant)
	if err != nil {
		return merchant, errors.New("No records found")
	}

	return merchant, nil
}

// Search list of merchants with paginated results in the database
func (c *MerchantClient) Search(page int64, limit int64) (models.MerchantList, error) {
	merchantList := models.MerchantList{
		TotalCount: 0,
		Merchants:  nil,
	}

	// set option filters
	opts := options.Find().SetLimit(limit).SetSkip(page * limit)

	cursor, err := c.Col.Find(c.Ctx, bson.M{}, opts)
	if err != nil {
		return merchantList, errors.New("No records found")
	}

	merchants := []models.Merchant{}
	for cursor.Next(c.Ctx) {
		row := models.Merchant{}
		cursor.Decode(&row)
		merchants = append(merchants, row)
	}

	// get Total number of merchants in the collection
	count, err := c.Col.CountDocuments(c.Ctx, bson.M{})
	if err != nil {
		return merchantList, err
	}

	merchantList.Merchants = merchants
	merchantList.TotalCount = count

	return merchantList, nil
}

// UploadFile saves the file in the database
func (c *MerchantClient) UploadFile(fileName string, file []byte) error {
	bucket, err := gridfs.NewBucket(c.Col.Database())
	if err != nil {
		return err
	}

	uploadStream, err := bucket.OpenUploadStream(fileName)
	if err != nil {
		return err
	}

	defer uploadStream.Close()

	fileSize, err := uploadStream.Write(file)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully Uploaded File, with size of: %d bytes\n", fileSize)

	return nil
}

// GenerateBase64Image creates file in the database and returns Base64 once its successful
func (c *MerchantClient) GenerateBase64Image(fileName string) (string, error) {
	bucket, err := gridfs.NewBucket(c.Col.Database())
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(fileName, &buf)
	if err != nil {
		return "", err
	}

	// ioutil.WriteFile(fileName, buf.Bytes(), 0600)

	fmt.Printf("File size to download: %v\n", dStream)

	// encode bytes to base64
	sEnc := b64.StdEncoding.EncodeToString(buf.Bytes())

	return sEnc, nil
}
