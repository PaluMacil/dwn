package shoppingapi

import (
	"encoding/json"
	"net/http"

	"github.com/PaluMacil/dwn/dwn"
	"github.com/PaluMacil/dwn/sections/shopping"
)

// api/shopping/item
func (rt *ShoppingRoute) handleItem() {
	if rt.API().Current == nil {
		http.Error(rt.W, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
	if spouse, err := rt.API().Current.Is(dwn.BuiltInGroupSpouse); err != nil {
		rt.API().ServeInternalServerError(err)
		return
	} else if !spouse {
		http.Error(rt.W, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	switch rt.R.Method {
	case "POST":
		var item shopping.Item
		err := json.NewDecoder(rt.R.Body).Decode(&item)
		if err != nil {
			rt.API().ServeBadRequest()
			return
		}
		err = rt.Db.Shopping.Items.Set(item)
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
	case "DELETE":
		var item shopping.Item
		err := json.NewDecoder(rt.R.Body).Decode(&item)
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
		err = rt.Db.Shopping.Items.Delete(item.Name)
		if err != nil {
			rt.API().ServeInternalServerError(err)
			return
		}
		rt.W.WriteHeader(http.StatusNoContent)
	default:
		rt.API().ServeMethodNotAllowed()
		return
	}
}
