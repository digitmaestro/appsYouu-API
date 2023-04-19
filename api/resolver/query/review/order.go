package review

import (
	"fmt"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func AssembleOrder(input *gqlgen.ReviewOrderByInput) *prisma.ReviewOrderByInput {
	if input == nil {
		return nil
	}

	r := prisma.ReviewOrderByInput(fmt.Sprintf("%s_%s", input.Field, input.Direction))
	return &r
}
