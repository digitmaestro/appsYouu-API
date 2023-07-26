package email_template

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/smtp"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/prisma"
)

var (
	errSendingEmail = errors.New("error sending email")
)

func SendEmail(
	ctx context.Context,
	prismaClient *prisma.Client,
	branchId string,
	toEmail string,
	subject string,
	body string,
) (*gqlgen.SendEmailPayload, error) {
	branch, err := prismaClient.Branch(prisma.BranchWhereUniqueInput{
		ID: &branchId,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if branch.FromEmail == nil {
		return nil, gqlerrors.NewValidationError("fromEmail missing for branch "+branchId, "MissingData")
	}

	if branch.SmtpUsername == nil {
		return nil, gqlerrors.NewValidationError("smtpUsername missing for branch "+branchId, "MissingData")
	}

	if branch.SmtpPassword == nil {
		return nil, gqlerrors.NewValidationError("smtpPassword missing for branch "+branchId, "MissingData")
	}

	if branch.SmtpSendHost == nil {
		return nil, gqlerrors.NewValidationError("smtpSendHost missing for branch "+branchId, "MissingData")
	}

	if branch.SmtpSendPort == nil {
		return nil, gqlerrors.NewValidationError("smtpSendPort missing for branch "+branchId, "MissingData")
	}
	header := make(map[string]string)
	header["From"] = *branch.FromEmail
	header["To"] = toEmail
	header["Subject"] = "=?utf-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	auth := smtp.PlainAuth("", *branch.SmtpUsername, *branch.SmtpPassword, *branch.SmtpSendHost)
	err = smtp.SendMail(*branch.SmtpSendHost+":"+*branch.SmtpSendPort, auth, *branch.FromEmail, []string{toEmail}, []byte(message))

	if err != nil {
		log.Printf("error sending email %v", err)

		return nil, errSendingEmail
	}

	return &gqlgen.SendEmailPayload{
		Status: "OK",
	}, nil
}
