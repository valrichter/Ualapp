# 🏦 Ualapp Fintech

Pryecto de prueba de concepto de una fintech basada en Next.js y Go

## 🔨 Tecnologias

- **Next.js** v14.1.0
  - Frontend
- **Go** v1.22.0 linux/amd64
  - Backend
- **Postgres** v16.1-alpine3.19
  - Almacenamiento

## 📦 Herramietas

- **Migrate** v4.14.1
  - Manejar migraciones de base de datos ya que los ORM son muy lentos para una fintech, por lo que es preferible usar SQL puro
- **sqlc** v1.25.0
  - SQL compiler
  - **Input**: Se escribe la consulta en SQL ---> **Blackbox**: [sqlc] ---> **Output**: Funciones en Golang con interfaces para poder utilizarlas y hacer consultas

## 📌 Proyecto

- CRUD de usuarios
- Autenticacion y verificacion la identidad de los usuarios
- Implemetacion del Registro de usuarios
- Implemetacion del Login de usuarios
- Hashing y verificacion de contraseñas
- Encriptacion y desencriptacion con JWT

## 📚 Documentacion

- [Data Base Design (DER)](https://dbdocs.io/valrichter/go-ualapp)

## 🚀 Desarrollo

- Diseño de la tabla de Usuario
- Creacion de los archivos sql de migracion para la base de datos
- Implementacion de docker compose para levantar el servicio de postgres
- Automatizacion de comandos con Makefile para ejecutar el contenedor de postgres, crear la base de datos e insertar las tablas
