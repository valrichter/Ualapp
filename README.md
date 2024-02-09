# 🏦 Ualapp Fintech <img alt="Nextjs" width="30" src="https://devicon-website.vercel.app/api/nextjs/original.svg?color=%23303030" /><img alt="Go" width="30" src="https://devicon-website.vercel.app/api/go/plain.svg?color=%2300ACD7" /><img alt="postgresql" width="30" src="https://devicon-website.vercel.app/api/postgresql/plain.svg?color=%23336791" /><img alt="postgresql" width="35" src="https://devicon-website.vercel.app/api/docker/plain.svg?color=%23019BC6" />

## 🔨 Tecnologias usadas

- **Next.js** v14.1.0
  - Frontend
- **Go** v1.22.0 linux/amd64
  - Backend
- **Postgres** v16.1-alpine3.19
  - Almacenamiento
- **Migrate CLI** v4.14.1
  - Manejar migraciones de base de datos ya que los ORM son muy lentos para una fintech, por lo que es preferible usar SQL puro

## 📦 Herramietas

- **sqlc** v1.25.0
  - SQL compiler
  - **Input**: Se escribe la consulta en SQL ---> **Blackbox**: [sqlc] ---> **Output**: Funciones en Golang con interfaces para poder utilizarlas y hacer consultas

## 🚀 Proyecto

- CRUD de Usuarios
- Autenticacion
  - Verficar la identidad de los usuarios
- Implemetacion del Registry de usuarios
- Implemetacion del Login de usuarios
- Hashing y verificacion de contrasenias
- Encriptacion y desencriptacion con JWT

## Informacion

[Data Base Design DOCS](https://dbdocs.io/valrichter/go-ualapp)

## Implementacion

- Diseño de la tabla de Usuario
- Creacion de los archivos sql de migracion para la base de datos
- Implementacion de docker compose para levantar el servicio de postgres
- Automatizacion de comandos con Makefile para ejecutar el contenedor de postgres, crear la base de datos e insertar las tablas
