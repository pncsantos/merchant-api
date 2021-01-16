package handlers

import (
	"encoding/json"
	"go-merchants/database"
	"go-merchants/models"
	"io/ioutil"
	"net/http"
)

// InsertMember handler for adding new members in the database
func InsertMember(membersDb database.MembersInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		member := models.Member{}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &member)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := membersDb.Insert(member)
		if err != nil {
			if err.Error() == "Email exist" {
				WriteResponse(w, http.StatusBadRequest, map[string]string{
					"message": "Email is already in used",
				})
				return
			}
			WriteResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		WriteResponse(w, http.StatusOK, res)
	}
}
