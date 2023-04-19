package prisma

import (
	"fmt"
	"math"
	"time"
)

var AllDayOfWeek = []DayOfWeek{
	DayOfWeekMo,
	DayOfWeekTu,
	DayOfWeekWe,
	DayOfWeekTh,
	DayOfWeekFr,
	DayOfWeekSa,
	DayOfWeekSu,
}

type DateFilter struct {
	Gt  *string
	Gte *string
	Lt  *string
	Lte *string
}

func Int32Ptr(i *int) *int32 {
	if i == nil {
		return nil
	}

	i32 := int32(*i)
	return &i32
}

func IntRaw(i int) *int32 {
	return Int32Ptr(&i)
}

func Int(i int) *int {
	return &i
}

// TimeDate is a temporary method to convert from Prisma TimeDate to Go Time
func TimeDate(loc *time.Location, str string) time.Time {
	out, err := time.Parse(time.RFC3339, str)

	if err != nil {
		panic(err)
	}

	return out.In(loc)
}

// TimeString is a temporary method to convert from Go Time to Prisma DateTime
func TimeString(t time.Time) *string {
	str := t.Format(time.RFC3339)
	return &str
}

// TimeStringPtr is a temporary method to convert from Go Time to Prisma DateTime
func TimeStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	str := t.Format(time.RFC3339)
	return &str
}

func Price(f *float64) *int32 {
	if f == nil {
		return nil
	}

	v := int32(math.Round(*f*100*100) / 100)

	return &v
}

type IUser interface {
	IsIUser()
}

func (u *User) Convert() IUser {
	switch u.Type {
	case UserTypeCustomer:
		return &Customer{u}
	case UserTypeEmployee:
		return &Employee{u}
	case UserTypeManager:
		return &Manager{u}
	case UserTypeAdministrator:
		return &Administrator{u}
	default:
		panic(fmt.Sprintf("could not convert type %#v to IUser", u))
	}
}

type IStaff interface {
	IsIStaff()
}

func (u *User) ConvertStaff() IStaff {
	switch u.Type {
	case UserTypeEmployee:
		return &Employee{u}
	case UserTypeManager:
		return &Manager{u}
	default:
		panic(fmt.Sprintf("could not convert type %#v to IStaff", u))
	}
}

type IPublicStaff interface {
	IsIPublicStaff()
}

func (u *User) ConvertPublicStaff() IPublicStaff {
	switch u.Type {
	case UserTypeEmployee:
		return &Employee{u}
	case UserTypeManager:
		return &Manager{u}
	default:
		panic(fmt.Sprintf("could not convert type %#v to IPublicStaff", u))
	}
}

type Customer struct {
	*User
}

func (Customer) IsNode()      {}
func (Customer) IsIUser()     {}
func (Customer) IsICustomer() {}

type Employee struct {
	*User
}

func (Employee) IsNode()         {}
func (Employee) IsIUser()        {}
func (Employee) IsIStaff()       {}
func (Employee) IsIPublicStaff() {}
func (Employee) IsIEmployee()    {}

type Manager struct {
	*User
}

func (Manager) IsNode()         {}
func (Manager) IsIUser()        {}
func (Manager) IsIStaff()       {}
func (Manager) IsIPublicStaff() {}
func (Manager) IsIManager()     {}

func (Appointment) IsNode() {}
func (Branch) IsNode()      {}
func (Company) IsNode()     {}
func (Product) IsNode()     {}
func (Service) IsNode()     {}
func (Session) IsNode()     {}
func (Favorite) IsNode()    {}
func (News) IsNode()        {}

type Administrator struct {
	*User
}

func (Administrator) IsNode()  {}
func (Administrator) IsIUser() {}

type IReview interface {
	IsIReview()
}

func (r *Review) Convert() IReview {
	switch r.Type {
	case ReviewTypeProduct:
		return &ProductReview{r}
	case ReviewTypeService:
		return &ServiceReview{r}
	case ReviewTypeAppointment:
		return &AppointmentReview{r}
	default:
		panic(fmt.Sprintf("could not convert type %#v to IReview", r))
	}
}

type ProductReview struct {
	*Review
}

func (ProductReview) IsNode()           {}
func (ProductReview) IsIReview()        {}
func (ProductReview) IsIProductReview() {}

type ServiceReview struct {
	*Review
}

func (ServiceReview) IsNode()           {}
func (ServiceReview) IsIReview()        {}
func (ServiceReview) IsIServiceReview() {}

type AppointmentReview struct {
	*Review
}

func (AppointmentReview) IsNode()               {}
func (AppointmentReview) IsIReview()            {}
func (AppointmentReview) IsIAppointmentReview() {}
