package docente

import (
	"testing"

	"DocenteEvaluacion/evaluacion"
)

// TestValidarEvaluacionInterno prueba el método privado
// validarEvaluacion directamente, ya que este archivo sí pertenece
// al paquete "docente" (no usa sufijo _test en el nombre del
// paquete). Sirve como contraste con TestValidarEvaluacionPrivada en
// docente_externo_test.go, que demuestra que desde afuera ese mismo
// método NO es accesible.
func TestValidarEvaluacionInterno(t *testing.T) {
	doc := NuevoDocente("D1", "Carlos Ruiz")

	valida := evaluacion.NuevaEvaluacion("E1", "D1")
	if !doc.validarEvaluacion(*valida) {
		t.Errorf("se esperaba que la evaluación con DocenteID correcto fuera válida")
	}

	deOtroDocente := evaluacion.NuevaEvaluacion("E2", "D2")
	if doc.validarEvaluacion(*deOtroDocente) {
		t.Errorf("se esperaba que la evaluación de otro docente fuera inválida")
	}

	sinID := evaluacion.NuevaEvaluacion("", "D1")
	if doc.validarEvaluacion(*sinID) {
		t.Errorf("se esperaba que una evaluación sin ID fuera inválida")
	}
}

func TestNuevoDocente(t *testing.T) {
	doc := NuevoDocente("D1", "Carlos Ruiz")
	if doc.ID != "D1" || doc.Nombre != "Carlos Ruiz" {
		t.Errorf("docente creado incorrectamente: %+v", doc)
	}
	if doc.TotalEvaluaciones() != 0 {
		t.Errorf("un docente nuevo no debería tener evaluaciones")
	}
}
