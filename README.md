# 🏦 Ualapp Fintech

Proyecto basico para simular el funcionamiento de una fintech basada en Next.js y Go

## 🔨 Tecnologias

- **React** v18.2.0
- **Go** v1.22.0
- **PostgreSQL** v16.1
- **Docker** v25.0.3
- **Docker Compose** v2.24.5
- **Make** v4.3

## 📦 Herramietas

- **Migrate CLI** v4.14.1
  - Manejar migraciones de base de datos ya que los ORM son muy lentos para una fintech, por lo que es preferible usar SQL puro
- **sqlc CLI** v1.25.0
  - SQL compiler
  - **Input**: Se escribe la consulta en SQL --> **Blackbox**: [sqlc] --> **Output**: Funciones en Golang con interfaces para poder utilizarlas y hacer consultas

## 📚 Documentacion

- [Frontend](https://dbdocs.io/valrichter/ualapp-frontend)
- [Diseño de la Base de Datos (DER)](https://dbdocs.io/valrichter/go-ualapp)

## 📌 Proyecto

- [x] CRUD de usuarios
- [x] Autenticacion y verificacion la identidad de los usuarios
- [x] Implemetacion del Registro de usuarios
- [x] Implemetacion del Login de usuarios
- [x] Hashing y verificacion de contraseñas
- [x] Encriptacion y desencriptacion con PASETO
- [x] Implementacion de la creacion de cuenta en ARS
- [x] Implementacion de las transferencias de dinero entre cuentas

## 🚀 Desarrollo

- Diseño de la tabla de Usuario
- Creacion de los archivos sql de migracion para la base de datos
- Implementacion de docker compose para levantar el servicio de postgres
- Automatizacion de comandos con Makefile para ejecutar el contenedor de postgres, crear la base de datos e insertar las tablas
- Agregado de random generators en `utils/random.go`
- Agregado de hashing de contraseñas en `utils/password.go` con la libreria bcrypt y testeo de la misma
- Testeadas todas la queries creadas con `sqlc` para la tablas `users`
- Impletancion de `Store` para conectarse a postgres
- Implemetacion del pool de conexiones para la base de datos `pgxpool.Pool` (Singleton)
- Confuguracion del server para la API
- Configuracion del archivo `app.env`
- API server recfator
- Se agrego el endopint `/list_users` para listar todos los usuarios
- - Se agrego el endopint `/create_user` para listar todos los usuarios
- Se agrego bases de datos dedicada para testeo
- Se agrego el endpoint `auth/login` para autenticar un usuario
- Se agrego token de autenticacion para el endpoint `auth/login`
- Manejo de errores de la base de datos de usuarios no existentes y de contraseñas incorrectas
- Se agrego auth middleware para autenticar el token
- Agregado de tablas `accounts`, `entries` y `transfers`
- Creacion de los archivos sqlc para las nuevas tablas
- Implementacion la api `accounts` para crear cuenta en ARS
- Implementacion de trasacciones SQL para la query de transferir dinero entre cuentas
- Implementacion de la api para transferencias de dinero entre cuentas
- Agregado de la columna `username` a la tabla `users`
- Implementacion de la api para actualizar el `username` de una cuenta

---

- Generacion de numero de cuentas
- Endpoint `/get-account-by-number` para obtener una cuenta por su numero

## 🧪 Tests

### Base de datos `ualapp`

- [x] Tests para todas las queries de la tabla `users`
- [x] Uso de go rutines para el `TestListUsers`

### Util
  
- [x] Password hashing & verification
- [ ] Random generator
