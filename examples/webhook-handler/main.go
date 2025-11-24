package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gstachniukrsk/govalidator"
)

// WebhookSecret is used to verify webhook signatures
const WebhookSecret = "your-webhook-secret"

// GitHubPushEvent webhook payload structures
type GitHubPushEvent struct {
	Ref        string         `json:"ref"`
	Before     string         `json:"before"`
	After      string         `json:"after"`
	Repository GitHubRepo     `json:"repository"`
	Pusher     GitHubUser     `json:"pusher"`
	Commits    []GitHubCommit `json:"commits"`
}

type GitHubRepo struct {
	Name     string     `json:"name"`
	FullName string     `json:"full_name"`
	Owner    GitHubUser `json:"owner"`
}

type GitHubUser struct {
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

type GitHubCommit struct {
	ID        string     `json:"id"`
	Message   string     `json:"message"`
	Timestamp string     `json:"timestamp"`
	Author    GitHubUser `json:"author"`
}

// StripeEvent webhook payload structures
type StripeEvent struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Created int64           `json:"created"`
	Data    StripeEventData `json:"data"`
}

type StripeEventData struct {
	Object map[string]any `json:"object"`
}

// Define validation schemas for different webhook types
var (
	// GitHub user schema
	githubUserSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("name").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("email").
			Optional().
			WithValidators(govalidator.IsStringValidator),
		govalidator.NewField("username").
			Optional().
			WithValidators(govalidator.IsStringValidator),
	)

	// GitHub repository schema
	githubRepoSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("name").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("full_name").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("owner").
			Required().
			WithSchema(githubUserSchema),
	).WithExtra(govalidator.ExtraIgnore)

	// GitHub commit schema
	githubCommitSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("id").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(40),
				govalidator.MaxLengthValidator(40),
			),
		govalidator.NewField("message").
			Required().
			WithValidators(govalidator.IsStringValidator),
		govalidator.NewField("timestamp").
			Required().
			WithValidators(govalidator.IsStringValidator),
		govalidator.NewField("author").
			Required().
			WithSchema(githubUserSchema),
	).WithExtra(govalidator.ExtraIgnore)

	// GitHub push event schema
	githubPushSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("ref").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("before").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(40),
				govalidator.MaxLengthValidator(40),
			),
		govalidator.NewField("after").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(40),
				govalidator.MaxLengthValidator(40),
			),
		govalidator.NewField("repository").
			Required().
			WithSchema(githubRepoSchema),
		govalidator.NewField("pusher").
			Required().
			WithSchema(githubUserSchema),
		govalidator.NewField("commits").
			Required().
			WithSchema(
				govalidator.Array(githubCommitSchema.Required()),
			),
	).WithExtra(govalidator.ExtraIgnore) // GitHub sends many extra fields

	// Stripe event data schema
	stripeEventDataSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("object").
			Required().
			WithValidators(govalidator.IsMapValidator),
	).WithExtra(govalidator.ExtraIgnore)

	// Stripe event schema
	stripeEventSchema = govalidator.NewSchema().WithFields(
		govalidator.NewField("id").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.MinLengthValidator(1),
			),
		govalidator.NewField("type").
			Required().
			WithValidators(
				govalidator.IsStringValidator,
				govalidator.OneOfValidator(
					"payment_intent.succeeded",
					"payment_intent.payment_failed",
					"customer.created",
					"customer.updated",
					"customer.deleted",
					"invoice.payment_succeeded",
					"invoice.payment_failed",
				),
			),
		govalidator.NewField("created").
			Required().
			WithValidators(
				govalidator.IsIntegerValidator,
				govalidator.MinFloatValidator(0),
			),
		govalidator.NewField("data").
			Required().
			WithSchema(stripeEventDataSchema),
	).WithExtra(govalidator.ExtraIgnore) // Stripe sends many extra fields
)

// verifySignature verifies the HMAC signature of the webhook
func verifySignature(payload []byte, signature string, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// handleGitHubWebhook handles GitHub webhook requests
func handleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Verify signature (in production, always verify!)
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature != "" {
		// Remove "sha256=" prefix
		signature = strings.TrimPrefix(signature, "sha256=")
		if !verifySignature(body, signature, WebhookSecret) {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Parse JSON
	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate payload structure
	valid, errs := githubPushSchema.ValidateFlat(
		context.Background(),
		data,
		govalidator.CombinedPresenter(".", ": "),
	)

	if !valid {
		log.Printf("GitHub webhook validation failed: %v", errs)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"error":  "Invalid webhook payload",
			"errors": errs,
		})
		return
	}

	// Parse into typed struct for processing
	var event GitHubPushEvent
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Failed to parse event", http.StatusInternalServerError)
		return
	}

	// Process the webhook
	log.Printf("✓ Received valid GitHub push event:")
	log.Printf("  Repository: %s", event.Repository.FullName)
	log.Printf("  Ref: %s", event.Ref)
	log.Printf("  Pusher: %s", event.Pusher.Name)
	log.Printf("  Commits: %d", len(event.Commits))

	// Send success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Webhook processed successfully")
}

// handleStripeWebhook handles Stripe webhook requests
func handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var data any
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate payload structure
	valid, errs := stripeEventSchema.ValidateFlat(
		context.Background(),
		data,
		govalidator.CombinedPresenter(".", ": "),
	)

	if !valid {
		log.Printf("Stripe webhook validation failed: %v", errs)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"error":  "Invalid webhook payload",
			"errors": errs,
		})
		return
	}

	// Parse into typed struct
	var event StripeEvent
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "Failed to parse event", http.StatusInternalServerError)
		return
	}

	// Process the webhook
	log.Printf("✓ Received valid Stripe event:")
	log.Printf("  Event ID: %s", event.ID)
	log.Printf("  Type: %s", event.Type)
	log.Printf("  Created: %s", time.Unix(event.Created, 0).Format(time.RFC3339))

	// Send success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Webhook processed successfully")
}

func main() {
	http.HandleFunc("/webhooks/github", handleGitHubWebhook)
	http.HandleFunc("/webhooks/stripe", handleStripeWebhook)

	fmt.Println("Webhook handler server starting on :8080")
	fmt.Println("\nEndpoints:")
	fmt.Println("  POST /webhooks/github - GitHub push events")
	fmt.Println("  POST /webhooks/stripe - Stripe payment events")

	fmt.Println("\nTest with curl:")
	fmt.Println("\n1. Valid GitHub push event:")
	fmt.Println(`curl -X POST http://localhost:8080/webhooks/github \
  -H "Content-Type: application/json" \
  -d '{
    "ref": "refs/heads/main",
    "before": "0000000000000000000000000000000000000000",
    "after": "1111111111111111111111111111111111111111",
    "repository": {
      "name": "my-repo",
      "full_name": "user/my-repo",
      "owner": {
        "name": "user"
      }
    },
    "pusher": {
      "name": "John Doe",
      "email": "john@example.com"
    },
    "commits": [
      {
        "id": "1111111111111111111111111111111111111111",
        "message": "Initial commit",
        "timestamp": "2024-01-01T12:00:00Z",
        "author": {
          "name": "John Doe",
          "email": "john@example.com"
        }
      }
    ]
  }'`)

	fmt.Println("\n2. Valid Stripe event:")
	fmt.Println(`curl -X POST http://localhost:8080/webhooks/stripe \
  -H "Content-Type: application/json" \
  -d '{
    "id": "evt_1234567890",
    "type": "payment_intent.succeeded",
    "created": 1704110400,
    "data": {
      "object": {
        "id": "pi_1234567890",
        "amount": 1000
      }
    }
  }'`)

	fmt.Println("\n3. Invalid Stripe event (wrong type):")
	fmt.Println(`curl -X POST http://localhost:8080/webhooks/stripe \
  -H "Content-Type: application/json" \
  -d '{
    "id": "evt_1234567890",
    "type": "invalid.event.type",
    "created": 1704110400,
    "data": {
      "object": {}
    }
  }'`)

	fmt.Println()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
