package gateway

import "errors"

var (
	ErrInvalidEmail                 = errors.New("invalid email")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrInvalidUsername              = errors.New("invalid username")
	ErrInvalidSessionID             = errors.New("invalid session id")
	ErrInvalidClientID              = errors.New("invalid client id")
	ErrInvalidVerificationCode      = errors.New("invalid verification code")
	ErrInvalidGigID                 = errors.New("invalid gig id")
	ErrInvalidFreelancerID          = errors.New("invalid freelancer id")
	ErrInvalidTitle                 = errors.New("invalid title")
	ErrInvalidDescription           = errors.New("invalid description")
	ErrInvalidCategoryID            = errors.New("invalid category id")
	ErrInvalidCurrency              = errors.New("invalid currency")
	ErrInvalidPackageTier           = errors.New("invalid package tier")
	ErrInvalidPackageDescription    = errors.New("invalid package description")
	ErrInvalidPackageDeliveryDays   = errors.New("invalid package delivery days")
	ErrInvalidPackagePriceCents     = errors.New("invalid package price cents")
	ErrInvalidQuestionContent       = errors.New("invalid question content")
	ErrInvalidMediaUpload           = errors.New("invalid media upload")
	ErrInvalidMediaRef              = ErrInvalidMediaUpload
	ErrInvalidPackageCount          = errors.New("invalid package count")
	ErrInvalidGigState              = errors.New("invalid gig state")
	ErrInvalidOrderConnectionID     = errors.New("invalid order connection id")
	ErrInvalidOrderBuyerEmail       = errors.New("invalid order buyer email")
	ErrInvalidOrderCurrency         = errors.New("invalid order currency")
	ErrInvalidOrderPackage          = errors.New("invalid order package")
	ErrInvalidOrderPrice            = errors.New("invalid order price")
	ErrInvalidOrderGigID            = errors.New("invalid order gig id")
	ErrInvalidOrderBuyerID          = errors.New("invalid order buyer id")
	ErrInvalidOrderTitle            = errors.New("invalid order title")
	ErrInvalidFreelancerOnboarding  = errors.New("invalid freelancer onboarding")
	ErrGigNotFound                  = errors.New("gig not found")
	ErrGigDraftIncomplete           = errors.New("gig draft is incomplete")
	ErrGigAlreadyPublished          = errors.New("gig already published")
	ErrFailedToCreateGig            = errors.New("failed to create gig")
	ErrFailedToUpdateGig            = errors.New("failed to update gig")
	ErrFailedToPublishGig           = errors.New("failed to publish gig")
	ErrFailedToGetGig               = errors.New("failed to get gig")
	ErrFailedToReplaceGigPackages   = errors.New("failed to replace gig packages")
	ErrFailedToReplaceGigQuestions  = errors.New("failed to replace gig questions")
	ErrFailedToReplaceGigMedia      = errors.New("failed to replace gig media")
	ErrFailedToCreateOrder          = errors.New("failed to create order")
	ErrRegistrationNotCompleted     = errors.New("registration is not completed")
	ErrRegistrationAlreadyClaimed   = errors.New("registration tokens already claimed")
	ErrFailedToStartRegistration    = errors.New("failed to start registration")
	ErrFailedToVerifyEmail          = errors.New("failed to verify email")
	ErrFailedToCompleteRegistration = errors.New("failed to complete registration")
)

// RegistrationConflictError reports that registration cannot continue because
// either an in-progress session already exists or the owned username/email is
// already bound by the underlying services.
type RegistrationConflictError struct {
	State         string
	UsernameTaken bool
	EmailTaken    bool
}

// Error implements the error interface for registration conflict responses.
func (e *RegistrationConflictError) Error() string {
	return "registration conflict"
}

// GigConflictError reports that a gig operation conflicts with the current
// draft state or ownership constraints.
type GigConflictError struct {
	State string
}

// Error implements the error interface for gig conflict responses.
func (e *GigConflictError) Error() string {
	return "gig conflict"
}
