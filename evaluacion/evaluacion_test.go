package evaluacion

import "testing"

// TestCalcularPromedioInterno prueba directamente el método PRIVADO
// calcularPromedioInterno. Esto solo es posible porque este archivo de
// test pertenece al mismo paquete "evaluacion" (no usa _test como
// sufijo de paquete, sino "package evaluacion"). Es una prueba interna
// (white-box testing): conoce y verifica los detalles de implementación.
func TestCalcularPromedioInterno(t *testing.T) {
	eval := NuevaEvaluacion("E1", "D1")

	// Caso 1: sin puntajes, el promedio debe ser 0
	if got := eval.calcularPromedioInterno(); got != 0 {
		t.Errorf("esperaba promedio 0 sin puntajes, obtuve %v", got)
	}

	// Caso 2: con varios puntajes registrados
	eval.RegistrarPuntaje(80)
	eval.RegistrarPuntaje(90)
	eval.RegistrarPuntaje(100)

	esperado := (80.0 + 90.0 + 100.0) / 3.0
	if got := eval.calcularPromedioInterno(); got != esperado {
		t.Errorf("esperaba promedio %v, obtuve %v", esperado, got)
	}

	// Caso 3: un puntaje fuera de rango no debería haberse registrado
	ok := eval.RegistrarPuntaje(150)
	if ok {
		t.Errorf("se esperaba que RegistrarPuntaje rechazara un valor fuera de rango (150)")
	}
}

// TestAgregarComentario prueba directamente el método PRIVADO
// agregarComentario, verificando que inserta el comentario en el
// slice interno sin aplicar validaciones (esas validaciones son
// responsabilidad del método público AgregarComentarioPublico).
func TestAgregarComentario(t *testing.T) {
	eval := NuevaEvaluacion("E2", "D1")

	eval.agregarComentario("Buen desempeño en clase")
	eval.agregarComentario("Necesita mejorar puntualidad")

	comentarios := eval.ObtenerComentarios()
	if len(comentarios) != 2 {
		t.Fatalf("esperaba 2 comentarios, obtuve %d", len(comentarios))
	}
	if comentarios[0] != "Buen desempeño en clase" {
		t.Errorf("comentario inesperado en posición 0: %v", comentarios[0])
	}
	if comentarios[1] != "Necesita mejorar puntualidad" {
		t.Errorf("comentario inesperado en posición 1: %v", comentarios[1])
	}

	// agregarComentario es "tonto" a propósito: no valida cadenas
	// vacías, eso lo hace el método público. Lo comprobamos aquí
	// para dejar documentado ese contrato interno.
	eval.agregarComentario("")
	comentarios = eval.ObtenerComentarios()
	if len(comentarios) != 3 {
		t.Errorf("agregarComentario privado no debería validar vacíos; esperaba 3 comentarios, obtuve %d", len(comentarios))
	}
}

// TestAgregarComentarioPublico prueba el método PÚBLICO
// AgregarComentarioPublico, verificando que sí aplica la validación
// (rechaza comentarios vacíos) antes de delegar al privado
// agregarComentario.
func TestAgregarComentarioPublico(t *testing.T) {
	eval := NuevaEvaluacion("E3", "D1")

	if ok := eval.AgregarComentarioPublico("Excelente manejo del aula"); !ok {
		t.Errorf("se esperaba poder agregar un comentario válido")
	}
	if ok := eval.AgregarComentarioPublico(""); ok {
		t.Errorf("se esperaba que un comentario vacío fuera rechazado")
	}

	comentarios := eval.ObtenerComentarios()
	if len(comentarios) != 1 {
		t.Fatalf("esperaba 1 comentario válido registrado, obtuve %d", len(comentarios))
	}
	if comentarios[0] != "Excelente manejo del aula" {
		t.Errorf("comentario inesperado: %v", comentarios[0])
	}
}
