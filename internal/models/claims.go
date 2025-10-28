package models

type Claims struct {
	Username string   `json:"username"`
	Groups   []string `json:"groups"`
}
