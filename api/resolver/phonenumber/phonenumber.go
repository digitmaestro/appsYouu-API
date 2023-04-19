package phonenumber

import (
	"fmt"

	"github.com/steebchen/keskin-api/gqlgen"
)

func Convert(p *string) *gqlgen.PhoneNumber {
	if p == nil {
		return nil
	}

	return &gqlgen.PhoneNumber{
		Raw:  *p,
		Href: fmt.Sprintf("tel:%s", *p),
	}
}
