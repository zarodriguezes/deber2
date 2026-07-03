package evaluacion

type EstadoEvaluacion string

const (
	Borrador   EstadoEvaluacion = "borrador"
	EnProceso  EstadoEvaluacion = "en_proceso"
	Finalizada EstadoEvaluacion = "finalizada"
)

// Evaluacion representa la evaluación de un docente.
// Los campos en minúscula (comentarios, puntajes) son privados:
// solo el propio paquete evaluacion puede modificarlos directamente.
// Desde fuera del paquete solo se puede interactuar a través de los
// métodos públicos (AgregarComentarioPublico, CalcularPromedio, etc.)
type Evaluacion struct {
	ID        string
	DocenteID string
	Estado    EstadoEvaluacion
	Puntaje   float64
	// Campo privado
	comentarios []string
	// puntajes almacena cada calificación individual registrada para
	// esta evaluación, y es la base real para calcular el promedio.
	puntajes []float64
}

func NuevaEvaluacion(id, docenteID string) *Evaluacion {
	return &Evaluacion{
		ID:          id,
		DocenteID:   docenteID,
		Estado:      Borrador,
		Puntaje:     0,
		comentarios: []string{},
		puntajes:    []float64{},
	}
}

func (e Evaluacion) GetEstado() string {
	return string(e.Estado)
}

func (e Evaluacion) EstaFinalizado() bool {
	return e.Estado == Finalizada
}

// RegistrarPuntaje agrega una calificación individual a la evaluación.
// Se valida que el puntaje esté en un rango razonable (0-100) antes de
// aceptarlo, evitando que datos inválidos contaminen el promedio.
func (e *Evaluacion) RegistrarPuntaje(puntaje float64) bool {
	if puntaje < 0 || puntaje > 100 {
		return false
	}
	e.puntajes = append(e.puntajes, puntaje)
	e.Puntaje = e.calcularPromedioInterno()
	return true
}

// calcularPromedioInterno es un método PRIVADO (empieza en minúscula).
// Encapsula la lógica real del cálculo del promedio aritmético de los
// puntajes registrados. Al ser privado, solo puede ser invocado desde
// código dentro del paquete evaluacion, lo que evita que código externo
// dependa de los detalles internos de cómo se calcula el promedio.
func (e Evaluacion) calcularPromedioInterno() float64 {
	if len(e.puntajes) == 0 {
		return 0
	}
	suma := 0.0
	for _, p := range e.puntajes {
		suma += p
	}
	return suma / float64(len(e.puntajes))
}

// agregarComentario es un método PRIVADO que realiza la inserción
// "cruda" del comentario en el slice interno, sin validaciones.
// Se mantiene privado porque asume que el comentario ya fue validado
// por el método público que lo invoca (AgregarComentarioPublico).
func (e *Evaluacion) agregarComentario(comentario string) {
	e.comentarios = append(e.comentarios, comentario)
}

// CalcularPromedio es el método PÚBLICO requerido por la interfaz
// Evaluable. Delega el cálculo real al método privado
// calcularPromedioInterno, exponiendo así una API simple hacia afuera
// mientras la lógica interna queda encapsulada.
func (e Evaluacion) CalcularPromedio() float64 {
	return e.calcularPromedioInterno()
}

// AgregarComentarioPublico es el método PÚBLICO para agregar
// comentarios. Valida la entrada (que no esté vacío) antes de
// delegar al método privado agregarComentario. Esta es la única
// puerta de entrada permitida para modificar los comentarios desde
// fuera del paquete.
func (e *Evaluacion) AgregarComentarioPublico(comentario string) bool {
	if comentario == "" {
		return false
	}
	e.agregarComentario(comentario)
	return true
}

// ObtenerComentarios devuelve una copia de los comentarios registrados.
// Se devuelve una copia (no el slice interno directamente) para que
// quien llame no pueda mutar el estado interno de la evaluación.
func (e Evaluacion) ObtenerComentarios() []string {
	copia := make([]string, len(e.comentarios))
	copy(copia, e.comentarios)
	return copia
}

// Finalizar cambia el estado de la evaluación a Finalizada.
func (e *Evaluacion) Finalizar() {
	e.Estado = Finalizada
}
