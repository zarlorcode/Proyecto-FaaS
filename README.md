# FaaS. Function as a Service


## Objetivo
El objetivo de este proyecto es diseñar e implementar un sistema FaaS (Function as a Service), una solución que permita a los usuarios registrar y ejecutar funciones bajo demanda a través de una API REST.

## Descripción Operativa del Proyecto
Este sistema FaaS básico permite:
- Registro de usuarios: El proyecto permite registrar usuarios a través de basic auth de apisix.
- Registro de funciones: Los usuarios registran sus propias funciones (referencias a imágenes de Docker) de manera privada.
-  Eliminación de funciones: Los usuarios pueden eliminar las funciones que previamente han registrado
- Ejecución de funciones: Las funciones se activan mediante llamadas API y deben aceptar un único parámetro de tipo cadena. Estas funciones corren dentro de contenedores que se eliminan automáticamente tras su ejecución.

### Sobre la Autenticación  
El sistema no cuenta con un sistema de inicio de sesión tradicional. En su lugar, la autenticación se realiza mediante **Basic Auth**, gestionada por APISIX.  
- Para realizar operaciones como registrar, eliminar o ejecutar una función, los usuarios deben incluir sus credenciales en la cabecera de cada petición.  
- Este enfoque garantiza que solo los usuarios autenticados puedan interactuar con el sistema y realizar acciones.  

Esto se ha gestionado mediante **dos tipos diferentes de rutas en APISIX**, que se detallarán más adelante en su correspondiente apartado.


## Requisitos para poder ejecutar el proyecto.
- Docker 
- Tener instalado Postman (Se pueden hacer con curl pero es más fácil y visual postman)

## Configuración para la ejecución del proyecto
### 1. Clona el repositorio:
   ```bash
   git clone https://github.com/zarlorcode/Proyecto-FaaS.git 
   ```
   ```bash
   cd Proyecto-FasS
   ```
   
### 2. Construir la aplicación y lanzar la aplicación 
Se lanzan a ejecución el proxy inverso Apisix, Nats, la API que ofrece el servicio y tantos workers como queramos tener, por defecto 1. En este ejemplo se lanzan 3 workers especificando con --scale.

```bash
docker compose up --build --scale worker=3
```

### 3. Añadir más workers en tiempo de ejecución 
Si quisieramos lanzar más workers en tiempo de ejecución podríamos ejecutar este comando en terminales distintas

```bash
docker run --rm   --name worker   -v /var/run/docker.sock:/var/run/docker.sock   --network apisix   worker-image
```
    
## Probar el proyecto
Una vez con el servicio lanzado: 
- Abrimos postman
### Probar la funcionalidad de Registrar Usuario: 
Hay que hacer una petición **POST** a la siguiente url:
```bash
http://localhost:9080/register
```
Únicamente hay que navegar a la sección Auth, en tipo hay que seleccionar "Basic Auth" e introducir los credenciales ahí. Esto podría ir en el header también pero es más directo y visual si usamos el apartado Auth. 

No hay que especificar nada en el body.
A continuación, se pulsa Send.

### Posibles salidas:
#### Registro de usuario exitoso:
```bash
{
    "message": "Usuario registrado exitosamente",
    "status": "success"
}
```
#### Si el usuario ya está registrado:
```bash
{
    "message": "El usuario ya existe",
    "status": "error"
}
```

### Probar la funcionalidad de Registrar Función:
Hay que hacer una petición **POST** a la siguiente url:
```bash
http://localhost:9080/functions/register
```
En el apartado Auth nos tenemos que autenticar. Seleccionamos el tipo "Basic Auth" y usamos los credenciales que hemos registrado antes.

En el body ponemos lo siguiente:
```bash
{
    "functionName": "reverse-string",
    "dockerImage": "lzarpor/reverse-string-function"
}
```
A continuación se pulsa Send.

Esta función se ha creado para realizar las pruebas del proyecto y lo que hace es darle la vuelta a los caracteres de una string. Es decir devuelve la string al revés.

### Posibles salidas:
#### Registro de función exitosa:
```bash
{
    "message": "Función registrada exitosamente",
    "status": "success"
}
```
#### Si el usuario ya ha registrado esa función:
```bash
{
    "message": "la función ya está registrada",
    "status": "error"
}
```

### Probar la funcionalidad de Eliminar Función:
Hay que hacer una petición **POST** a la siguiente url:
```bash
http://localhost:9080/functions/deregister
```
En el apartado Auth nos tenemos que autenticar. Seleccionamos el tipo "Basic Auth" y usamos los credenciales que hemos registrado antes.

En el body ponemos lo siguiente:
```bash
{
    "functionName": "reverse-string"
}
```
A continuación se pulsa Send.

### Posibles salidas:
#### Eliminación de función exitosa:
```bash
{
    "message": "Función eliminada exitosamente",
    "status": "success"
}
```
#### Si el usuario no es propietario de una función con dicho nombre:
```bash
{
    "message": "función no encontrada",
    "status": "error"
}
```

### Probar la funcionalidad de Activar Función:

Hay que hacer una petición **POST** a la siguiente url:

```bash
    http://localhost:9080/functions/activate
```

En el apartado Auth nos tenemos que autenticar. Seleccionamos el tipo "Basic Auth" y usamos los credenciales que hemos registrado antes.

En el body ponemos lo siguiente:
```bash
{
    "functionName": "reverse-string",
    "parameter": "hello"
}
```
### Posibles salidas:
#### Si el usuario es propietario de esa función y la petición se procesa correctamente:
```bash
{
    "result": "olleh\n",
    "status": "success"
}
```
#### Si el usuario no es propietario de una función con dicho nombre:
```bash
{
    "message": "función no encontrada",
    "status": "error"
}
```



## Componentes de la Aplicación

El sistema FaaS está compuesto por varios elementos clave que trabajan en conjunto para garantizar su funcionalidad y escalabilidad:

### 1. APISIX

[Apache APISIX](https://apisix.apache.org/) es un proxy inverso que actúa como puerta de enlace para la aplicación. Este componente gestiona todas las solicitudes API y ofrece funciones esenciales para la operación segura y eficiente del sistema.  

#### Responsabilidades de APISIX
- **Gestión de solicitudes API:** Redirige las solicitudes entrantes al componente correspondiente dentro de la arquitectura.  
- **Autenticación y autorización:** Configurable para verificar la identidad de los usuarios antes de conceder acceso a funciones específicas.  

#### Integración y Configuración
Para integrar APISIX, se ha utilizado su **Dashboard** como herramienta de configuración. La persistencia de las rutas se asegura almacenándolas en un volumen basado en **etcd**.  

#### Rutas Configuradas
Se han definido dos rutas principales, cada una con un propósito específico:  

   **Ruta abierta (RutaRegister):**  
   - **Descripción:** Esta ruta es de acceso público, no requiere credenciales para ser utilizada.  
   - **Propósito:** Permitir el registro de nuevos usuarios en el sistema sin necesidad de autenticarse previamente.  
   - **Prioridad:** Tiene una prioridad más alta para asegurar que las solicitudes de registro siempre sean atendidas primero.  

  **Ruta cerrada (Ruta1):**  
   - **Descripción:** Esta ruta está protegida mediante autenticación básica (**Basic Auth**), requiriendo un nombre de usuario y una contraseña válidos.  
   - **Propósito:** Gestionar el resto de las operaciones del FaaS, tales como:  
     - Registro de funciones.  
     - Eliminación de funciones.  
     - Activación de funciones.  
   - **Seguridad:** Garantiza que solo los usuarios autenticados puedan realizar acciones críticas en el sistema.  

Con esta configuración, APISIX asegura un control robusto sobre las solicitudes, diferenciando claramente entre las operaciones públicas y las protegidas.  


### 2. NATS  
[NATS](https://nats.io/) es un sistema de mensajería ligera que actúa como la columna vertebral de la comunicación entre los microservicios. Este componente desempeña un papel crucial en la arquitectura del sistema FaaS, facilitando la coordinación y el flujo de datos entre los distintos componentes.  

#### Responsabilidades de NATS  
- **Cola de mensajes:** Permite la comunicación asincrónica entre el servidor de la API y los workers.  
- **Almacenamiento clave-valor:** Se utiliza como base de datos para almacenar información esencial del sistema. Por ejemplo:  
  - **Bucket Key-Value:**  
    - **Key:** Una tupla que combina el nombre del usuario y el nombre de la función.  
    - **Value:** La referencia a la imagen de Docker asociada a la función.  
- **Escalabilidad:** Soporta múltiples conexiones simultáneas, manteniendo una comunicación eficiente incluso en entornos de alta carga.  

#### JetStream y Gestión de Mensajes   

##### Stream `activations`  
- **Propósito:**  
  Este stream almacena las solicitudes de activación de funciones enviadas desde el servidor de la API.  
- **Funcionamiento:**  
  - Cada solicitud contiene:  
    - Nombre del usuario que activa la función.  
    - Nombre de la función registrada.  
    - Parámetro de entrada como una cadena de texto.  
  - Se genera un ID único (`requestID`) para cada solicitud y se publica como un mensaje en el stream `activations.<requestID>`.  
  - Los workers actúan como consumidores de este stream, procesando las solicitudes en paralelo según la disponibilidad.  

##### Stream `results`  
- **Propósito:**  
  Este stream permite que los workers envíen los resultados de la ejecución de funciones de vuelta al servidor de la API.  
- **Funcionamiento:**  
  - Una vez que un worker completa la ejecución de una función, publica el resultado como un mensaje en el stream `results.<requestID>`.  
  - El servidor de la API escucha este stream y consume el resultado, devolviéndolo al usuario que activó la función.  
  
#### Consumidores y Paralelismo  

**Consumidor del stream `activations`:**  
   - Los workers están configurados como consumidores duraderos de este stream.  
   - Cada worker procesa mensajes, ejecuta la función especificada utilizando Docker y publica los resultados en el stream `results`.  
   - Para garantizar robustez, se configuran reintentos en caso de errores (por ejemplo, 5 intentos con tiempos de espera crecientes).  

**Procesamiento de Mensajes:**  
   - El mensaje de activación se descompone en sus componentes (`username`, `functionName`, `parameter`).  
   - El worker ejecuta la función usando un contenedor Docker, asegurando que el resultado sea enviado a través de `stdout`.  
   - En caso de error, el worker publica un mensaje de error en el stream `results`.  

**Consumidor del stream `results`:**  
   - La API se suscribe a mensajes en `results.<requestID>`.  
   - Espera de forma síncrona hasta recibir el resultado, que es devuelto al usuario.  
   - Si no se recibe un mensaje dentro de un tiempo límite predefinido (por ejemplo, 40 segundos), se notifica un error de tiempo de espera.  

#### Ejemplo de Flujo Completo  

**Activación de la Función:**  
   - La API publica un mensaje en `activations.<requestID>` con los detalles de la solicitud.  

**Procesamiento por el Worker:**  
   - Un worker disponible consume el mensaje, extrae los datos y ejecuta la función especificada mediante Docker.  
   - El resultado de la ejecución (o un error, si ocurre) se publica en `results.<requestID>`.  

**Recepción del Resultado:**  
   - La API escucha el mensaje de resultado en `results.<requestID>` y lo procesa para devolverlo al usuario final.  

#### Ventajas de esta Configuración  
- **Paralelismo y Escalabilidad:**  
  - Los workers pueden escalar horizontalmente para manejar mayor carga.  
  - El uso de consumidores asegura que los mensajes sean procesados de manera eficiente sin interferencias.  

- **Manejo de Errores:**  
  - Los reintentos y el manejo de errores en los workers aseguran la fiabilidad del sistema incluso ante fallos puntuales.  


### 3. API Server
El servidor de la API es el núcleo lógico de la aplicación. Sus funciones principales son:  
- **Gestión de usuarios:**  
  - Registro de nuevos usuarios.  
  - Autenticación y autorización según el método implementado.
- **Registro y gestión de funciones:**  
  - Permite a los usuarios registrar, listar y eliminar funciones (referencias a imágenes de Docker).  
  - Verifica que solo los usuarios propietarios puedan gestionar sus funciones.  
- **Coordinación:** Recibe solicitudes API, valida los parámetros y comunica las tareas a los trabajadores a través de NATS.  

### 4. Workers  

Los Workers son los encargados de procesar las funciones solicitadas por los usuarios, ejecutándolas en contenedores Docker y gestionando la comunicación con el sistema mediante NATS y JetStream.  

#### Principales Responsabilidades  

**Ejecución de Contenedores:**  
   - Utilizan la imagen de Docker especificada por el usuario para crear y ejecutar un contenedor.  
   - Configuran el contenedor para recibir un único parámetro de tipo cadena, que es proporcionado en la solicitud de activación.  
   - Ejecutan la función y capturan tanto la salida (`stdout`) como cualquier error (`stderr`).  

**Gestión de Mensajes:**  
   - Los Workers están suscritos al stream `activations` de JetStream, donde reciben solicitudes para activar funciones.  
   - Procesan las solicitudes, ejecutan las funciones, y publican los resultados en el stream `results`.  

**Gestión de Recursos:**  
   - Cada contenedor se elimina inmediatamente después de su ejecución para optimizar el uso de recursos.  
   - Los Workers pueden escalar horizontalmente, lo que permite manejar múltiples solicitudes simultáneamente en entornos de alta carga.  

#### Flujo de Trabajo  

**Conexión a NATS:**  
   - Los Workers se conectan al servidor NATS, implementando lógica de reintento para garantizar la conexión incluso si el servidor no está disponible inicialmente.  

**Suscripción al Stream `activations`:**  
   - Configuran un consumidor duradero para el stream `activations`, lo que asegura que los mensajes no procesados persistan en caso de reinicios o fallos.  
   - Cada mensaje contiene los datos necesarios para la activación:  
     - Usuario (`username`)  
     - Nombre de la función (`functionName`)  
     - Parámetro (`parameter`)  

**Procesamiento de Mensajes:**  
   - Los Workers extraen los datos del mensaje y ejecutan la función especificada utilizando Docker.  
   - Los resultados (o errores) se publican en el stream `results.<requestID>`.  

**Publicación de Resultados:**  
   - La respuesta, que puede ser la salida de la función o un mensaje de error, se envía al servidor de la API mediante el stream `results`.  

#### Configuración del Consumidor  

- **Durabilidad:**  
  Los consumidores son configurados como duraderos para asegurar la persistencia de los mensajes.  

- **Manejo de Errores:**  
  - Cada mensaje tiene un máximo de 5 intentos de entrega en caso de fallos en su procesamiento.  
  - Se implementan intervalos de backoff entre intentos (por ejemplo, 5 y 10 segundos).  

#### Ejemplo de Ejecución  

**Recepción de Solicitud:**  
   - Un mensaje es recibido desde `activations.<requestID>` con los datos del usuario, función y parámetro.  

**Ejecución en Docker:**  
   - El Worker ejecuta el siguiente comando:  
     ```bash
     docker run --rm <functionName> <parameter>
     ```  
   - Captura la salida estándar y cualquier error producido durante la ejecución.  

**Publicación del Resultado:**  
   - Si la ejecución es exitosa, el resultado es publicado en `results.<requestID>`.  
   - En caso de error, un mensaje de error detallado es enviado al mismo stream.  

#### Ventajas del Diseño  

- **Escalabilidad:**  
  - El sistema permite añadir más Workers según sea necesario, distribuyendo la carga y mejorando el rendimiento.  

- **Resiliencia:**  
  - El manejo de errores y los consumidores duraderos aseguran que las solicitudes no se pierdan incluso en condiciones adversas.  

- **Eficiencia en la Gestión de Recursos:**  
  - Los contenedores son eliminados inmediatamente después de la ejecución, minimizando el uso de recursos.  

#### Código Relacionado  
El siguiente fragmento de código representa la implementación de un Worker:  

```go
// Fragmento destacado del Worker
func processFunction(workerMsgsId, functionName, parameter string) (string, error) {
    log.Printf("[%s] PROCESANDO la función %s", workerMsgsId, functionName)
    cmd := exec.Command("docker", "run", "--rm", functionName, parameter)
    var out, stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        return "", fmt.Errorf("error ejecutando la función: %s", stderr.String())
    }
    return out.String(), nil
}

```

# Resumen de Preguntas Clave Resueltas

## 1. ¿Cómo se implementa la arquitectura de microservicios?
El sistema se organiza en una arquitectura de microservicios compuesta por los siguientes componentes:
- **APISIX:** Un proxy inverso que gestiona la entrada de solicitudes API y autentica usuarios.
- **NATS:** Un sistema de mensajería que conecta los componentes de manera asincrónica.
- **API Server:** El núcleo lógico que coordina las operaciones del sistema.
- **Workers:** Procesan las funciones registradas ejecutándolas en contenedores Docker.

## 2. ¿Cómo realizan los Workers su trabajo?
- Los Workers consumen mensajes del stream `activations` en NATS, que contienen los detalles de las solicitudes.
- Ejecutan funciones en contenedores Docker, procesando un único parámetro de entrada y capturando el resultado de la ejecución.
- Publican los resultados en el stream `results` para que el servidor de la API pueda enviarlos al usuario.

## 3. ¿Cómo se escalan los Workers?
- Los Workers pueden escalarse horizontalmente añadiendo más instancias. Esto se puede hacer al iniciar el sistema con `--scale worker=N` en Docker Compose o lanzándolos manualmente con un comando de Docker.
- El uso de NATS como sistema de mensajería asegura que los nuevos Workers se integren automáticamente y comiencen a procesar tareas pendientes.

## 4. ¿Cómo escala la base de datos?
- La base de datos se implementa utilizando el almacenamiento de clave-valor de NATS (JetStream). Esto permite que los datos se distribuyan y repliquen en entornos de alta disponibilidad.
- NATS garantiza que las claves estén disponibles de manera eficiente incluso bajo alta carga.

## 5. ¿Cómo se configura el sistema para adaptarse a cambios en la carga?
- **Workers:** Se pueden añadir o quitar dinámicamente para manejar incrementos o disminuciones en la carga de trabajo.
- **Streams en NATS:** JetStream maneja colas y persistencia de mensajes, lo que permite absorber picos de solicitudes.
- **APISIX:** Configurable para manejar el equilibrio de carga y proporcionar autenticación segura.

## 6. ¿Cómo se gestiona el acceso a servicios externos en la implementación?
- Cada función registrada hace referencia a una imagen de Docker personalizada que se ejecuta en un contenedor aislado. 
- Los Workers son responsables de configurar y ejecutar estas imágenes de forma segura, sin exponer información sensible ni requerir acceso directo a los servicios externos desde el sistema central.

## 7. ¿Qué ventajas y limitaciones presenta el sistema?
### Ventajas:
- **Escalabilidad:** Los componentes del sistema, como los Workers y NATS, son altamente escalables.
- **Modularidad:** La arquitectura basada en microservicios permite un mantenimiento y desarrollo más simples.
- **Eficiencia:** La eliminación de contenedores tras la ejecución asegura un uso óptimo de los recursos.
### Limitaciones:
- **Latencia:** La comunicación asincrónica y la creación de contenedores puede introducir retrasos en comparación con soluciones sin contenedores.
- **Complejidad Operativa:** Requiere experiencia en orquestación de contenedores y configuración de sistemas distribuidos.

## 8. ¿Qué datos adicionales podrían incluirse en las cadenas codificadas?
- **Timestamp de creación de la solicitud:** Para facilitar el monitoreo y diagnóstico de problemas.
- **Identificador de región o servidor:** Para implementar estrategias de despliegue multi-región.
- **Metadatos de la función:** Por ejemplo, la versión de la función o parámetros de configuración adicionales.

## 9. ¿Algún otro aspecto importante?
- **Seguridad:** El uso de contenedores garantiza un entorno aislado para ejecutar las funciones, minimizando riesgos de seguridad.
- **Monitoreo:** El sistema puede mejorarse añadiendo herramientas de monitoreo como Prometheus para supervisar el rendimiento de los componentes.
- **Timeouts configurables:** Los tiempos de espera son críticos para evitar que las solicitudes queden pendientes indefinidamente.

## 10. Comparación con otros sistemas FaaS de código abierto 
| Característica            | Este Sistema          | OpenFaaS                  | Knative                    |
|---------------------------|-----------------------|---------------------------|----------------------------|
| Escalabilidad             | Manual/Docker Compose| Automática mediante K8s   | Automática mediante K8s    |
| Almacenamiento de Estado  | NATS (JetStream)      | Configurable (DB/Filesystem)| Configurable (DB/Filesystem)|
| Complejidad de Configuración | Moderada             | Alta                      | Muy Alta                  |
| Recursos por ejecución    | Docker (por función)  | Docker (por función)      | Kubernetes Pod            |
