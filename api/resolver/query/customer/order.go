package customer

import (
	"fmt"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func AssembleOrder(input *gqlgen.CustomerOrderByInput) *prisma.UserOrderByInput {
	if input == nil {
		return nil
	}

	r := prisma.UserOrderByInput(fmt.Sprintf("%s_%s", input.Field, input.Direction))
	return &r
}
