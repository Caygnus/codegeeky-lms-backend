package types

import (
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/samber/lo"
)

type PaymentStatus string

const (
	PaymentStatusPending           PaymentStatus = "pending"
	PaymentStatusSuccess           PaymentStatus = "success"
	PaymentStatusFailed            PaymentStatus = "failed"
	PaymentStatusPendingRefund     PaymentStatus = "pending_refund"
	PaymentStatusRefunding         PaymentStatus = "refunding"
	PaymentStatusProcessing        PaymentStatus = "processing"
	PaymentStatusPartiallyRefunded PaymentStatus = "partially_refunded"
	PaymentStatusCancelled         PaymentStatus = "cancelled"
	PaymentStatusExpired           PaymentStatus = "expired"
	PaymentStatusRefunded          PaymentStatus = "refunded"
)

func (s PaymentStatus) String() string {
	return string(s)
}

func (s PaymentStatus) Validate() error {
	allowed := []PaymentStatus{
		PaymentStatusPending,
		PaymentStatusSuccess,
		PaymentStatusFailed,
	}
	if !lo.Contains(allowed, s) {
		return ierr.NewError("invalid payment status").
			WithHint("Please provide a valid payment status").
			WithReportableDetails(map[string]any{
				"allowed": allowed,
			}).
			Mark(ierr.ErrValidation)
	}
	return nil
}

type PaymentMethodType string

const (
	PaymentMethodTypeCard         PaymentMethodType = "card"
	PaymentMethodTypeUPI          PaymentMethodType = "upi"
	PaymentMethodTypeNetbanking   PaymentMethodType = "netbanking"
	PaymentMethodTypeOffline      PaymentMethodType = "offline"
	PaymentMethodTypeWallet       PaymentMethodType = "wallet"
	PaymentMethodTypeBankTransfer PaymentMethodType = "bank_transfer"
)

func (s PaymentMethodType) String() string {
	return string(s)
}

func (s PaymentMethodType) Validate() error {
	allowed := []PaymentMethodType{
		PaymentMethodTypeCard,
		PaymentMethodTypeUPI,
		PaymentMethodTypeNetbanking,
		PaymentMethodTypeOffline,
		PaymentMethodTypeWallet,
		PaymentMethodTypeBankTransfer,
	}
	if !lo.Contains(allowed, s) {
		return ierr.NewError("invalid payment method type").
			WithHint("Please provide a valid payment method type").
			WithReportableDetails(map[string]any{
				"allowed": allowed,
			}).
			Mark(ierr.ErrValidation)
	}
	return nil
}

type PaymentGatewayProvider string

const (
	PaymentGatewayProviderRazorpay PaymentGatewayProvider = "razorpay"
)

func (s PaymentGatewayProvider) String() string {
	return string(s)
}

func (s PaymentGatewayProvider) Validate() error {
	allowed := []PaymentGatewayProvider{
		PaymentGatewayProviderRazorpay,
	}
	if !lo.Contains(allowed, s) {
		return ierr.NewError("invalid payment provider").
			WithHint("Please provide a valid payment provider").
			WithReportableDetails(map[string]any{
				"allowed": allowed,
			}).
			Mark(ierr.ErrValidation)
	}
	return nil
}

type PaymentDestinationType string

const (
	PaymentDestinationTypeInternship PaymentDestinationType = "internship"
	PaymentDestinationTypeEnrollment PaymentDestinationType = "enrollment"
	PaymentDestinationTypeCourse     PaymentDestinationType = "course"
	PaymentDestinationTypeOrder      PaymentDestinationType = "order"
)

func (s PaymentDestinationType) Validate() error {
	allowed := []PaymentDestinationType{
		PaymentDestinationTypeInternship,
		PaymentDestinationTypeEnrollment,
		PaymentDestinationTypeCourse,
		PaymentDestinationTypeOrder,
	}
	if !lo.Contains(allowed, s) {
		return ierr.NewError("invalid payment destination type").
			WithHint("Please provide a valid payment destination type").
			WithReportableDetails(map[string]any{
				"allowed": allowed,
			}).
			Mark(ierr.ErrValidation)
	}
	return nil
}

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
	PaymentMethodType PaymentMethodType
	PaymentGateway    PaymentGatewayProvider
	PaymentStatus     PaymentStatus
	PaymentAmount     int64
	PaymentCurrency   string
	PaymentCountry    string
}

// PaymentFilter represents the filter for listing payments
type PaymentFilter struct {
	*QueryFilter
	*TimeRangeFilter

	PaymentIDs        []string `form:"payment_ids"`
	DestinationType   *string  `form:"destination_type"`
	DestinationID     *string  `form:"destination_id"`
	PaymentMethodType *string  `form:"payment_method_type"`
	PaymentStatus     *string  `form:"payment_status"`
	PaymentGateway    *string  `form:"payment_gateway"`
	Currency          *string  `form:"currency"`
}

// NewNoLimitPaymentFilter creates a new payment filter with no limit
func NewNoLimitPaymentFilter() *PaymentFilter {
	return &PaymentFilter{
		QueryFilter: NewNoLimitQueryFilter(),
	}
}

// Validate validates the payment filter
func (f *PaymentFilter) Validate() error {
	if f == nil {
		return nil
	}

	if err := f.QueryFilter.Validate(); err != nil {
		return err
	}

	if err := f.TimeRangeFilter.Validate(); err != nil {
		return err
	}

	return nil
}

// GetLimit implements BaseFilter interface
func (f *PaymentFilter) GetLimit() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetLimit()
	}
	return f.QueryFilter.GetLimit()
}

// GetOffset implements BaseFilter interface
func (f *PaymentFilter) GetOffset() int {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOffset()
	}
	return f.QueryFilter.GetOffset()
}

// GetSort implements BaseFilter interface
func (f *PaymentFilter) GetSort() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetSort()
	}
	return f.QueryFilter.GetSort()
}

func (f *PaymentFilter) GetOrder() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetOrder()
	}
	return f.QueryFilter.GetOrder()
}

func (f *PaymentFilter) GetStatus() string {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetStatus()
	}
	return f.QueryFilter.GetStatus()
}

func (f *PaymentFilter) GetExpand() Expand {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().GetExpand()
	}
	return f.QueryFilter.GetExpand()
}

// IsUnlimited returns true if the filter has no limit
func (f *PaymentFilter) IsUnlimited() bool {
	if f.QueryFilter == nil {
		return NewDefaultQueryFilter().IsUnlimited()
	}
	return f.QueryFilter.IsUnlimited()
}
