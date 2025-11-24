package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gstachniukrsk/govalidator"
)

// UserRegistration represents a user registration request
type UserRegistration struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Role     string `json:"role"`
}

// ProductCreation represents a product creation request
type ProductCreation struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Categories  []string `json:"categories"`
	Tags        []string `json:"tags"`
}

// Define validation schemas
var (
	// Email regex pattern
	emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// User registration schema
	userSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("username").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(3),
				govalidator.MaxLengthValidator(20),
			),
		govalidator.NewField("email").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.RegexpValidator(*emailPattern),
			),
		govalidator.NewField("password").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(8),
			),
		govalidator.NewField("age").
			Optional().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(13),
				govalidator.MaxFloatValidator(120),
			),
		govalidator.NewField("role").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.OneOfValidator("admin", "user", "guest"),
			),
	).WithExtra(govalidator.ExtraForbid)

	// Product creation schema
	productSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("name").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(3),
				govalidator.MaxLengthValidator(100),
			),
		govalidator.NewField("description").
			Optional().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MaxLengthValidator(1000),
			),
		govalidator.NewField("price").
			Required().
			WithValidators(
				govalidator.NumberValidator,
				govalidator.MinFloatValidator(0.01),
			),
		govalidator.NewField("categories").
			Required().
			WithSchema(
				govalidator.Array(
					govalidator.NewSchema(
						govalidator.IsStringValidator,
						govalidator.MinLengthValidator(1),
					).Required(),
				).Required(),
			),
		govalidator.NewField("tags").
			Optional().
			WithSchema(
				govalidator.Array(
					govalidator.NewSchema(govalidator.IsStringValidator).Required(),
				),
			),
	).WithExtra(govalidator.ExtraForbid)
)

// validateRequest validates HTTP request body against a schema
func validateRequest(r *http.Request, schema *govalidator.Schema) (any, []string, error) {
	var data any
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Validate using flat error list for API responses
	valid, errs := schema.ValidateFlat(
		context.Background(),
		data,
		govalidator.CombinedPresenter(".", ": "),
	)

	if !valid {
		return data, errs, nil
	}

	return data, nil, nil
}

// writeJSONResponse writes a JSON response
func writeJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// handleUserRegistration handles user registration endpoint
func handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	data, validationErrors, err := validateRequest(r, userSchema)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if len(validationErrors) > 0 {
		writeJSONResponse(w, http.StatusUnprocessableEntity, map[string]any{
			"error":  "Validation failed",
			"errors": validationErrors,
		})
		return
	}

	// In a real application, you would save the user to a database
	writeJSONResponse(w, http.StatusCreated, map[string]any{
		"message": "User registered successfully",
		"data":    data,
	})
}

// handleProductCreation handles product creation endpoint
func handleProductCreation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "Method not allowed",
		})
		return
	}

	data, validationErrors, err := validateRequest(r, productSchema)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if len(validationErrors) > 0 {
		writeJSONResponse(w, http.StatusUnprocessableEntity, map[string]any{
			"error":  "Validation failed",
			"errors": validationErrors,
		})
		return
	}

	// In a real application, you would save the product to a database
	writeJSONResponse(w, http.StatusCreated, map[string]any{
		"message": "Product created successfully",
		"data":    data,
	})
}

func main() {
	http.HandleFunc("/api/users/register", handleUserRegistration)
	http.HandleFunc("/api/products", handleProductCreation)

	fmt.Println("Server starting on :8080")
	fmt.Println("\nTry these curl commands:")
	fmt.Println("\n1. Valid user registration:")
	fmt.Println(`curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john_doe","email":"john@example.com","password":"secret123","role":"user","age":25}'`)

	fmt.Println("\n2. Invalid user registration (missing fields):")
	fmt.Println(`curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"jo","email":"invalid-email","password":"short"}'`)

	fmt.Println("\n3. Valid product creation:")
	fmt.Println(`curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"High-performance laptop","price":999.99,"categories":["electronics","computers"],"tags":["new","featured"]}'`)

	fmt.Println("\n4. Invalid product creation (negative price):")
	fmt.Println(`curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","price":-100,"categories":[]}'`)

	fmt.Println()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
