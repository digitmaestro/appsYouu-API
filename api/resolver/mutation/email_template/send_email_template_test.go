package email_template

import (
	"context"
	"sync"
	"testing"

	"github.com/steebchen/keskin-api/prisma"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func TestSendEmailTemplate(t *testing.T) {
	cconfig, err := prisma.NewConfig()
	panicIf(err)

	c, err := prisma.NewClient(cconfig)
	panicIf(err)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		SendEmailTemplate(context.TODO(),
			c, "register",
			"clkjlvegj001j0839fww97vt1",
			"rahat_murtaza@outlook.com",
			prisma.GenderMale,
			"Rahat",
			"Murtaza",
			prisma.Str("Some date"),
			prisma.Str("Some date"),
			prisma.Str("Some token"),
			prisma.Str("Some activationToken"))

		wg.Done()
	}()

	wg.Wait()
}
