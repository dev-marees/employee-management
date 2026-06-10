// Package docs holds the generated OpenAPI specification.
//
// This file is a minimal hand-written placeholder so the project compiles and
// serves Swagger UI out of the box. Run `make swagger` (which calls
// `swag init -g cmd/api/main.go -o docs`) to regenerate it from the handler
// annotations with the full schema.
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "Employee Management System API",
    "description": "Backend service for the EMS frontend (auth, employees, dashboard).",
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/api/v1",
  "schemes": ["http", "https"],
  "paths": {},
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  }
}`

// SwaggerInfo registers the spec with the swag runtime. Fields are overwritten
// by the generated file after `swag init`.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{"http", "https"},
	Title:            "Employee Management System API",
	Description:      "Backend service for the EMS frontend (auth, employees, dashboard).",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
