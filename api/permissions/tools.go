package permissions

import (
	"github.com/steebchen/keskin-api/prisma"
)

func in(userType prisma.UserType, types []prisma.UserType) bool {
	for _, t := range types {
		if t == userType {
			return true
		}
	}

	return false
}
