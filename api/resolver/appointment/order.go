package appointment

import (
	"fmt"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func AssembleOrder(input *gqlgen.AppointmentOrderByInput) *prisma.AppointmentOrderByInput {
	if input == nil {
		return nil
	}

	r := prisma.AppointmentOrderByInput(fmt.Sprintf("%s_%s", input.Field, input.Direction))
	return &r
}
