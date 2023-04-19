package user

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/api/resolver/user/administrator"
	"github.com/steebchen/keskin-api/api/resolver/user/customer"
	"github.com/steebchen/keskin-api/api/resolver/user/employee"
	"github.com/steebchen/keskin-api/api/resolver/user/manager"
	"github.com/steebchen/keskin-api/prisma"
)

type User struct {
	Prisma                *prisma.Client
	EmployeeResolver      *employee.Employee
	CustomerResolver      *customer.Customer
	ManagerResolver       *manager.Manager
	AdministratorResolver *administrator.Administrator
}

func New(
	client *prisma.Client,
	employee *employee.Employee,
	customer *customer.Customer,
	manager *manager.Manager,
	administrator *administrator.Administrator,
) *User {
	return &User{
		Prisma:                client,
		EmployeeResolver:      employee,
		CustomerResolver:      customer,
		ManagerResolver:       manager,
		AdministratorResolver: administrator,
	}
}

var ProviderSet = wire.NewSet(
	New,
	employee.New,
	customer.New,
	manager.New,
	administrator.New,
)
