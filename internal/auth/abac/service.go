package abac

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/omkar273/codegeeky/internal/domain/auth"
	"github.com/omkar273/codegeeky/internal/logger"
)

// Service provides Attribute-Based Access Control functionality
type Service interface {
	// Evaluate evaluates all applicable policies for an access request
	Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error)

	// RegisterPolicy adds a new policy to the evaluation chain
	RegisterPolicy(policy Policy) error

	// UnregisterPolicy removes a policy from the evaluation chain
	UnregisterPolicy(policyName string) error

	// GetPolicies returns all registered policies
	GetPolicies() []Policy

	// RegisterAttributeProvider adds an attribute provider
	RegisterAttributeProvider(provider AttributeProvider) error

	// SetPolicyCombiner sets the strategy for combining policy decisions
	SetPolicyCombiner(combiner PolicyCombiner)
}

// Policy interface for ABAC policies
type Policy interface {
	// Evaluate returns a decision for the given access request
	Evaluate(ctx context.Context, request *auth.AccessRequest) (Decision, error)

	// GetName returns the policy name for identification
	GetName() string

	// GetPriority returns the policy priority (higher = more important)
	GetPriority() int

	// Applies checks if this policy is applicable to the request
	Applies(request *auth.AccessRequest) bool
}

// Decision represents the result of a policy evaluation
type Decision struct {
	Allow    bool                   `json:"allow"`
	Reason   string                 `json:"reason"`
	Priority int                    `json:"priority"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// AttributeProvider interface for loading dynamic attributes
type AttributeProvider interface {
	// LoadUserAttributes loads additional attributes for a user
	LoadUserAttributes(ctx context.Context, userID string) (map[string]interface{}, error)

	// LoadResourceAttributes loads additional attributes for a resource
	LoadResourceAttributes(ctx context.Context, resourceType, resourceID string) (map[string]interface{}, error)

	// GetName returns the provider name
	GetName() string
}

// PolicyCombiner interface for combining multiple policy decisions
type PolicyCombiner interface {
	// Combine takes multiple policy decisions and returns a final decision
	Combine(decisions []Decision) Decision

	// GetName returns the combiner name
	GetName() string
}

type service struct {
	logger             *logger.Logger
	policies           []Policy
	attributeProviders []AttributeProvider
	policyCombiner     PolicyCombiner
	decisionCache      map[string]CachedDecision
	cacheTTL           time.Duration
	mu                 sync.RWMutex
}

// CachedDecision represents a cached policy decision
type CachedDecision struct {
	Decision  Decision
	ExpiresAt time.Time
}

// NewService creates a new ABAC service with default policies
func NewService(logger *logger.Logger) Service {
	s := &service{
		logger:             logger,
		policies:           make([]Policy, 0),
		attributeProviders: make([]AttributeProvider, 0),
		policyCombiner:     NewAllMustAllowCombiner(),
		decisionCache:      make(map[string]CachedDecision),
		cacheTTL:           5 * time.Minute, // Cache decisions for 5 minutes
	}

	// Register default policies
	s.registerDefaultPolicies()

	return s
}

// Evaluate evaluates all applicable policies for an access request
func (s *service) Evaluate(ctx context.Context, request *auth.AccessRequest) (bool, error) {
	// Check cache first
	if cached, found := s.getCachedDecision(request); found {
		s.logger.Debugw("ABAC decision served from cache",
			"user_id", request.Subject.UserID,
			"action", request.Action,
			"resource_type", request.Resource.Type,
			"allow", cached.Allow)
		return cached.Allow, nil
	}

	// Load additional attributes if needed
	if err := s.enrichRequestWithAttributes(ctx, request); err != nil {
		s.logger.Errorw("Failed to load attributes for ABAC evaluation",
			"error", err,
			"user_id", request.Subject.UserID)
		return false, fmt.Errorf("failed to load attributes: %w", err)
	}

	// Evaluate applicable policies
	decisions := make([]Decision, 0)

	for _, policy := range s.policies {
		if !policy.Applies(request) {
			s.logger.Debugw("Policy not applicable",
				"policy", policy.GetName(),
				"user_id", request.Subject.UserID,
				"action", request.Action)
			continue
		}

		decision, err := policy.Evaluate(ctx, request)
		if err != nil {
			s.logger.Errorw("Policy evaluation failed",
				"policy", policy.GetName(),
				"error", err,
				"user_id", request.Subject.UserID)
			// Continue with other policies instead of failing completely
			continue
		}

		decision.Priority = policy.GetPriority()
		decisions = append(decisions, decision)

		s.logger.Debugw("Policy evaluated",
			"policy", policy.GetName(),
			"decision", decision.Allow,
			"reason", decision.Reason,
			"user_id", request.Subject.UserID)
	}

	// If no applicable policies, allow (RBAC has already been checked)
	if len(decisions) == 0 {
		finalDecision := Decision{
			Allow:  true,
			Reason: "No applicable ABAC policies",
		}
		s.cacheDecision(request, finalDecision)
		return true, nil
	}

	// Combine policy decisions
	finalDecision := s.policyCombiner.Combine(decisions)

	// Cache the decision
	s.cacheDecision(request, finalDecision)

	s.logger.Infow("ABAC evaluation completed",
		"user_id", request.Subject.UserID,
		"action", request.Action,
		"resource_type", request.Resource.Type,
		"resource_id", request.Resource.ID,
		"allow", finalDecision.Allow,
		"reason", finalDecision.Reason,
		"policies_evaluated", len(decisions))

	return finalDecision.Allow, nil
}

// RegisterPolicy adds a new policy to the evaluation chain
func (s *service) RegisterPolicy(policy Policy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate policy names
	for _, existingPolicy := range s.policies {
		if existingPolicy.GetName() == policy.GetName() {
			return fmt.Errorf("policy with name '%s' already exists", policy.GetName())
		}
	}

	s.policies = append(s.policies, policy)

	s.logger.Infow("ABAC policy registered",
		"policy", policy.GetName(),
		"priority", policy.GetPriority())

	return nil
}

// UnregisterPolicy removes a policy from the evaluation chain
func (s *service) UnregisterPolicy(policyName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, policy := range s.policies {
		if policy.GetName() == policyName {
			// Remove policy from slice
			s.policies = append(s.policies[:i], s.policies[i+1:]...)

			s.logger.Infow("ABAC policy unregistered", "policy", policyName)
			return nil
		}
	}

	return fmt.Errorf("policy with name '%s' not found", policyName)
}

// GetPolicies returns all registered policies
func (s *service) GetPolicies() []Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]Policy, len(s.policies))
	copy(result, s.policies)
	return result
}

// RegisterAttributeProvider adds an attribute provider
func (s *service) RegisterAttributeProvider(provider AttributeProvider) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate provider names
	for _, existingProvider := range s.attributeProviders {
		if existingProvider.GetName() == provider.GetName() {
			return fmt.Errorf("attribute provider with name '%s' already exists", provider.GetName())
		}
	}

	s.attributeProviders = append(s.attributeProviders, provider)

	s.logger.Infow("ABAC attribute provider registered",
		"provider", provider.GetName())

	return nil
}

// SetPolicyCombiner sets the strategy for combining policy decisions
func (s *service) SetPolicyCombiner(combiner PolicyCombiner) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.policyCombiner = combiner

	s.logger.Infow("ABAC policy combiner set",
		"combiner", combiner.GetName())
}

// enrichRequestWithAttributes loads additional attributes from providers
func (s *service) enrichRequestWithAttributes(ctx context.Context, request *auth.AccessRequest) error {
	s.mu.RLock()
	providers := make([]AttributeProvider, len(s.attributeProviders))
	copy(providers, s.attributeProviders)
	s.mu.RUnlock()

	for _, provider := range providers {
		// Load user attributes
		userAttrs, err := provider.LoadUserAttributes(ctx, request.Subject.UserID)
		if err != nil {
			s.logger.Warnw("Failed to load user attributes",
				"provider", provider.GetName(),
				"user_id", request.Subject.UserID,
				"error", err)
			continue // Continue with other providers
		}

		// Merge attributes (existing attributes take precedence)
		for key, value := range userAttrs {
			if _, exists := request.Subject.Attributes[key]; !exists {
				request.Subject.Attributes[key] = value
			}
		}

		// Load resource attributes
		if request.Resource.ID != "" {
			resourceAttrs, err := provider.LoadResourceAttributes(ctx, request.Resource.Type, request.Resource.ID)
			if err != nil {
				s.logger.Warnw("Failed to load resource attributes",
					"provider", provider.GetName(),
					"resource_type", request.Resource.Type,
					"resource_id", request.Resource.ID,
					"error", err)
				continue
			}

			// Merge resource attributes
			for key, value := range resourceAttrs {
				if _, exists := request.Resource.Attributes[key]; !exists {
					request.Resource.Attributes[key] = value
				}
			}
		}
	}

	return nil
}

// getCachedDecision retrieves a cached decision if available and not expired
func (s *service) getCachedDecision(request *auth.AccessRequest) (Decision, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cacheKey := s.buildCacheKey(request)
	cached, exists := s.decisionCache[cacheKey]

	if !exists || time.Now().After(cached.ExpiresAt) {
		return Decision{}, false
	}

	return cached.Decision, true
}

// cacheDecision stores a decision in cache
func (s *service) cacheDecision(request *auth.AccessRequest, decision Decision) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cacheKey := s.buildCacheKey(request)
	s.decisionCache[cacheKey] = CachedDecision{
		Decision:  decision,
		ExpiresAt: time.Now().Add(s.cacheTTL),
	}

	// Clean expired entries periodically (simple cleanup)
	if len(s.decisionCache) > 1000 { // Arbitrary limit
		s.cleanExpiredCacheEntries()
	}
}

// buildCacheKey creates a cache key for a request
func (s *service) buildCacheKey(request *auth.AccessRequest) string {
	return fmt.Sprintf("%s:%s:%s:%s",
		request.Subject.UserID,
		request.Action,
		request.Resource.Type,
		request.Resource.ID)
}

// cleanExpiredCacheEntries removes expired cache entries
func (s *service) cleanExpiredCacheEntries() {
	now := time.Now()
	for key, cached := range s.decisionCache {
		if now.After(cached.ExpiresAt) {
			delete(s.decisionCache, key)
		}
	}
}

// registerDefaultPolicies registers the built-in policies
func (s *service) registerDefaultPolicies() {
	defaultPolicies := []Policy{
		NewEnrollmentBasedAccessPolicy(),
		NewOwnershipPolicy(),
		NewTimeBasedAccessPolicy(),
		NewProgressBasedPolicy(),
	}

	for _, policy := range defaultPolicies {
		if err := s.RegisterPolicy(policy); err != nil {
			s.logger.Errorw("Failed to register default policy",
				"policy", policy.GetName(),
				"error", err)
		}
	}
}
