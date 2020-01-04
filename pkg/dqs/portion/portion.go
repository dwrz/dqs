package portion

type Amount string

const (
	None = ""
	Full = "full"
	Half = "half"
)

type Portion struct {
	Points  int   `json:"points"`
	Amount Amount `json:"amount"`
}

func (p *Portion) CalculateScore() float64 {
	switch p.Amount {
	case Full:
		return float64(p.Points)
	case Half:
		return float64(p.Points) / 2
	case None:
		fallthrough
	default:
		return float64(0)
	}
}
