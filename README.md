# Okapi Example

A simple example demonstrating Okapi API

[Github: https://github.com/jkaninda/okapi](https://github.com/jkaninda/okapi)

## Prerequisites

- Go installed
- Git installed

## Features

- Basic Okapi implementation example
- Ready-to-run code structure
- Minimal dependencies

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

Visit [`http://localhost:8080`](http://localhost:8080) to see the response:

```json
{"message": "Welcome to Okapi!"}
```

Visit [`http://localhost:8080/docs/`](http://localhost:8080/docs/) to see the documentation

## Project Structure

```
.
├── main.go          # Main application file
├── middleware.go    # Middleware file
├── go.mod           # Go module file
└── README.md        # Project documentation
```

## License

[MIT](LICENSE) - Feel free to use and modify this example.

