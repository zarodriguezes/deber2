// Este archivo usa "package docente_test" (con sufijo _test), lo que
// lo convierte en un paquete EXTERNO respecto a "docente" y
// "evaluacion". Esto simula exactamente cómo vería el código un
// usuario cualquiera de estos paquetes: solo tiene acceso a lo que
// es público (mayúscula inicial).
package docente_test

import (
	"testing"

	"DocenteEvaluacion/docente"
	"DocenteEvaluacion/evaluacion"
	"DocenteEvaluacion/interfaces"
)

// TestAgregarYListarEvaluaciones prueba el flujo completo
// Docente <-> Evaluacion usando ÚNICAMENTE métodos públicos, tal como
// lo haría cualquier consumidor externo de estos paquetes.
func TestAgregarYListarEvaluaciones(t *testing.T) {
	doc := docente.NuevoDocente("D1", "Ana Pérez")

	eval1 := evaluacion.NuevaEvaluacion("E1", "D1")
	eval1.RegistrarPuntaje(90)
	eval1.Finalizar()

	eval2 := evaluacion.NuevaEvaluacion("E2", "D1")
	eval2.RegistrarPuntaje(70)
	// eval2 se deja en estado Borrador, no finalizada

	if ok := doc.AgregarEvaluacion(*eval1); !ok {
		t.Fatalf("se esperaba poder agregar eval1 al docente")
	}
	if ok := doc.AgregarEvaluacion(*eval2); !ok {
		t.Fatalf("se esperaba poder agregar eval2 al docente")
	}

	if total := doc.TotalEvaluaciones(); total != 2 {
		t.Errorf("esperaba 2 evaluaciones totales, obtuve %d", total)
	}

	finalizadas := doc.ObtenerEvaluaciones()
	if len(finalizadas) != 1 {
		t.Fatalf("esperaba 1 evaluación finalizada, obtuve %d", len(finalizadas))
	}
	if finalizadas[0].ID != "E1" {
		t.Errorf("esperaba que la evaluación finalizada fuera E1, obtuve %v", finalizadas[0].ID)
	}

	// También se prueba el rechazo de una evaluación que no
	// corresponde a este docente (regla validada internamente por
	// el método privado validarEvaluacion, accedido solo a través
	// de AgregarEvaluacion).
	evalOtroDocente := evaluacion.NuevaEvaluacion("E3", "D2")
	if ok := doc.AgregarEvaluacion(*evalOtroDocente); ok {
		t.Errorf("no se debería poder agregar una evaluación de otro docente")
	}
}

// TestValidarEvaluacionPrivada documenta y demuestra que el método
// validarEvaluacion NO es accesible desde un paquete externo.
//
// En Go, intentar invocar un identificador no exportado (minúscula)
// de otro paquete es un ERROR DE COMPILACIÓN, no un error en tiempo
// de ejecución. Por eso este test no puede "ejecutar" el intento de
// acceso: si se descomenta la línea de abajo, el paquete completo
// deja de compilar y `go test` falla con un mensaje como:
//
//	./docente_externo_test.go:XX:6: doc.validarEvaluacion undefined
//	(cannot refer to unexported field or method docente.validarEvaluacion)
//
// La prueba real de la encapsulación es justamente esa: el compilador
// impide el acceso. Dejamos la línea comentada como evidencia, y el
// test pasa demostrando que el código compila y funciona SIN
// necesitar acceso al método privado.
func TestValidarEvaluacionPrivada(t *testing.T) {
	doc := docente.NuevoDocente("D1", "Ana Pérez")
	eval := evaluacion.NuevaEvaluacion("E1", "D1")

	// La siguiente línea, si se descomenta, NO COMPILA porque
	// validarEvaluacion es un método privado del paquete docente:
	//
	// doc.validarEvaluacion(*eval)
	//
	// error esperado del compilador:
	// doc.validarEvaluacion undefined (cannot refer to unexported
	// field or method docente.validarEvaluacion)

	// En su lugar, verificamos que la única forma pública y
	// correcta de acceder a esa validación es AgregarEvaluacion.
	if ok := doc.AgregarEvaluacion(*eval); !ok {
		t.Fatalf("se esperaba que AgregarEvaluacion validara y aceptara la evaluación")
	}

	_ = doc
}

// TestInterfaceEvaluable prueba que *evaluacion.Evaluacion implemente
// la interfaz interfaces.Evaluable. Se hace de dos maneras: una
// aserción de tipo en tiempo de compilación y, además, usando el
// valor a través de una variable de tipo interfaz para confirmar que
// los métodos responden correctamente en tiempo de ejecución.
func TestInterfaceEvaluable(t *testing.T) {
	eval := evaluacion.NuevaEvaluacion("E1", "D1")
	eval.RegistrarPuntaje(85)
	eval.RegistrarPuntaje(95)
	eval.Finalizar()

	// Aserción en tiempo de compilación: si *Evaluacion dejara de
	// implementar Evaluable, esta línea por sí sola haría fallar la
	// compilación del test.
	var _ interfaces.Evaluable = eval

	// Verificación en tiempo de ejecución usando la variable como
	// interfaz, no como el tipo concreto.
	var ev interfaces.Evaluable = eval

	if ev.GetEstado() != "finalizada" {
		t.Errorf("esperaba estado 'finalizada', obtuve %v", ev.GetEstado())
	}
	if !ev.EstaFinalizado() {
		t.Errorf("esperaba que EstaFinalizado() devolviera true")
	}

	promedioEsperado := (85.0 + 95.0) / 2.0
	if got := ev.CalcularPromedio(); got != promedioEsperado {
		t.Errorf("esperaba promedio %v, obtuve %v", promedioEsperado, got)
	}
}
