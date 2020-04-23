package core

type Me struct {
	User    UserInfo `json:"user"`
	Session Session  `json:"session"`
	Groups  []Group  `json:"groups"`
}
