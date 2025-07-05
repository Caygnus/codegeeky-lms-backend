package service

import (
	"context"
	"time"

	"github.com/omkar273/codegeeky/internal/api/dto"
	domainCart "github.com/omkar273/codegeeky/internal/domain/cart"
	ierr "github.com/omkar273/codegeeky/internal/errors"
	"github.com/omkar273/codegeeky/internal/types"
	"github.com/samber/lo"
	"github.com/shopspring/decimal"
)

type CartService interface {
	CreateCart(ctx context.Context, req *dto.CreateCartRequest) (*domainCart.Cart, error)
	GetCart(ctx context.Context, id string) (*domainCart.Cart, error)
	UpdateCart(ctx context.Context, id string, req *dto.UpdateCartRequest) (*domainCart.Cart, error)
	DeleteCart(ctx context.Context, id string) error
	ListCarts(ctx context.Context, filter *types.CartFilter) ([]*domainCart.Cart, error)
	GetCartLineItems(ctx context.Context, cartId string) ([]*domainCart.CartLineItem, error)

	// line item service
	AddLineItem(ctx context.Context, req *dto.CreateCartLineItemRequest) (*domainCart.CartLineItem, error)
	RemoveLineItem(ctx context.Context, id string) error
	GetLineItem(ctx context.Context, id string) (*domainCart.CartLineItem, error)
	ListLineItems(ctx context.Context, cartId string) ([]*domainCart.CartLineItem, error)
}

type cartService struct {
	ServiceParams
}

func NewCartService(params ServiceParams) CartService {
	return &cartService{
		ServiceParams: params,
	}
}

func (s *cartService) CreateCart(ctx context.Context, req *dto.CreateCartRequest) (*domainCart.Cart, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	cart := req.ToCart(ctx)

	// if cart type is default, then we need to check if the user has a cart
	if req.Type == types.CartTypeDefault {
		cart, err := s.CartRepo.GetUserDefaultCart(ctx, cart.UserID)
		if err == nil && cart != nil {
			return nil, ierr.NewError("user already has a default cart").
				WithHint("user already has a default cart").
				WithReportableDetails(map[string]any{
					"cart_id": cart.ID,
				}).
				Mark(ierr.ErrValidation)
		}
	}

	// resolve the line items
	lineItems := make([]*domainCart.CartLineItem, 0)
	for _, lineItemReq := range req.LineItems {
		if lineItemReq.EntityType == types.CartLineItemEntityTypeInternshipBatch {

			internshipBatch, err := s.ServiceParams.InternshipBatchRepo.Get(ctx, lineItemReq.EntityID)
			if err != nil {
				return nil, err
			}

			internship, err := s.ServiceParams.InternshipRepo.Get(ctx, internshipBatch.InternshipID)
			if err != nil {
				return nil, err
			}

			if internship.Status != types.StatusPublished {
				return nil, ierr.NewError("internship is not published").
					WithHint("internship is not published").
					WithReportableDetails(map[string]any{
						"internship_id": internship.ID,
					}).
					Mark(ierr.ErrValidation)
			}

			lineItem := lineItemReq.ToCartLineItem(ctx)

			lineItem.EntityID = internship.ID
			lineItem.EntityType = types.CartLineItemEntityTypeInternshipBatch
			lineItem.CartID = cart.ID

			// PerUnitPrice should be the original/base price (internship.Price)
			// This shows customers the "regular price" before discounts
			lineItem.PerUnitPrice = internship.Price

			// Calculate subtotal: original price * quantity
			lineItem.Subtotal = internship.Price.Mul(decimal.NewFromInt(int64(lineItemReq.Quantity)))

			// Calculate discount amount: (original price - final price) * quantity
			// This reflects the internship's built-in discount
			discountPerUnit := internship.Price.Sub(internship.Total)
			lineItem.DiscountAmount = discountPerUnit.Mul(decimal.NewFromInt(int64(lineItemReq.Quantity)))

			// For now, tax amount is zero since it's not implemented in the internship model
			// In a real system, you'd calculate tax based on the final price
			lineItem.TaxAmount = decimal.Zero

			// Calculate total: subtotal - discount + tax (should equal internship.Total * quantity)
			lineItem.Total = lineItem.Subtotal.Sub(lineItem.DiscountAmount).Add(lineItem.TaxAmount)

			lineItems = append(lineItems, lineItem)
		}
	}

	// update
	cart.LineItems = lineItems

	err := s.ServiceParams.CartRepo.CreateWithLineItems(ctx, cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

// GetCart retrieves a cart by its ID
func (s *cartService) GetCart(ctx context.Context, id string) (*domainCart.Cart, error) {
	if id == "" {
		return nil, ierr.NewError("cart ID is required").
			WithHint("Please provide a valid cart ID").
			Mark(ierr.ErrValidation)
	}

	cart, err := s.ServiceParams.CartRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	// Here we are not checking if the user has permission to access this cart
	// because userId check is already enforced in the repository

	return cart, nil
}

// UpdateCart updates an existing cart
func (s *cartService) UpdateCart(ctx context.Context, id string, req *dto.UpdateCartRequest) (*domainCart.Cart, error) {
	if id == "" {
		return nil, ierr.NewError("cart ID is required").
			WithHint("Please provide a valid cart ID").
			Mark(ierr.ErrValidation)
	}

	// Get existing cart
	existingCart, err := s.ServiceParams.CartRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if !req.ExpiresAt.IsZero() {
		existingCart.ExpiresAt = lo.ToPtr(req.ExpiresAt)
	}

	if req.Metadata != nil {
		existingCart.Metadata = req.Metadata
	}

	// Update cart
	err = s.ServiceParams.CartRepo.Update(ctx, existingCart)
	if err != nil {
		return nil, err
	}

	return existingCart, nil
}

// DeleteCart deletes a cart by its ID
func (s *cartService) DeleteCart(ctx context.Context, id string) error {
	if id == "" {
		return ierr.NewError("cart ID is required").
			WithHint("Please provide a valid cart ID").
			Mark(ierr.ErrValidation)
	}

	// Get existing cart to check permissions
	existingCart, err := s.ServiceParams.CartRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	// Check if user has permission to delete this cart
	userID := types.GetUserID(ctx)
	if existingCart.UserID != userID {
		return ierr.NewError("access denied").
			WithHint("You can only delete your own carts").
			WithReportableDetails(map[string]any{
				"cart_id": id,
				"user_id": userID,
			}).
			Mark(ierr.ErrPermissionDenied)
	}

	return s.CartRepo.Delete(ctx, id)
}

// ListCarts retrieves a list of carts based on filter criteria
func (s *cartService) ListCarts(ctx context.Context, filter *types.CartFilter) ([]*domainCart.Cart, error) {
	if filter == nil {
		filter = types.NewCartFilter()
	}

	filter.UserID = types.GetUserID(ctx)

	// Validate filter
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	return s.CartRepo.List(ctx, filter)
}

// GetCartLineItems retrieves all line items for a specific cart
func (s *cartService) GetCartLineItems(ctx context.Context, cartId string) ([]*domainCart.CartLineItem, error) {
	if cartId == "" {
		return nil, ierr.NewError("cart ID is required").
			WithHint("Please provide a valid cart ID").
			Mark(ierr.ErrValidation)
	}

	// Verify cart exists and user has access
	_, err := s.CartRepo.Get(ctx, cartId)
	if err != nil {
		return nil, err
	}

	return s.CartRepo.ListCartLineItems(ctx, cartId)
}

// AddLineItem adds a new line item to a cart
func (s *cartService) AddLineItem(ctx context.Context, req *dto.CreateCartLineItemRequest) (*domainCart.CartLineItem, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var cart *domainCart.Cart
	var err error
	// TODO: This is a temporary solution to find the user's default cart
	// We need to find a better way to handle this
	// If cart ID is not provided, try to find user's default cart
	if req.CartID == "" {
		userID := types.GetUserID(ctx)
		cart, err = s.CartRepo.GetUserDefaultCart(ctx, userID)
		if err != nil {
			return nil, err
		}

		if cart == nil {
			return nil, ierr.NewError("no default cart found").
				WithHint("Please create a cart first or specify a cart ID").
				Mark(ierr.ErrValidation)
		}

		req.CartID = cart.ID
	} else {
		cart, err = s.CartRepo.Get(ctx, req.CartID)
		if err != nil {
			return nil, err
		}
	}

	// Check if user has permission to modify this cart
	userID := types.GetUserID(ctx)
	if cart.UserID != userID {
		return nil, ierr.NewError("access denied").
			WithHint("You can only add items to your own carts").
			WithReportableDetails(map[string]any{
				"cart_id": req.CartID,
				"user_id": userID,
			}).
			Mark(ierr.ErrPermissionDenied)
	}

	// Validate entity exists and is available
	if req.EntityType == types.CartLineItemEntityTypeInternshipBatch {
		internship, err := s.InternshipRepo.Get(ctx, req.EntityID)
		if err != nil {
			return nil, err
		}

		if internship == nil {
			return nil, ierr.NewError("internship not found").
				WithHint("The specified internship does not exist").
				WithReportableDetails(map[string]any{
					"internship_id": req.EntityID,
				}).
				Mark(ierr.ErrValidation)
		}

		if internship.Status != types.StatusPublished {
			return nil, ierr.NewError("internship is not available").
				WithHint("The specified internship is not currently available for enrollment").
				WithReportableDetails(map[string]any{
					"internship_id": internship.ID,
					"status":        internship.Status,
				}).
				Mark(ierr.ErrValidation)
		}
	}

	// Create line item
	lineItem := req.ToCartLineItem(ctx)
	lineItem.CartID = req.CartID

	// Set pricing based on entity type
	if req.EntityType == types.CartLineItemEntityTypeInternshipBatch {
		internship, _ := s.InternshipRepo.Get(ctx, req.EntityID)
		lineItem.PerUnitPrice = internship.Price
		lineItem.Subtotal = internship.Price.Mul(decimal.NewFromInt(int64(req.Quantity)))

		// Calculate discount amount: (original price - final price) * quantity
		discountPerUnit := internship.Price.Sub(internship.Total)
		lineItem.DiscountAmount = discountPerUnit.Mul(decimal.NewFromInt(int64(req.Quantity)))

		// For now, tax amount is zero
		lineItem.TaxAmount = decimal.Zero

		// Calculate total: subtotal - discount + tax
		lineItem.Total = lineItem.Subtotal.Sub(lineItem.DiscountAmount).Add(lineItem.TaxAmount)
	}

	// Create the line item
	err = s.CartRepo.CreateCartLineItem(ctx, lineItem)
	if err != nil {
		return nil, err
	}

	// Update cart totals
	err = s.updateCartTotals(ctx, req.CartID)
	if err != nil {
		s.Logger.Errorw("failed to update cart totals after adding line item",
			"cart_id", req.CartID, "error", err)
		// Don't fail the operation for total update failure
	}

	return lineItem, nil
}

// RemoveLineItem removes a line item from a cart
func (s *cartService) RemoveLineItem(ctx context.Context, id string) error {
	if id == "" {
		return ierr.NewError("line item ID is required").
			WithHint("Please provide a valid line item ID").
			Mark(ierr.ErrValidation)
	}

	// Get line item to check permissions
	lineItem, err := s.CartRepo.GetCartLineItem(ctx, id)
	if err != nil {
		return err
	}

	// Verify cart exists and user has access
	cart, err := s.CartRepo.Get(ctx, lineItem.CartID)
	if err != nil {
		return err
	}

	// Check if user has permission to modify this cart
	userID := types.GetUserID(ctx)
	if cart.UserID != userID {
		return ierr.NewError("access denied").
			WithHint("You can only remove items from your own carts").
			WithReportableDetails(map[string]any{
				"cart_id":      lineItem.CartID,
				"line_item_id": id,
				"user_id":      userID,
			}).
			Mark(ierr.ErrPermissionDenied)
	}

	// Delete the line item
	err = s.CartRepo.DeleteCartLineItem(ctx, id)
	if err != nil {
		return err
	}

	// Update cart totals
	err = s.updateCartTotals(ctx, lineItem.CartID)
	if err != nil {
		s.Logger.Errorw("failed to update cart totals after removing line item",
			"cart_id", lineItem.CartID, "error", err)
		// Don't fail the operation for total update failure
	}

	return nil
}

// GetLineItem retrieves a specific line item by its ID
func (s *cartService) GetLineItem(ctx context.Context, id string) (*domainCart.CartLineItem, error) {
	if id == "" {
		return nil, ierr.NewError("line item ID is required").
			WithHint("Please provide a valid line item ID").
			Mark(ierr.ErrValidation)
	}

	lineItem, err := s.CartRepo.GetCartLineItem(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify cart exists and user has access
	cart, err := s.CartRepo.Get(ctx, lineItem.CartID)
	if err != nil {
		return nil, err
	}

	// Check if user has permission to access this cart
	userID := types.GetUserID(ctx)
	if cart.UserID != userID {
		return nil, ierr.NewError("access denied").
			WithHint("You can only access line items from your own carts").
			WithReportableDetails(map[string]any{
				"cart_id":      lineItem.CartID,
				"line_item_id": id,
				"user_id":      userID,
			}).
			Mark(ierr.ErrPermissionDenied)
	}

	return lineItem, nil
}

// ListLineItems retrieves all line items for a specific cart
func (s *cartService) ListLineItems(ctx context.Context, cartId string) ([]*domainCart.CartLineItem, error) {
	return s.GetCartLineItems(ctx, cartId)
}

// updateCartTotals recalculates and updates the cart totals based on line items
func (s *cartService) updateCartTotals(ctx context.Context, cartID string) error {
	// Get all line items for the cart
	lineItems, err := s.CartRepo.ListCartLineItems(ctx, cartID)
	if err != nil {
		return err
	}

	// Calculate totals
	var subtotal, discountAmount, taxAmount, total decimal.Decimal

	for _, item := range lineItems {
		subtotal = subtotal.Add(item.Subtotal)
		discountAmount = discountAmount.Add(item.DiscountAmount)
		taxAmount = taxAmount.Add(item.TaxAmount)
		total = total.Add(item.Total)
	}

	// Get cart and update totals
	cart, err := s.CartRepo.Get(ctx, cartID)
	if err != nil {
		return err
	}

	cart.Subtotal = subtotal
	cart.DiscountAmount = discountAmount
	cart.TaxAmount = taxAmount
	cart.Total = total
	cart.UpdatedAt = time.Now().UTC()

	return s.CartRepo.Update(ctx, cart)
}
