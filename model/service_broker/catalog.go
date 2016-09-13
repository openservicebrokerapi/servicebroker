package model

type GetCatalogResponse struct {
	Services []Service `json:"services"`
}
