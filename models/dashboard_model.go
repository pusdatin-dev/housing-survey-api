package models

type DashboardResource struct {
	Name  string `json:"name"`
	Total int64  `json:"total"`
}

type DashboardProgramType struct {
	Name    string  `json:"name"`
	Total   int     `json:"total"`
	Percent float64 `json:"percent"`
}
