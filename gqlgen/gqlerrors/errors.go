package gqlerrors

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const prismaUniqueConstraint = "graphql: A unique constraint would be violated on User. Details: Field name = "

// IsUniqueConstraintError returns true if there is a unique constraint error
func IsUniqueConstraintError(err error, field string) bool {
	msg := prismaUniqueConstraint + field
	isConstraintError := err.Error() == msg
	return isConstraintError
}

const PrismaNotFound = "query returned no result"

// NewFormatNodeError returns a formatted GraphQL error if it matches a specific error message from the DB
func NewFormatNodeError(err error, node string) error {
	if err.Error() == PrismaNotFound {
		return NewNotFoundError(node)
	}

	return err
}

// NewNotFoundError returns a NotFoundError
func NewNotFoundError(node string) error {
	return &gqlerror.Error{
		Message:    "No node found with id " + node,
		Extensions: NewExtensions("Internal", "NoSuchNode"),
	}
}

// NewValidationError returns a formatted GraphQL error for form validation errors
func NewValidationError(message string, code string) error {
	return &gqlerror.Error{
		Message:    message,
		Extensions: NewExtensions("Validation", code),
	}
}

// NewInternalError returns a formatted GraphQL error for internal errors
func NewInternalError(message string, code string) error {
	return &gqlerror.Error{
		Message:    message,
		Extensions: NewExtensions("Internal", code),
	}
}

const InvalidUserType = "user type is invalid"

// NewPermissionError returns a formatted GraphQL error for permissions
func NewPermissionError(message string) error {
	return NewInternalError("unauthorized access. details: "+message, "Permission")
}

// NewVerificationError returns formatted GraphQL error for user not activated
func NewVerificationError(message string, code string) error {
	return &gqlerror.Error{
		Message:    message,
		Extensions: NewExtensions("Verification", code),
	}
}
