package razorpay

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omkar273/codegeeky/internal/api/dto"
	"github.com/omkar273/codegeeky/internal/config"
	"github.com/omkar273/codegeeky/internal/httpclient"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/razorpay/razorpay-go"
)

type RazorpayProvider struct {
	client         *httpclient.Client
	razorpayClient *razorpay.Client
}

func NewRazorpayProvider(client *httpclient.Client, config *config.Configuration) *RazorpayProvider {
	razorpayClient := razorpay.NewClient(config.Razorpay.APIKey, config.Razorpay.APISecret)
	return &RazorpayProvider{
		client:         client,
		razorpayClient: razorpayClient,
	}
}

func (r *RazorpayProvider) ProviderName() types.PaymentGatewayProvider {
	return types.PaymentGatewayProviderRazorpay
}

func (r *RazorpayProvider) SupportedFeatures() []types.PaymentGatewayFeatures {
	return []types.PaymentGatewayFeatures{
		types.PaymentGatewayFeaturesWebhooks,
		types.PaymentGatewayFeaturesUPI,
		types.PaymentGatewayFeaturesRefunds,
		types.PaymentGatewayFeaturesPayments,
		types.PaymentGatewayFeaturesPayouts,
		types.PaymentGatewayFeaturesPayments,
	}
}

func (r *RazorpayProvider) Initialize(ctx context.Context) error {
	// verify credentials
	// TODO: Implement this
	return nil
}

func (r *RazorpayProvider) ProcessWebhook(ctx context.Context, payload []byte, headers map[string]string) (*dto.WebhookResult, error) {
	var event map[string]interface{}
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}

	eventType := event["event"].(string)

	return &dto.WebhookResult{
		EventName: eventType,
		EventID:   event["id"].(string),
		Payload:   event,
		Headers:   headers,
		Raw:       event,
	}, nil
}

func (r *RazorpayProvider) CreatePaymentRequest(ctx context.Context, input *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	// create order
	order, err := r.razorpayClient.Order.Create(map[string]interface{}{
		"amount":          input.Amount,
		"currency":        input.Currency,
		"payment_capture": true,
		"notes": map[string]string{
			"notes": "notes",
		},
	}, map[string]string{
		"Content-Type": "application/json",
	})

	fmt.Printf("Razorpay Order: %+v\n", order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// return payment details
	return &dto.PaymentResponse{
		ProviderPaymentID: order["id"].(string),
		RedirectURL:       order["short_url"].(string),
		Status:            "success",
		Raw:               map[string]any{"order_id": order["id"].(string)},
	}, nil
}

func (r *RazorpayProvider) VerifyPaymentStatus(ctx context.Context, providerPaymentID string) (*dto.PaymentStatus, error) {
	// Call Razorpay API to verify payment
	return &dto.PaymentStatus{
		Status:            "success",
		ProviderPaymentID: providerPaymentID,
		Raw:               map[string]any{"payment_id": providerPaymentID},
	}, nil
}
