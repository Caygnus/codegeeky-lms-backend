package subscriber

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omkar273/codegeeky/internal/config"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/logger"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/omkar273/codegeeky/internal/webhook/publisher"
)

type RazorpayReceiver struct {
	config    *config.Configuration
	publisher publisher.WebhookPublisher
	logger    *logger.Logger
}

func NewRazorpayReceiver(
	config *config.Configuration,
	publisher publisher.WebhookPublisher,
	logger *logger.Logger,
) *RazorpayReceiver {
	return &RazorpayReceiver{
		config:    config,
		publisher: publisher,
		logger:    logger,
	}
}

func (r *RazorpayReceiver) HandleWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		r.logger.Errorw("failed to read webhook body", "error", err)
		c.Error(
			ierr.NewError("failed to read webhook body").
				WithHint(err.Error()).
				Mark(ierr.ErrBadRequest),
		)
		return
	}

	// Validate signature
	signature := c.Request.Header.Get("X-Razorpay-Signature")
	if !r.validateSignature(body, signature) {
		r.logger.Warnw("invalid webhook signature", "signature", signature)
		c.Error(
			ierr.NewError("invalid webhook signature").
				WithHint(signature).
				Mark(ierr.ErrUnauthorized),
		)
		return
	}

	// Parse webhook payload
	var payload types.RazorpayWebhookPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		r.logger.Errorw("failed to parse webhook payload", "error", err)

		c.Error(
			ierr.NewError("failed to parse webhook payload").
				WithHint(err.Error()).
				Mark(ierr.ErrBadRequest),
		)
		return
	}

	// Convert to internal event format
	event := r.convertToInternalEvent(payload)

	// Publish to internal event bus
	if err := r.publisher.PublishWebhook(c.Request.Context(), event); err != nil {
		r.logger.Errorw("failed to publish webhook event", "error", err)
		c.Error(
			ierr.NewError("failed to publish webhook event").
				WithHint(err.Error()).
				Mark(ierr.ErrInternal),
		)
		return
	}

	r.logger.Infow("razorpay webhook processed successfully",
		"event", payload.Event,
		"entity", payload.Entity,
	)

	c.Status(http.StatusOK)
}

func (r *RazorpayReceiver) validateSignature(payload []byte, signature string) bool {
	expectedSignature := r.generateSignature(payload)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func (r *RazorpayReceiver) generateSignature(payload []byte) string {
	// TODO: Use the actual razorpay webhook secret
	// h := hmac.New(sha256.New, []byte(r.config.WebhookSecret))
	h := hmac.New(sha256.New, []byte(r.config.Secrets.EncryptionKey))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func (r *RazorpayReceiver) convertToInternalEvent(payload types.RazorpayWebhookPayload) *types.WebhookEvent {
	return &types.WebhookEvent{
		ID:        fmt.Sprintf("rzp_%d", time.Now().UnixNano()),
		EventName: fmt.Sprintf("razorpay.%s", payload.Event),
		UserID:    &payload.Account.ID, // Map to your user ID logic
		Payload:   json.RawMessage(mustMarshal(payload)),
		Timestamp: time.Now(),
	}
}

// mustMarshal is a helper function to safely marshal the payload
func mustMarshal(payload types.RazorpayWebhookPayload) []byte {
	data, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal razorpay payload: %v", err))
	}
	return data
}
