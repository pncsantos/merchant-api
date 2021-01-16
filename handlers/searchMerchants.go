package handlers

import (
	"go-merchants/database"
	"go-merchants/models"
	"net/http"
	"strconv"
)

// SearchMerchants handler for getting list of merchants
func SearchMerchants(db database.MerchantsInterface, membersDb database.MembersInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
		}

		limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, err.Error())
		}

		res, err := db.Search(page, limit)
		if err != nil {
			if err.Error() == "No records found" {
				WriteResponse(w, http.StatusBadRequest, map[string]string{
					"message": "No records found",
				})
				return
			}
			WriteResponse(w, http.StatusBadGateway, err.Error())
			return
		}

		// get members count per merchant account
		merchants := []models.Merchant{}
		for _, merchant := range res.Merchants {
			merchant.MembersCount, _ = membersDb.GetMembersCountByMerchantID(merchant.ID)
			merchants = append(merchants, merchant)
		}

		WriteResponse(w, http.StatusOK, models.MerchantList{
			TotalCount: res.TotalCount,
			Merchants:  merchants,
		})
	}
}
