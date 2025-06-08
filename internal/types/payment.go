package types

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type PaymentMode string

const (
	PaymentModeCard       PaymentMode = "card"
	PaymentModeUPI        PaymentMode = "upi"
	PaymentModeNetbanking PaymentMode = "netbanking"
)

type PaymentProvider string

const (
	PaymentProviderRazorpay PaymentProvider = "razorpay"
)

type PaymentGatewayFeatures string

const (
	PaymentGatewayFeaturesWebhooks      PaymentGatewayFeatures = "webhooks"
	PaymentGatewayFeaturesRefunds       PaymentGatewayFeatures = "refunds"
	PaymentGatewayFeaturesUPI           PaymentGatewayFeatures = "upi"
	PaymentGatewayFeaturesPayments      PaymentGatewayFeatures = "payments"
	PaymentGatewayFeaturesPayouts       PaymentGatewayFeatures = "payouts"
	PaymentGatewayFeaturesInvoices      PaymentGatewayFeatures = "invoices"
	PaymentGatewayFeaturesSubscriptions PaymentGatewayFeatures = "subscriptions"
)

type SelectionAttributes struct {
	CountryCode     string
	Currency        string
	Amount          int64
	PaymentMode     PaymentMode
	PaymentProvider PaymentProvider
}
