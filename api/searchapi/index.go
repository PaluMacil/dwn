package searchapi

import (
	"github.com/PaluMacil/dwn/core"
)

// api/search/index/{name|list}
func (rt *SearchRoute) handleIndex() {
	if rt.API().ServeCannot(core.PermissionManageIndexes) {
		return
	}
	switch rt.R.Method {
	case "GET":
	case "POST":
	default:
		rt.API().ServeMethodNotAllowed()
	}
}
