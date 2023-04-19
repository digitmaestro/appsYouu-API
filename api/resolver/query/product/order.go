package product

import (
	"fmt"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func assembleOrder(input *gqlgen.ProductOrderByInput) *prisma.ProductOrderByInput {
	if input == nil {
		return nil
	}

	r := prisma.ProductOrderByInput(fmt.Sprintf("%s_%s", input.Field, input.Direction))
	return &r
}
