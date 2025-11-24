# Webhook Handler Example

This example demonstrates how to use govalidator to validate incoming webhook payloads from external services like GitHub and Stripe, ensuring data integrity and security.

## Features

- Validates webhook payloads from multiple providers
- Demonstrates signature verification for webhook security
- Handles nested, complex payload structures
- Uses `ExtraIgnore` mode to handle evolving webhook schemas
- Validates arrays of objects (commits, line items, etc.)
- Provides detailed validation error responses
- Shows real-world webhook payload structures

## Running the Example

```bash
cd examples/webhook-handler
go run main.go
```

The server will start on `http://localhost:8080` with two webhook endpoints.

## Webhook Endpoints

### POST /webhooks/github

Handles GitHub push event webhooks with validation for:
- Repository information
- Commit data (SHA validation)
- Author/pusher information
- Reference (branch/tag) names

**Example Valid Request:**
```bash
curl -X POST http://localhost:8080/webhooks/github \
  -H "Content-Type: application/json" \
  -d '{
    "ref": "refs/heads/main",
    "before": "0000000000000000000000000000000000000000",
    "after": "1111111111111111111111111111111111111111",
    "repository": {
      "name": "my-repo",
      "full_name": "user/my-repo",
      "owner": {"name": "user"}
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
  }'
```

### POST /webhooks/stripe

Handles Stripe event webhooks with validation for:
- Event ID and type
- Timestamp validation
- Event data structure
- Allowed event types (enum validation)

**Example Valid Request:**
```bash
curl -X POST http://localhost:8080/webhooks/stripe \
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
  }'
```

**Example Invalid Request (wrong event type):**
```bash
curl -X POST http://localhost:8080/webhooks/stripe \
  -H "Content-Type: application/json" \
  -d '{
    "id": "evt_1234567890",
    "type": "invalid.event.type",
    "created": 1704110400,
    "data": {"object": {}}
  }'
```

## Validation Rules

### GitHub Push Events
- **ref**: Required, non-empty string (branch reference)
- **before**: Required, exactly 40 characters (commit SHA)
- **after**: Required, exactly 40 characters (commit SHA)
- **repository**: Required object with name, full_name, and owner
- **pusher**: Required user object
- **commits**: Required array of commit objects
- Extra fields are ignored (GitHub sends many additional fields)

### Stripe Events
- **id**: Required, non-empty string
- **type**: Required, must be one of predefined event types:
  - payment_intent.succeeded
  - payment_intent.payment_failed
  - customer.created/updated/deleted
  - invoice.payment_succeeded/failed
- **created**: Required, non-negative integer (Unix timestamp)
- **data**: Required object containing event data
- Extra fields are ignored (Stripe's schema evolves over time)

## Security Features

The example includes HMAC signature verification:

```go
signature := r.Header.Get("X-Hub-Signature-256")
if !verifySignature(body, signature, WebhookSecret) {
    http.Error(w, "Invalid signature", http.StatusUnauthorized)
    return
}
```

**Important**: Always verify webhook signatures in production!

## Key Patterns Demonstrated

1. **Schema Composition**: Building complex schemas from reusable components
2. **Array Validation**: Validating arrays of nested objects
3. **Extra Fields Handling**: Using `ExtraIgnore` for evolving APIs
4. **SHA Validation**: Using exact length validators for commit hashes
5. **Enum Validation**: Restricting event types to known values
6. **Nested Object Validation**: Multi-level object structures
7. **Security**: Signature verification before validation
8. **Error Responses**: Returning detailed validation errors to webhook senders

## Use Cases

This webhook validation pattern is essential for:
- GitHub/GitLab webhooks (CI/CD triggers)
- Stripe/PayPal payment notifications
- SendGrid/Mailgun email delivery webhooks
- Slack/Discord bot webhooks
- Custom webhook integrations
- Third-party service callbacks
- Event-driven architectures

## Production Considerations

When deploying webhook handlers:

1. **Always verify signatures** - Prevents forged requests
2. **Use HTTPS** - Encrypt webhook data in transit
3. **Implement idempotency** - Handle duplicate webhook deliveries
4. **Log validation failures** - Debug integration issues
5. **Return appropriate status codes**:
   - 200: Successfully processed
   - 400: Bad request (malformed JSON)
   - 401: Invalid signature
   - 422: Validation failed
6. **Process asynchronously** - Queue webhooks for background processing
7. **Monitor webhook endpoints** - Track success/failure rates

## Extending the Example

To add support for more webhook providers:

1. Define the payload structure as Go structs
2. Create a validation schema matching the structure
3. Use `ExtraIgnore` for flexibility with API changes
4. Add signature verification for the specific provider
5. Implement the HTTP handler following the same pattern

Example for adding Twilio webhooks:

```go
twilioWebhookSchema := govalidator.NewSchema().WithFields(
    govalidator.NewField("MessageSid").Required().WithValidators(govalidator.IsStringValidator),
    govalidator.NewField("From").Required().WithValidators(govalidator.IsStringValidator),
    govalidator.NewField("Body").Required().WithValidators(govalidator.IsStringValidator),
).WithExtra(govalidator.ExtraIgnore)
```
