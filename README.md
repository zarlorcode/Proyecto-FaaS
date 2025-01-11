# FaaS. Function as a Service


## Objetivo
El objetivo de este proyecto es diseñar e implementar un sistema FaaS (Function as a Service), una solución que permita a los usuarios registrar y ejecutar funciones bajo demanda a través de una API REST.

## Descripción Operativa del Proyecto
Este sistema FaaS básico permite:
- Registro de usuarios: El proyecto permite registrar usuarios a través de basic auth de apisix.
- Registro de funciones: Los usuarios registran sus propias funciones (referencias a imágenes de Docker) de manera privada.
-  Eliminación de funciones: Los usuarios pueden eliminar las funciones que previamente han registrado
- Ejecución de funciones: Las funciones se activan mediante llamadas API y deben aceptar un único parámetro de tipo cadena. Estas funciones corren dentro de contenedores que se eliminan automáticamente tras su ejecución.

Sobre la Autenticación: No se ha definido un login como tal, la autenticación de los usuarios se ha implementado con Basic Auth a través de Apisix se ha hecho de manera que cada vez que el usuario quiera realizar cualquiera de las acciones como registrar una función, eliminar una función y ejecutar una función deberá enviar sus credenciales en la cabecera de la petición para poder realizar esta operación.

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
Hay que hacer una petición POST a la siguiente url:
```bash
http://localhost:9080/register
```
Únicamente hay que navegar a la sección Auth, en tipo hay que seleccionar "Basic Auth" e introducir los credenciales ahí. Esto podría ir en el header también pero es más directo y visual si usamos el apartado Auth. 

No hay que especificar nada en el body.
A continuación se pulsa Send.

### Probar la funcionalidad de Registrar Función:
Hay que hacer una petición POST a la siguiente url:
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

### Probar la funcionalidad de Eliminar Función:
Hay que hacer una petición POST a la siguiente url:
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

### Probar la funcionalidad de Activar Función:

Hay que hacer una petición POST a la siguiente url:

```bash
    http://localhost:9080/functions/functions/activate
```

En el apartado Auth nos tenemos que autenticar. Seleccionamos el tipo "Basic Auth" y usamos los credenciales que hemos registrado antes.

En el body ponemos lo siguiente:
```bash
{
    "functionName": "reverse-string",
    "parameter": "hello"
}
```

## Componentes de la Aplicación

El sistema FaaS está compuesto por varios elementos clave que trabajan en conjunto para garantizar su funcionalidad y escalabilidad:

### 1. APISIX
[Apache APISIX](https://apisix.apache.org/) es un proxy inverso que actúa como puerta de enlace para la aplicación. Sus responsabilidades incluyen:  
- **Gestión de solicitudes API:** Redirige las solicitudes entrantes al componente correspondiente.  
- **Autenticación y autorización:** Puede configurarse para verificar la identidad de los usuarios antes de permitirles acceder al sistema.  
- **Balanceo de carga:** Distribuye las solicitudes entre los diferentes trabajadores (`Workers`) para optimizar el uso de recursos.  
- **Seguridad:** Maneja el tráfico HTTPS para proteger la comunicación entre los usuarios y el sistema.  

Para integrar y configurar APISIX SE HA HECHO MEDIANTE EL dashboard y para garantizar la persistencia de las rutas se guarda en un volumen etcd
Se han definido 2 rutas (FALTA) Explicar para qué sirve cada una de ellas.
Una que es abierta (sin necesidad de introducir credenciales) y de prioridad más alta para gestionar el register
Otra ruta que es cerrada y solicita Basic Auth con usuario y contraseña para el resto de funciones del FaaS (registrar funcion, eliminar funcion y activar función)

### 2. NATS 
[NATS](https://nats.io/) es un sistema de mensajería ligera que actúa como la columna vertebral de la comunicación entre los microservicios. Su rol incluye:  
- **Cola de mensajes:** Permite la comunicación asincrónica entre el servidor de la API y los trabajadores.  
- **Almacenamiento clave-valor:** Se utiliza como base de datos para almacenar información esencial del sistema, como registros de usuarios y funciones.  
- **Escalabilidad:** Admite la conexión simultánea de múltiples componentes, manteniendo una comunicación eficiente incluso bajo alta carga.  

También os ha servido para crear una base de datos Key-Value donde guardar en un bucket la tupla usuario/nombre de función como Key y docker image de la función como value.

Explicar jeatstream y la cola de mensajes con un consumidor para varios workers

### 3. API Server
El servidor de la API es el núcleo lógico de la aplicación. Sus funciones principales son:  
- **Gestión de usuarios:**  
  - Registro de nuevos usuarios.  
  - Autenticación y autorización según el método implementado (por ejemplo, tokens JWT).  
- **Registro y gestión de funciones:**  
  - Permite a los usuarios registrar, listar y eliminar funciones (referencias a imágenes de Docker).  
  - Verifica que solo los usuarios propietarios puedan gestionar sus funciones.  
- **Coordinación:** Recibe solicitudes API, valida los parámetros y comunica las tareas a los trabajadores a través de NATS.  

### 4. Workers
Los trabajadores (`Workers`) son los responsables de ejecutar las funciones registradas. Sus tareas incluyen:  
- **Ejecución de contenedores:**  
  - Usan la referencia de imagen de Docker proporcionada por el usuario para crear y ejecutar un contenedor.  
  - Configuran el contenedor para aceptar un único parámetro de tipo cadena.  
- **Procesamiento de resultados:**  
  - Capturan la salida del contenedor (`stdout`) y la devuelven al usuario como respuesta.  
  - Emiten mensajes de log a través de `stderr` para fines de depuración.  
- **Gestión de recursos:**  
  - Eliminan los contenedores tras la ejecución para liberar recursos.  
  - Pueden escalar horizontalmente según la carga, permitiendo múltiples ejecuciones simultáneas.  
