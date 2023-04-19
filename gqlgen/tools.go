package gqlgen

import (
	"github.com/steebchen/keskin-api/prisma"
)

func TimeFilter(input *DateFilter) prisma.DateFilter {
	if input == nil {
		return prisma.DateFilter{}
	}

	return prisma.DateFilter{
		Gt:  prisma.TimeStringPtr(input.Gt),
		Gte: prisma.TimeStringPtr(input.Gte),
		Lt:  prisma.TimeStringPtr(input.Lt),
		Lte: prisma.TimeStringPtr(input.Lte),
	}
}
