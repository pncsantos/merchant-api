package handlers

import (
	"encoding/json"
	"go-merchants/database"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// UpdateMember handler for updating member details by ID in the database
func UpdateMember(db database.MembersInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var merchant interface{}
		err = json.Unmarshal(body, &merchant)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := db.Update(id, merchant)
		if err != nil {
			if err.Error() == "Email exist" {
				WriteResponse(w, http.StatusBadRequest, map[string]string{
					"message": "Email is already in used",
				})
				return
			} else if err.Error() == "No records found" {
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
