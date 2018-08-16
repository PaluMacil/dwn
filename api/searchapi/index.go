package searchapi

import (
	"github.com/PaluMacil/dwn/dwn"
)

// api/search/index/{name|list}
func (rt *SearchRoute) handleIndex() {
	if rt.API().ServeCannot(dwn.PermissionManageIndexes) {
		return
	}
	switch rt.R.Method {
	case "GET":
	case "POST":
	default:
		rt.API().ServeMethodNotAllowed()
	}
}
