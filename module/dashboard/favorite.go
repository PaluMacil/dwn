package dashboard

import (
	"github.com/PaluMacil/dwn/database/store"
)

type Favorite struct {
	ProjectID store.Identity `json:"project_id"`
}
