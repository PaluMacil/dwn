package dashboard

const DashboardPrefix = "DASHBOARD:BOARD"

type Dashboard struct {
	Projects []Project `json:"projects"`
}

func (d Dashboard) Key() []byte {
	return d.Prefix()
}

func (d Dashboard) Prefix() []byte {
	return []byte(DashboardPrefix)
}
