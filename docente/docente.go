package docente

import "DocenteEvaluacion/evaluacion"

// Docente representa a un profesor y mantiene una relación de
// "uno a muchos" con sus evaluaciones. El slice de evaluaciones es
// privado (minúscula) para forzar que toda modificación pase por los
// métodos públicos definidos abajo (AgregarEvaluacion, ObtenerEvaluaciones),
// en lugar de permitir que código externo manipule el slice directamente.
type Docente struct {
	ID     string
	Nombre string
	// evaluaciones es la relación (composición) hacia Evaluacion.
	// Es privada: el paquete docente controla cómo se agregan y
	// cómo se consultan las evaluaciones de un docente.
	evaluaciones []evaluacion.Evaluacion
}

func NuevoDocente(id, nombre string) *Docente {
	return &Docente{
		ID:           id,
		Nombre:       nombre,
		evaluaciones: []evaluacion.Evaluacion{},
	}
}

// validarEvaluacion es un método PRIVADO que aplica las reglas de
// negocio mínimas antes de aceptar una evaluación: que pertenezca al
// docente correcto y que tenga un ID no vacío. Al ser privado, esta
// regla de validación no puede ser invocada (ni por tanto saltada de
// forma directa) desde fuera del paquete docente; solo se ejecuta como
// parte del flujo controlado de AgregarEvaluacion.
func (d *Docente) validarEvaluacion(eval evaluacion.Evaluacion) bool {
	if eval.ID == "" {
		return false
	}
	if eval.DocenteID != d.ID {
		return false
	}
	return true
}

// AgregarEvaluacion es el método PÚBLICO para asociar una evaluación
// a este docente. Es la única forma permitida de modificar la relación
// Docente -> []Evaluacion desde fuera del paquete. Internamente delega
// la validación al método privado validarEvaluacion.
func (d *Docente) AgregarEvaluacion(eval evaluacion.Evaluacion) bool {
	if !d.validarEvaluacion(eval) {
		return false
	}
	d.evaluaciones = append(d.evaluaciones, eval)
	return true
}

// ObtenerEvaluaciones es un método PÚBLICO de solo lectura que
// devuelve únicamente las evaluaciones cuyo estado es "Finalizada".
// Devuelve un slice nuevo (no el interno) para que quien lo reciba no
// pueda mutar el estado real del docente desde afuera.
func (d Docente) ObtenerEvaluaciones() []evaluacion.Evaluacion {
	finalizadas := make([]evaluacion.Evaluacion, 0)
	for _, e := range d.evaluaciones {
		if e.EstaFinalizado() {
			finalizadas = append(finalizadas, e)
		}
	}
	return finalizadas
}

// TotalEvaluaciones devuelve la cantidad total de evaluaciones
// asociadas al docente, sin importar su estado. Útil para pruebas y
// para mostrar estadísticas sin exponer el slice interno.
func (d Docente) TotalEvaluaciones() int {
	return len(d.evaluaciones)
}
