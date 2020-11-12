package entities

type KLine struct {
	Labels []string    `json:"labels"`
	KLine  [][]float64 `json:"k_line"`
}
