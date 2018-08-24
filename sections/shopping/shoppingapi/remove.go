package shoppingapi

import (
	"encoding/json"
	"net/http"

	"github.com/PaluMacil/dwn/dwn"
	"github.com/PaluMacil/dwn/sections/shopping"
)

// api/shopping/remove
func (rt *ShoppingRoute) handleRemove() {
	if spouse, err := rt.API().Current.Is(dwn.BuiltInGroupSpouse); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	} else if !spouse {
		http.Error(rt.W, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	if rt.R.Method != "DELETE" {
		rt.API().ServeMethodNotAllowed()
		return
	}
	if rt.R.Body == nil {
		rt.API().ServeBadRequest()
		return
	}
	var item shopping.Item
	err := json.NewDecoder(rt.R.Body).Decode(&item)
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	rt.Db.Shopping.Items.Delete(item.Name)
	rt.W.WriteHeader(http.StatusNoContent)
}
