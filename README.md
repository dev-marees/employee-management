# Employee Management System

A full-stack Employee Management System built with React and Golang.

## Overview

This application allows organizations to manage employees efficiently through a modern web interface and a scalable backend API.

### Features

* Employee Management
* Dashboard Analytics
* Search Employees
* Filter Employees
* Sort Employees
* Pagination
* Authentication & Authorization
* Role-Based Access Control
* Soft Delete Support
* Swagger API Documentation
* Dockerized Deployment

---

# Tech Stack

## Frontend

* React
* React Router
* Axios
* Context API / Redux Toolkit
* Vite
* CSS / Tailwind CSS

## Backend

* Golang
* Gin Framework
* PostgreSQL
* GORM
* JWT Authentication
* Swagger
* Docker

---

# Project Structure

```text
employee-management-system/

├── frontend/
│   ├── public/
│   ├── src/
│   │   ├── assets/
│   │   ├── components/
│   │   │   ├── common/
│   │   │   ├── layout/
│   │   │   └── employee/
│   │   ├── pages/
│   │   │   ├── Dashboard/
│   │   │   ├── Employees/
│   │   │   ├── EmployeeDetails/
│   │   │   ├── AddEmployee/
│   │   │   └── EditEmployee/
│   │   ├── hooks/
│   │   ├── services/
│   │   ├── routes/
│   │   ├── utils/
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── vite.config.js
│
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── auth/
│   │   ├── employee/
│   │   ├── dashboard/
│   │   ├── middleware/
│   │   ├── config/
│   │   ├── database/
│   │   └── logger/
│   │
│   ├── migrations/
│   ├── docs/
│   ├── pkg/
│   ├── go.mod
│   └── Dockerfile
│
├── docker-compose.yml
└── README.md
```

---

# Functional Requirements

## Authentication

### Register

Create a new user account.

### Login

Authenticate users and generate JWT tokens.

### Refresh Token

Generate a new access token using a valid refresh token.

---

## Employee Management

### Create Employee

Add a new employee.

### Update Employee

Modify employee details.

### Delete Employee

Soft delete employee records.

### Employee Details

View complete employee information.

### Employee Listing

View all employees with:

* Pagination
* Search
* Filters
* Sorting

---

# Dashboard

Dashboard provides quick insights.

### Metrics

* Total Employees
* Active Employees
* Inactive Employees
* Department Wise Employee Count

Example:

```json
{
  "totalEmployees": 150,
  "activeEmployees": 120,
  "inactiveEmployees": 30
}
```

---

# Database Design

## Users

| Column        | Type      |
| ------------- | --------- |
| id            | UUID      |
| name          | VARCHAR   |
| email         | VARCHAR   |
| password_hash | VARCHAR   |
| role          | VARCHAR   |
| created_at    | TIMESTAMP |
| updated_at    | TIMESTAMP |

---

## Employees

| Column        | Type      |
| ------------- | --------- |
| id            | UUID      |
| employee_code | VARCHAR   |
| first_name    | VARCHAR   |
| last_name     | VARCHAR   |
| email         | VARCHAR   |
| phone         | VARCHAR   |
| department    | VARCHAR   |
| designation   | VARCHAR   |
| salary        | DECIMAL   |
| joining_date  | DATE      |
| status        | VARCHAR   |
| created_at    | TIMESTAMP |
| updated_at    | TIMESTAMP |
| deleted_at    | TIMESTAMP |

---

# API Endpoints

## Authentication

### Register

POST /api/v1/auth/register

### Login

POST /api/v1/auth/login

### Refresh Token

POST /api/v1/auth/refresh

---

## Dashboard

### Dashboard Summary

GET /api/v1/dashboard

---

## Employees

### Get Employees

GET /api/v1/employees

Query Parameters:

```text
?page=1
&limit=10
&search=john
&department=engineering
&status=active
&sort=name
&order=asc
```

### Get Employee

GET /api/v1/employees/{id}

### Create Employee

POST /api/v1/employees

### Update Employee

PUT /api/v1/employees/{id}

### Delete Employee

DELETE /api/v1/employees/{id}

---

# API Response Format

## Success

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {}
}
```

## Error

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {}
}
```

---

# Frontend Setup

## Prerequisites

* Node.js
* npm

## Installation

```bash
cd frontend
npm install
```

## Environment Variables

Create:

```text
frontend/.env
```

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## Run Application

```bash
npm run dev
```

Application URL:

```text
http://localhost:5173
```

---

# Backend Setup

## Prerequisites

* Go 1.25+
* PostgreSQL

## Installation

```bash
cd backend
go mod tidy
```

## Environment Variables

Create:

```text
backend/.env
```

```env
APP_ENV=development

SERVER_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=employee_db

JWT_SECRET=super-secret-key
JWT_EXPIRY=24h
```

## Run Application

```bash
go run cmd/server/main.go
```

API URL:

```text
http://localhost:8080
```

---

# Docker Setup

## Build Containers

```bash
docker-compose build
```

## Start Containers

```bash
docker-compose up -d
```

## Stop Containers

```bash
docker-compose down
```

---

# Swagger Documentation

Generate Swagger docs:

```bash
swag init
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

---

# Testing

## Backend

```bash
go test ./...
```

## Frontend

```bash
npm run test
```

---

# Future Enhancements

* Redis Caching
* Audit Logging
* AWS S3 File Uploads
* Role Management
* Employee Profile Photos
* Export CSV
* Export Excel
* Email Notifications
* OpenTelemetry Tracing
* Prometheus Monitoring
* GitHub Actions CI/CD

---

# User Roles

## Admin

* Create Employee
* Update Employee
* Delete Employee
* View Dashboard

## HR

* Create Employee
* Update Employee
* View Employee Details

## Employee

* View Own Profile

---

# License

MIT License

---

# Author

Built using React, Golang, PostgreSQL, Gin, GORM, JWT, and Docker.
