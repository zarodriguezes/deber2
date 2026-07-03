# DocenteEvaluacion

Este proyecto fue realizado en **Go** con el objetivo de poner en práctica conceptos como la **encapsulación**, las **interfaces** y la relación entre estructuras. La aplicación permite administrar evaluaciones de docentes y demuestra cómo organizar un proyecto utilizando varios paquetes.

## Estructura del proyecto

```text
DocenteEvaluacion/
├── docente/
│   ├── docente.go
│   ├── docente_test.go            # Pruebas internas
│   └── docente_externo_test.go    # Pruebas externas
├── evaluacion/
│   ├── evaluacion.go
│   └── evaluacion_test.go
├── interfaces/
│   └── evaluable.go
├── go.mod
└── README.md
```

## Comandos ejecutados

### Compilación

```bash
go build ./...
```

El proyecto compiló correctamente y no presentó errores.

### Revisión del código

```bash
go vet ./...
gofmt -l .
```

Con estos comandos se comprobó que el código no tenía advertencias y que todos los archivos estaban correctamente formateados.

### Ejecución de las pruebas

```bash
go test ./... -v
```

Todas las pruebas se ejecutaron correctamente. Tanto el paquete `docente` como `evaluacion` pasaron todos los casos de prueba sin inconvenientes. El paquete `interfaces` no incluye pruebas porque únicamente contiene la definición de una interfaz.

### Cobertura

Para medir la cobertura total se ejecutaron los siguientes comandos:

```bash
go test -coverpkg=./... -coverprofile=cover.out ./...
go tool cover -func=cover.out
```

El resultado obtenido fue una cobertura del **100 %** en todas las funciones del proyecto.

Cabe mencionar que al ejecutar únicamente `go test ./... -cover` el porcentaje del paquete `evaluacion` puede aparecer más bajo. Esto ocurre porque varias funciones son utilizadas desde las pruebas del paquete `docente`, por lo que la opción `-coverpkg=./...` refleja de mejor manera la cobertura real del proyecto completo.

### Comprobación de la encapsulación

Para verificar que la encapsulación funcionaba correctamente, se descomentó temporalmente una llamada al método privado `validarEvaluacion()` desde el archivo `docente_externo_test.go`.

Al volver a ejecutar las herramientas de verificación, el compilador mostró un error indicando que ese método no podía ser utilizado desde otro paquete. Después de comprobar este comportamiento, la línea volvió a comentarse para que el proyecto siguiera compilando normalmente.

## Análisis

### Relación entre Docente y Evaluacion

La estructura `Docente` almacena varias evaluaciones mediante un slice privado. Esto significa que ningún paquete externo puede modificar directamente esa información.

Para agregar una nueva evaluación es obligatorio utilizar el método `AgregarEvaluacion()`, el cual realiza primero una validación antes de guardarla. De esta forma se evita que se ingresen datos incorrectos o inconsistentes.

Cuando se necesitan consultar las evaluaciones, también se utilizan métodos públicos que controlan la información que se devuelve, manteniendo protegidos los datos internos.

### Uso de la interfaz `Evaluable`

La interfaz `Evaluable` permite trabajar con cualquier estructura que implemente los métodos definidos en ella, sin importar cómo esté construida internamente.

Esto hace que el código sea más flexible y fácil de ampliar en el futuro. Además, facilita la reutilización de funciones y la realización de pruebas utilizando otros tipos que implementen la misma interfaz.

### Verificación de la interfaz

La implementación de la interfaz se comprobó de dos maneras.

Primero, durante la compilación, verificando que la estructura `Evaluacion` cumpliera con todos los métodos requeridos.

Después, durante las pruebas, se creó una variable del tipo `Evaluable` utilizando una instancia de `Evaluacion` y se comprobó que todos sus métodos funcionaran correctamente.

### Métodos públicos y privados

Durante el desarrollo se decidió mantener privados aquellos métodos que solamente forman parte del funcionamiento interno del programa, por ejemplo:

* `validarEvaluacion()`
* `calcularPromedioInterno()`
* `agregarComentario()`

Estos métodos no necesitan ser utilizados por otros paquetes, por lo que es mejor mantenerlos ocultos para evitar modificaciones incorrectas.

En cambio, métodos como `AgregarEvaluacion()`, `ObtenerEvaluaciones()`, `CalcularPromedio()` o `Finalizar()` sí forman parte de la funcionalidad que otros paquetes necesitan utilizar, por lo que fueron declarados como públicos.

### Importancia de la encapsulación

La encapsulación ayuda a mantener el programa más seguro y organizado. Como el slice de evaluaciones es privado, ningún paquete externo puede agregar o modificar datos directamente.

Toda modificación debe realizarse mediante los métodos públicos, que verifican que la información sea válida antes de guardarla. Esto evita errores y hace que el estado del programa se mantenga consistente.

## Reflexión personal

Este proyecto me ayudó a entender mejor cómo funciona la encapsulación en Go. Antes conocía la teoría, pero al implementarla pude ver realmente cómo el compilador impide acceder a métodos o atributos privados desde otros paquetes.

Una de las partes que más tiempo tomó fue demostrar ese comportamiento sin dañar el proyecto. Para hacerlo, descomenté temporalmente una línea que llamaba a un método privado desde un paquete externo. Al compilar apareció el error esperado, comprobando que la encapsulación realmente funciona. Después volví a comentar esa línea para que todo siguiera funcionando.

También fue interesante diseñar la estructura `Evaluacion`. Al principio solo almacenaba un puntaje, pero luego se decidió guardar varios puntajes y calcular automáticamente el promedio. De esa manera el ejemplo tenía más sentido y permitía trabajar mejor con métodos públicos y privados.

Otra cosa que aprendí fue la diferencia entre las pruebas internas y las externas. Las internas permiten acceder a la lógica privada para comprobar que funciona correctamente, mientras que las externas solo utilizan la parte pública del paquete, simulando el uso que tendría cualquier otro programador.

Finalmente, también entendí que la cobertura depende de cómo se ejecuten las pruebas. Al principio pensé que un porcentaje menor significaba que faltaban pruebas, pero después comprendí que algunas funciones estaban siendo utilizadas desde otro paquete y que era necesario utilizar `-coverpkg` para obtener una medición más completa.

## Uso de inteligencia artificial

Durante el desarrollo utilicé Claude como apoyo para resolver algunas dudas sobre la estructura del proyecto y para obtener una primera versión del código y de la documentación.

Sin embargo, antes de entregar el trabajo revisé todo el código, ejecuté las pruebas, comprobé que compilara correctamente y verifiqué los resultados de cobertura. Además, me aseguré de entender qué hacía cada método, por qué algunos eran públicos y otros privados, y cómo interactúan entre sí los diferentes paquetes del proyecto.
