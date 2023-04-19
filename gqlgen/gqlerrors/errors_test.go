package gqlerrors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func TestNewInternalError(t *testing.T) {
	err := NewInternalError("MyType", "MyCode")

	assert.Equal(t, err, &gqlerror.Error{
		Message: "MyType",
		Extensions: map[string]interface{}{
			"code": "MyCode",
			"type": "Internal",
		},
	})
}

func TestCheckUniqueConstraintError(t *testing.T) {
	field := "email"
	errUnique := fmt.Errorf(prismaUniqueConstraint + field)
	isErrUnique := IsUniqueConstraintError(errUnique, field)
	assert.Equal(t, isErrUnique, true)

	errOther := fmt.Errorf("some other error")
	isErrOther := IsUniqueConstraintError(errOther, field)
	assert.Equal(t, isErrOther, false)
}
