package dtos

type IndicationsDTO struct {
	Id string `json:"id"`
	Indications map[string]float64 `json:"indications"`
}