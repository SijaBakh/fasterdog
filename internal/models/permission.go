package models

type PermissionsResult struct {
	Routes []RoutePermission `json:"routes"`
	Groups []string          `json:"groups"`
	CI     []string          `json:"ci"`
}

type RoutePermission struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}
