package service

import (
	"fmt"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func assembleOrder(input *gqlgen.ServiceOrderByInput) *prisma.ServiceOrderByInput {
	if input == nil {
		return nil
	}

	r := prisma.ServiceOrderByInput(fmt.Sprintf("%s_%s", input.Field, input.Direction))
	return &r
}
