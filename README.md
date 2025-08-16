# Okapi Example

A simple example demonstrating the Okapi API Framework

Okapi is a modern, minimalist HTTP web framework for Go, inspired by FastAPI's elegance. Designed for simplicity, performance, and developer happiness, it helps you build fast, scalable, and well-documented APIs with minimal boilerplate.

* [Okapi](https://github.com/jkaninda/okapi)
* [Source Code](https://github.com/jkaninda/okapi-example)
* [Docker Hub](https://hub.docker.com/r/jkaninda/okapi-example)

## Prerequisites

- Go installed
- Git installed

## Features

- Basic Okapi implementation example
- Okapi middlewares
- Okapi Route Definition
- Ready-to-run code structure

## Getting Started

### Clone the Repository

```shell
git clone https://github.com/jkaninda/okapi-example
cd okapi-example
```

### Install Dependencies

```shell
go mod tidy
```

### Run the Application

```shell
go run .
```

### Using Docker

```shell
docker run --rm --name okapi-example -p 8080:8080 jkaninda/okapi-example
```
Use `JWT_SIGNING_SECRET` environment variable if you want to change JWT secret, default: `supersecret`

Visit [`http://localhost:8080`](http://localhost:8080) to see the response:

```json
{"message": "Welcome to the Okapi Web Framework!"}
```

Visit [`http://localhost:8080/docs/`](http://localhost:8080/docs/) to see the documentation

## Project Structure

```
.
├── main.go          # Main application file
├── middlewares      # Middlewares package
├── controllers      # Controllers package
├── routes           # Routes package
├── models           # Models package
└── README.md        # Project documentation
```

### Swagger UI Preview

Okapi automatically generates Swagger UI for all routes:


![Okapi Swagger Interface](https://raw.githubusercontent.com/jkaninda/okapi-example/main/swagger.png)

---
## License

[MIT](LICENSE) - Feel free to use and modify this example.

