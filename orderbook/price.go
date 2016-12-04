package orderbook

type Price float64

func (p Price) Value() float64 {
	return float64(p)
}

type Volume float64

func (v Volume) Value() float64 {
	return float64(v)
}

type Prices []Price

func (p Prices) Len() int {
	return len(p)
}

func (p Prices) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Prices) Less(i, j int) bool {
	return p[i] < p[j]
}
