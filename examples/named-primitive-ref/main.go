package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

// CustomHeader is a custom type based on a primitive string type
type CustomHeader string

// CustomQueryParam is a custom type based on a primitive int type
type CustomQueryParam int

// CreateUserBody models a request body for creating a user
type CreateUserBody struct {
	// Header is a custom header type
	Header CustomHeader `json:"header" description:"A custom header parameter"`

	// Name is the user's name
	Name string `json:"name" description:"The user's name"`
}

// UpdateUserBody models a request body for updating a user
type UpdateUserBody struct {
	// AnotherHeader is the same custom header type, to demonstrate reuse
	AnotherHeader CustomHeader `json:"another_header" description:"Another use of the same custom header"`

	// Query is a custom query parameter type
	Query CustomQueryParam `json:"query" description:"A custom query parameter"`
}

// ResponseBody models the response for both operations
type ResponseBody struct {
	// Success indicates the operation was successful
	Success bool `json:"success"`

	// Message provides details about the operation
	Message string `json:"message"`
}

func main() {
	// Create a router and configure the API
	router := chi.NewMux()

	// Enable the reuse of primitive-based types
	config := huma.DefaultConfig("My API", "1.0.0")
	config.ReuseNamedPrimitiveTypes = true

	// Create the API
	api := humachi.New(router, config)

	// Register the create user operation
	huma.Register(api, huma.Operation{
		OperationID: "createUser",
		Method:      http.MethodPost,
		Path:        "/users",
		Summary:     "Create a new user",
		Description: "Creates a new user with the provided details.",
		Tags:        []string{"users"},
	}, func(ctx context.Context, req *CreateUserBody) (*ResponseBody, error) {
		// Process the request (in a real app)

		// Return a response
		return &ResponseBody{
			Success: true,
			Message: fmt.Sprintf("User created with header: %s", req.Header),
		}, nil
	})

	// Register the update user operation (using the same custom types)
	huma.Register(api, huma.Operation{
		OperationID: "updateUser",
		Method:      http.MethodPut,
		Path:        "/users/{id}",
		Summary:     "Update an existing user",
		Description: "Updates an existing user with the provided details.",
		Tags:        []string{"users"},
	}, func(ctx context.Context, req *UpdateUserBody) (*ResponseBody, error) {
		// Process the request (in a real app)

		// Return a response
		return &ResponseBody{
			Success: true,
			Message: fmt.Sprintf("User updated with header: %s and query: %d",
				req.AnotherHeader, req.Query),
		}, nil
	})

	// Print an info message
	log.Printf("Starting server on http://localhost:8888")
	log.Printf("View OpenAPI at http://localhost:8888/openapi.json")
	log.Printf("View docs at http://localhost:8888/docs")

	// Start the HTTP server
	http.ListenAndServe(":8888", router)
}
