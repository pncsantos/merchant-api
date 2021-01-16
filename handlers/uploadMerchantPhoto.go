package handlers

import (
	"encoding/json"
	"go-merchants/database"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// UploadMerchantPhoto handler for uploading new photo of merchant in the database
func UploadMerchantPhoto(db database.MerchantsInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// return error if file is more than 50kb
		if len(body) > 1000*50 {
			WriteResponse(w, http.StatusBadRequest, map[string]string{
				"message": "Image is too big!",
			})
			return
		}

		// get merchant details and check if user exist
		merchant, err := db.Get(id)
		if err != nil {
			if err.Error() == "No records found" {
				WriteResponse(w, http.StatusBadRequest, map[string]string{
					"message": "Merchant doesn't exist",
				})
				return
			}

			WriteResponse(w, http.StatusBadGateway, err.Error())
			return
		}

		// create unique filename
		fileName := "merchant_" + id + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)

		// upload file
		var uploadErr = db.UploadFile(fileName, body)
		if uploadErr != nil {
			WriteResponse(w, http.StatusBadGateway, uploadErr.Error())
			return
		}

		base64, err := db.GenerateBase64Image(fileName)
		if err != nil {
			WriteResponse(w, http.StatusBadGateway, err.Error())
			return
		}

		// add merchant logo details
		merchant.LogoFileName = fileName
		merchant.LogoBase64 = base64

		var merchantMap map[string]interface{}
		// encode basic data types to JSON strings
		b, err := json.Marshal(merchant)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		//decoding
		json.Unmarshal(b, &merchantMap)

		// update merchant details with photo details
		res, err := db.Update(id, merchantMap)
		if err != nil {
			if err.Error() == "No records found" {
				WriteResponse(w, http.StatusBadRequest, map[string]string{
					"message": "Merchant doesn't exist",
				})
				return
			}
			WriteResponse(w, http.StatusBadGateway, err.Error())
			return
		}

		WriteResponse(w, http.StatusOK, res)
	}
}
