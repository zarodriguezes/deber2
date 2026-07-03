package interfaces

type Evaluable interface {
	GetEstado() string
	EstaFinalizado() bool
	CalcularPromedio() float64
}
