package tc

type KeyParams struct {
	Src  string
	Dest string
}

type Proxy interface {
	SetRule(params KeyParams, rule Rule)
}
