package models

// Error ...
type Error struct {
	Message string `json:"message"`
}

const (
	WrongDateMessage = "Incorrect Date"
	WrongInfoMessage = "Incorrect Data"
	AlreadyAdded     = "Already have"

	NotFoundMessage   = "Data not found"
	NotCreatedMessage = "Data not created"
	NotUpdatedMessage = "Data not updated"
	NotDeletedMessage = "Data not deleted"
	NotAddedMessage   = "Data not added"
	InternalMessage   = "Something went wrong"
	NotAvailable      = "Not available"
	WrongPassword     = "Password must be at least 8 long"
	EmailAlreadyInUse = "Email already in use, please use another email address"

	TokenExpired = "Token expired"
)
