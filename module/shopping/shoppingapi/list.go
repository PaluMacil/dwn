package shoppingapi

import (
	"encoding/json"
	"net/http"

	"github.com/PaluMacil/dwn/core"
)

// api/shopping/list
func (rt *ShoppingRoute) handleList() {
	if spouse, err := rt.API().Current.Is(core.BuiltInGroupSpouse); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	} else if !spouse {
		http.Error(rt.W, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if rt.R.Method != "GET" {
		rt.API().ServeMethodNotAllowed()
		return
	}
	if rt.R.Body == nil {
		rt.API().ServeBadRequest()
		return
	}
	items, err := rt.Db.Shopping.Items.All()
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	if err := json.NewEncoder(rt.W).Encode(items); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
}
