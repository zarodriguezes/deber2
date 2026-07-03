# DocenteEvaluacion

Proyecto en Go que implementa encapsulación, relaciones entre estructuras
e interfaces para manejar evaluaciones de docentes.

## Estructura del proyecto

```
DocenteEvaluacion/
├── docente/
│   ├── docente.go
│   ├── docente_test.go            (pruebas internas, package docente)
│   └── docente_externo_test.go    (pruebas externas, package docente_test)
├── evaluacion/
│   ├── evaluacion.go
│   └── evaluacion_test.go         (pruebas internas, package evaluacion)
├── interfaces/
│   └── evaluable.go
├── go.mod
└── README.md
```

## Comandos ejecutados y resultados

### 1. Compilación

```
$ go build ./...
```
Resultado: sin errores ni salida (compila correctamente).

### 2. Verificación estática y formato

```
$ go vet ./...
$ gofmt -l .
```
Resultado: sin advertencias de `vet` y sin archivos pendientes de formatear.

### 3. Batería de pruebas

```
$ go test ./... -v
```

Salida real obtenida:

```
=== RUN   TestValidarEvaluacionInterno
--- PASS: TestValidarEvaluacionInterno (0.00s)
=== RUN   TestNuevoDocente
--- PASS: TestNuevoDocente (0.00s)
=== RUN   TestAgregarYListarEvaluaciones
--- PASS: TestAgregarYListarEvaluaciones (0.00s)
=== RUN   TestValidarEvaluacionPrivada
--- PASS: TestValidarEvaluacionPrivada (0.00s)
=== RUN   TestInterfaceEvaluable
--- PASS: TestInterfaceEvaluable (0.00s)
PASS
ok      DocenteEvaluacion/docente      0.002s
=== RUN   TestCalcularPromedioInterno
--- PASS: TestCalcularPromedioInterno (0.00s)
=== RUN   TestAgregarComentario
--- PASS: TestAgregarComentario (0.00s)
=== RUN   TestAgregarComentarioPublico
--- PASS: TestAgregarComentarioPublico (0.00s)
PASS
ok      DocenteEvaluacion/evaluacion   0.002s
?       DocenteEvaluacion/interfaces   [no test files]
```

### 4. Cobertura de pruebas

Como las pruebas externas de `docente_externo_test.go` (paquete
`docente_test`) también ejercitan código del paquete `evaluacion`, se
usó `-coverpkg=./...` para medir la cobertura real combinada:

```
$ go test -coverpkg=./... -coverprofile=cover.out ./...
$ go tool cover -func=cover.out
```

Resultado:

```
DocenteEvaluacion/docente/docente.go:19:    NuevoDocente             100.0%
DocenteEvaluacion/docente/docente.go:33:    validarEvaluacion        100.0%
DocenteEvaluacion/docente/docente.go:47:    AgregarEvaluacion        100.0%
DocenteEvaluacion/docente/docente.go:59:    ObtenerEvaluaciones      100.0%
DocenteEvaluacion/docente/docente.go:72:    TotalEvaluaciones        100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:28:  NuevaEvaluacion          100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:39:  GetEstado                100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:43:  EstaFinalizado           100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:50:  RegistrarPuntaje         100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:64:  calcularPromedioInterno  100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:79:  agregarComentario        100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:87:  CalcularPromedio         100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:96:  AgregarComentarioPublico 100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:107: ObtenerComentarios       100.0%
DocenteEvaluacion/evaluacion/evaluacion.go:114: Finalizar                100.0%
total:                                          (statements)             100.0%
```

**Cobertura total: 100.0% de las sentencias.**

> Nota: si se ejecuta `go test ./... -cover` (sin `-coverpkg`), el
> paquete `evaluacion` mostrará un porcentaje más bajo (66.7%), porque
> esa medición solo cuenta lo cubierto por los tests *dentro* del
> propio paquete `evaluacion` y no ve que varios de sus métodos
> públicos (`GetEstado`, `EstaFinalizado`, `CalcularPromedio`,
> `Finalizar`) ya están cubiertos por las pruebas externas en
> `docente/docente_externo_test.go`. La cifra de 100% con
> `-coverpkg=./...` es la medición correcta del proyecto completo.

### 5. Demostración de que el acceso a un método privado no compila

Como evidencia adicional, se descomentó temporalmente la línea
`doc.validarEvaluacion(*eval)` dentro de
`docente_externo_test.go` (paquete externo `docente_test`) y se corrió
`go vet ./...`. El compilador rechazó el código con:

```
docente/docente_externo_test.go:82:6: doc.validarEvaluacion undefined
(type *docente.Docente has no field or method validarEvaluacion)
```

Esto confirma en la práctica que Go aplica la encapsulación a nivel de
compilador: un identificador en minúscula (no exportado) es
literalmente invisible fuera de su paquete. La línea se dejó comentada
en el archivo final como documentación de esta prueba.

## Análisis detallado de encapsulación e interfaces

### 1. ¿Cómo se relacionan las estructuras `Docente` y `Evaluacion`?

La relación es de **composición / agregación "uno a muchos"**: un
`Docente` mantiene un slice privado `evaluaciones []evaluacion.Evaluacion`.
No es herencia (Go no tiene herencia de structs), sino que `Docente`
contiene y administra un conjunto de valores `Evaluacion`. El acceso a
ese slice está controlado completamente por el paquete `docente`: para
agregar una evaluación hay que pasar por `AgregarEvaluacion`, que a su
vez aplica una regla de negocio (`validarEvaluacion`) antes de aceptar
el dato. Para leerlas, solo existe `ObtenerEvaluaciones` (que filtra y
devuelve solo las finalizadas) y `TotalEvaluaciones`. En ningún momento
el paquete externo puede tocar el slice `evaluaciones` directamente, lo
que evita que se inserten evaluaciones inconsistentes (por ejemplo, de
otro docente, o sin ID) sin pasar por la validación.

### 2. ¿Qué ventajas ofrece la interfaz `Evaluable`?

La interfaz desacopla el "qué se puede hacer" del "cómo está
implementado". Cualquier función que reciba un parámetro de tipo
`interfaces.Evaluable` puede operar sobre cualquier estructura que
implemente esos tres métodos (`GetEstado`, `EstaFinalizado`,
`CalcularPromedio`), sin necesitar saber si por debajo hay una
`Evaluacion` u otro tipo futuro (por ejemplo, una `AutoEvaluacion` o
una `EvaluacionPorPares`). Esto facilita: (a) escribir código genérico
y reusable, por ejemplo un reporte que reciba `[]interfaces.Evaluable`
y calcule estadísticas sin importar el tipo concreto; (b) hacer testing
con dobles/mocks que implementen la interfaz; y (c) extender el sistema
en el futuro sin romper el código existente, siempre que el nuevo tipo
cumpla el contrato.

### 3. ¿Cómo se prueba la implementación de una interfaz?

En Go esto se prueba de dos formas complementarias, y ambas se usaron
en `TestInterfaceEvaluable`:

* **En tiempo de compilación**, con una aserción de tipo del estilo
  `var _ interfaces.Evaluable = eval`. Si `*Evaluacion` no
  implementara los tres métodos de la interfaz con la firma exacta,
  esta línea por sí sola haría fallar la compilación, sin necesidad de
  ejecutar nada.
* **En tiempo de ejecución**, asignando el valor concreto a una
  variable declarada con el tipo de la interfaz (`var ev
  interfaces.Evaluable = eval`) y luego invocando los métodos a través
  de esa variable, verificando que el comportamiento (valores
  devueltos) sea el esperado. Esto confirma no solo que el tipo
  "calza" en la interfaz, sino que su comportamiento real es correcto.

### 4. ¿Qué métodos deben ser públicos y cuáles privados? ¿Por qué?

La regla aplicada fue: **público lo que forma parte del contrato que
otros paquetes necesitan usar; privado todo detalle de implementación
o regla de negocio que, si se expone, permitiría saltarse
validaciones o acoplar a otros paquetes con la forma interna de
hacer las cosas.**

* `calcularPromedioInterno()` (privado): es el algoritmo concreto del
  promedio. Si mañana cambia (por ejemplo, a un promedio ponderado),
  no debe romper a nadie que use `CalcularPromedio()` desde afuera.
* `agregarComentario()` (privado): hace la inserción "cruda", sin
  validar. Exponerlo permitiría a cualquier paquete externo insertar
  comentarios vacíos o saltarse reglas futuras.
* `validarEvaluacion()` (privado, en `Docente`): contiene la regla de
  negocio de qué evaluación es aceptable para un docente. Si fuera
  público, cualquiera podría llamarla sin después insertar el dato
  (inconsistencia entre "validar" y "agregar"), o peor, alguien podría
  asumir que validar implica agregar.
* `CalcularPromedio()`, `AgregarComentarioPublico()`,
  `AgregarEvaluacion()`, `ObtenerEvaluaciones()`,
  `TotalEvaluaciones()`, `GetEstado()`, `EstaFinalizado()`,
  `Finalizar()` (públicos): son las únicas operaciones que el resto
  del sistema necesita para usar estos tipos de forma segura y
  completa, sin tener que conocer los detalles internos.

En resumen: los métodos públicos son la **API estable** del paquete;
los privados son el **cómo**, libre de cambiar internamente sin
afectar a quien consume el paquete.

### 5. ¿Cómo afecta la encapsulación a la relación entre Docente y Evaluación?

La encapsulación hace que la relación `Docente -> []Evaluacion` sea
**segura por construcción**. Al ser `evaluaciones` un campo privado del
struct `Docente`, ningún paquete externo puede hacer
`doc.evaluaciones = append(...)` directamente; está obligado a usar
`AgregarEvaluacion`, que internamente llama a `validarEvaluacion`. Esto
garantiza que toda evaluación dentro de un `Docente` cumple siempre la
regla de negocio (pertenece a ese docente y tiene ID), sin tener que
confiar en que cada llamador externo recuerde validar antes de
insertar. De la misma forma, `ObtenerEvaluaciones` no devuelve el
slice interno sino una copia filtrada, evitando que alguien mute el
estado real del docente desde afuera. La encapsulación convierte así
una relación que podría ser frágil (cualquiera mete cualquier cosa) en
una relación con invariantes garantizados por el propio paquete.

## Reflexión personal

**Dificultades:** el punto más delicado fue decidir cómo demostrar
"que no compila" al acceder a un método privado desde un test, ya que
un test que realmente no compila rompe `go test` para todo el módulo.
Se resolvió dejando la línea problemática comentada dentro del test
(`TestValidarEvaluacionPrivada`), documentando el mensaje de error
exacto del compilador, y verificándolo de forma manual descomentándola
temporalmente con `go vet` antes de revertir el cambio. También costó
un poco decidir el diseño de `Evaluacion`: el struct original solo
tenía un campo `Puntaje float64` único, así que para que
"calcular el promedio" tuviera sentido real (y no fuera simplemente
devolver el mismo número que ya se guardó) se agregó un slice privado
`puntajes []float64` que acumula calificaciones individuales, y
`Puntaje` pasó a reflejar el promedio actualizado automáticamente.

**Aprendizajes:** quedó muy claro en la práctica que en Go la
encapsulación no es una convención de estilo sino una regla del
compilador (mayúscula/minúscula), lo cual es distinto a otros lenguajes
donde `private` es más una anotación. También fue útil distinguir entre
pruebas de caja blanca (mismo paquete, pueden tocar lo privado, sirven
para validar el algoritmo interno) y pruebas de caja negra (paquete
`_test`, solo ven lo público, sirven para validar el contrato real que
verá cualquier usuario del paquete). Medir cobertura correctamente
también dejó una lección concreta: `go test ./... -cover` por sí solo
puede subestimar la cobertura real de un paquete si otro paquete
externo es quien termina ejercitando parte de su código; para eso sirve
`-coverpkg=./...`.

**Uso de IA:** se usó asistencia de IA (Claude) como apoyo para
estructurar el código, redactar los comentarios explicativos y generar
este README, partiendo de la estructura de archivos y firmas de
funciones que ya estaban definidas en el enunciado. Todo el código fue
revisado, compilado y probado localmente (`go build`, `go vet`,
`go test -v`, medición de cobertura) antes de incluirlo en el
repositorio, y se puede explicar línea por línea: qué hace cada método,
por qué es público o privado, y cómo se relacionan los tres paquetes
entre sí.
