package types

// lifecycle of enrollment
// pending -> enrolled -> completed -> refund
//
//		   -> failed
//	       -> cancelled
type EnrollmentStatus string

const (
	EnrollmentStatusPending   = "pending"
	EnrollmentStatusEnrolled  = "enrolled"
	EnrollmentStatusCompleted = "completed"
	EnrollmentStatusRefunded  = "refunded"
	EnrollmentStatusCancelled = "cancelled"
	EnrollmentStatusFailed    = "failed"
)
