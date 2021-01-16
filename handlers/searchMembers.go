package handlers

import (
	"go-merchants/database"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// SearchMembers handler for getting list of members by merchant's ID
func SearchMembers(db database.MembersInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
		}

		limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
		}

		res, err := db.Search(id, page, limit)
		if err != nil {
			WriteResponse(w, http.StatusBadGateway, err.Error())
			return
		}

		WriteResponse(w, http.StatusOK, res)
	}
}
