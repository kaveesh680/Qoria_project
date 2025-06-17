# 📊 ABT Dashboard API

This project is a high-performance analytics backend built using Go. It exposes RESTful endpoints to serve business insights for ABT Corporation, powered by transactional data from a MySQL database.

---

## 📌 Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [API Endpoints](#api-endpoints)
- [Technologies Used](#technologies-used)
- [Prerequisites](#prerequisites)
- [Setup Instructions](#setup-instructions)
- [Running the Server](#running-the-server)
- [Running Tests](#running-tests)
- [Swagger Documentation](#swagger-documentation)
- [Frontend Integration](#frontend-integration)
- [Test Coverage Report](#test-coverage-report)

---

## 📖 Project Overview

The ABT Dashboard API provides RESTful access to:

- Country-level revenue metrics
- Monthly sales volumes
- Top-selling products with stock
- Top revenue-generating regions

The API supports pagination and is optimized for fast analytics querying.

---

## 🏗️ Architecture

- Clean architecture with handler, service, and repository layers
- MySQL-based backend
- Swagger/OpenAPI documentation

---

## 🔗 API Endpoints

| Endpoint                      | Method | Description                                  |
| ----------------------------- | ------ | -------------------------------------------- |
| `/ping`                       | GET    | Health check                                 |
| `/v1/metrics/country-revenue` | GET    | Get country-wise revenue and transactions    |
| `/v1/metrics/monthly-sales`   | GET    | Get highest monthly sales volume             |
| `/v1/metrics/top-products`    | GET    | Get top 20 frequently purchased products     |
| `/v1/metrics/top-regions`     | GET    | Get top 30 regions by revenue and item sales |

---

## 💡 Technologies Used

- Go (Golang)
- MySQL
- chi router
- Swaggo (Swagger generation)

---

## ✅ Prerequisites

- Go 1.18 or higher
- MySQL 8+
- A dataset loaded into a table called `transactions` (via [CSV Import Script](../csv-importer/README.md))

---

## 🛠️ Setup Instructions

1. Clone the repository:

```bash
git clone https://github.com/your-username/abt-dashboard-api.git
cd abt-dashboard-api
```

2. Set up MySQL and load the transactions data.
3. Edit DB credentials in `main.go` or use environment variables.

---

## ▶️ Running the Server

```bash
go mod tidy
go mod download
go get
go mod vendor
go run main.go
```

The server will run at: `http://localhost:8080`

---

## ✅ Running Tests

```bash
go test ./... -coverprofile=coverage.out
```

To view the coverage as HTML:

```bash
go tool cover -html=coverage.out -o coverage.html
```

Then open `coverage.html` in your browser to see detailed test coverage.

---

## 📄 Swagger Documentation

Once the server is running, access API docs at:

```
http://localhost:8080/swagger/index.html
```

---

## 🧩 Frontend Integration

A basic `frontend.html` file is available that fetches data from the backend and visualizes it using charts and tables. Make sure CORS is enabled in the backend to allow local access.

---

## 📊 Test Coverage Report

You can generate a visual HTML test coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` to explore which parts of the code are covered by tests.

---

## 📬 Contact

For any questions or support, please open an issue or contact the maintainer.

---

ABT Corporation © 2025