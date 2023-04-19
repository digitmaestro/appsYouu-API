package main

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/steebchen/keskin-api/prisma"
)

const Password = "asdfasdf"

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func hour(hour time.Duration, minute time.Duration) string {
	return time.Now().
		Truncate(time.Hour * 24).
		Add(hour * time.Hour).
		Add(-2 * time.Hour).
		Add(minute * time.Minute).
		Format(time.RFC3339)
}

func main() {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	panicIf(err)

	password := string(hashBytes)

	config, err := prisma.NewConfig()
	panicIf(err)

	c, err := prisma.NewClient(config)
	panicIf(err)
	ctx := context.Background()

	company1, err := c.CreateCompany(prisma.CompanyCreateInput{
		Name: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Hairdesign"),
				En: prisma.Str("Hairdesign"),
				Tr: prisma.Str("Hairdesign"),
			},
		},
		PwaShortName: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Hairdesign"),
				En: prisma.Str("Hairdesign"),
				Tr: prisma.Str("Hairdesign"),
			},
		},
		PwaThemeColor:      "white",
		PwaBackgroundColor: "white",
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateBranch(prisma.BranchCreateInput{
		Name: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Stuttgart"),
				En: prisma.Str("Stuttgart"),
				Tr: prisma.Str("Stuttgart"),
			},
		},
		Address:     prisma.Str("Friseurstraße 3, 71145 Stuttgart"),
		PhoneNumber: prisma.Str("+4971568435916"),
		WelcomeMessage: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Willkommen im Barbershop"),
				En: prisma.Str("Welcome to the barber shop"),
				Tr: prisma.Str(""),
			},
		},
		WebsiteUrl: "https://www.appsyouu.com",

		Company: prisma.CompanyCreateOneWithoutBranchesInput{
			Connect: &prisma.CompanyWhereUniqueInput{
				ID: &company1.ID,
			},
		},
	}).Exec(ctx)
	panicIf(err)

	genderMale := prisma.GenderMale
	activated := true

	_, err = c.CreateUser(prisma.UserCreateInput{
		FirstName:    "",
		LastName:     "Administrator",
		Email:        "admin@appsyouu.de",
		PasswordHash: password,
		Type:         prisma.UserTypeAdministrator,
		PhoneNumber:  prisma.Str(""),
		Gender:       &genderMale,
		Activated:    &activated,
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateEmailTemplate(prisma.EmailTemplateCreateInput{
		Name: "invite",
		Title: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("App-Einladung"),
				En: prisma.Str("App invite"),
				Tr: prisma.Str("App-Einladung"),
			},
		},
		Content: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Hallo {{firstName}},\n\nvielen Dank, dass Sie bei uns waren.\n\nRegistrieren Sie sich unter {{customerSubdomain}}, um Ihre zukünftigen Termine zu vereinbaren und bequem Leistungen und Produkte zu bestellen.\n\nMit freundlichen Grüßen\nIhr {{companyName}}-Team"),
				En: prisma.Str("Hello {{firstName}},\n\nThank you for choosing us.\n\nRegister at {{customerSubdomain}}, to schedule future appointments and order services and products comfortably.\n\nBest regards,\n{{companyName}}"),
				Tr: prisma.Str("Hallo {{firstName}},\n\nvielen Dank, dass Sie bei uns waren.\n\nRegistrieren Sie sich unter {{customerSubdomain}}, um Ihre zukünftigen Termine zu vereinbaren und bequem Leistungen und Produkte zu bestellen.\n\nSaygılarımızla,\n{{companyName}}"),
			},
		},
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateEmailTemplate(prisma.EmailTemplateCreateInput{
		Name: "register",
		Title: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Ihre Registrierung"),
				En: prisma.Str("Your registration information"),
				Tr: prisma.Str("Kaydınız"),
			},
		},
		Content: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("{{salutation}} {{lastName}},\n\nvielen Dank, dass Sie sich bei uns registriert haben.\n\nSie können mit unserer App in Zukunft ganz bequem Termine vereinbaren, Ihr eigenes Kundenprofil anlegen sowie vieles mehr.\n\nBitte aktivieren Sie Ihren Account unter: {{customerSubdomain}}/activate-account/{{activateToken}}\n\nBesuchen Sie uns hier: {{customerSubdomain}}\n\nWir freuen uns auf Sie.\n\nMit freundlichen Grüßen\nIhr {{companyName}}-Team"),
				En: prisma.Str("{{salutation}} {{firstName}},\n\nThank you for registering to our {{companyName}} app.\n\nWe are exited to have you join us.\n\nWith our app you can easily arrange appointments, create your own customer profile and much more.\n\nPlease click here to activate your account: {{customerSubdomain}}/activate-account/{{activateToken}}\n\nPlease click here: {{customerSubdomain}}.\n\nWe look forward to seeing you.\n\nBest regards,\n{{companyName}}"),
				Tr: prisma.Str("Merhaba {{firstname}} {{salutation}},\n\nuygulamamıza kayıt olduğunuz için teşekkürler.\n\nUygulamamızın sayesinde bizimle randevularınızı ayarlayıp ve kendi müşteri profilinizi ekliye bilirsiniz.\n\nHesabınızı etkinleştirmek için lütfen burayı tıklayın: {{customerSubdomain}}/activate-account/{{activateToken}}\n\nUygulamamıza giriş yapmak icin lütfen tıklayın: {{customerSubdomain}}\n\nSaygılarımızla,\n{{companyName}}"),
			},
		},
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateEmailTemplate(prisma.EmailTemplateCreateInput{
		Name: "appointmentReminder",
		Title: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Terminerinnerung"),
				En: prisma.Str("Appointment reminder"),
				Tr: prisma.Str("Randevunuzu unutmayın"),
			},
		},
		Content: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("{{salutation}} {{lastName}},\n\nwir möchten Sie daran erinnern, dass Sie morgen um {{appointmentTime}} Uhr einen Termin haben.\n\nWir freuen uns auf Ihren Besuch.\n\nMit freundlichen Grüßen\nIhr {{companyName}}-Team"),
				En: prisma.Str("{{salutation}} {{firstName}},\n\nWe would like to remind you that you have an appointment tomorrow at {{appointmentTime}}.\n\nThank you for choosing us.\n\nBest regards,\n{{companyName}}"),
				Tr: prisma.Str("Merhaba {{firstname}} {{salutation}},\n\nYarın saat {{appointmentTime}} randevunuzun olduğunu hatırlatmak isteriz.\n\nBizi tercih ettiğiniz için teşekkür ederiz.\n\nSaygılarımızla,\n{{companyName}}"),
			},
		},
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateEmailTemplate(prisma.EmailTemplateCreateInput{
		Name: "appointmentConfirmed",
		Title: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Ihre Terminbestätigung"),
				En: prisma.Str("Appointment confirmed"),
				Tr: prisma.Str("Randevu onayınız"),
			},
		},
		Content: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("{{salutation}} {{lastName}},\n\nIhr Termin am {{appointmentDate}} um {{appointmentTime}} Uhr wurde bestätigt.\n\nWir freuen uns auf Ihren Besuch.\n\nMit freundlichen Grüßen\nIhr {{companyName}}-Team"),
				En: prisma.Str("{{salutation}} {{firstName}},\n\nYour appointment is confirmed for the {{appointmentDate}} at {{appointmentTime}}\n\nWe look forward to your visit.\n\nBest regards,\n{{companyName}}"),
				Tr: prisma.Str("Merhaba {{firstname}} {{salutation}},\n\nRandevunuz onaylandı. Tarih: {{appointmentDate}}, saat: {{appointmentTime}}\n\nBizi tercih ettiğiniz için teşekkür ederiz.\n\nSaygılarımızla,\n{{companyName}}"),
			},
		},
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateEmailTemplate(prisma.EmailTemplateCreateInput{
		Name: "appointmentCanceled",
		Title: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Terminabsage"),
				En: prisma.Str("Appointment cacelation"),
				Tr: prisma.Str("Randevu iptaliniz"),
			},
		},
		Content: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("{{salutation}} {{lastName}},\n\nmit Bedauern müssen wir Ihren Termin am {{appointmentDate}} um {{appointmentTime}} Uhr bei uns absagen.\nBitte vereinbaren Sie einen neuen Termin mit uns.\n\nMit freundlichen Grüßen\nIhr {{companyName}}-Team"),
				En: prisma.Str("{{salutation}} {{firstName}},\n\nRegretabbly, we will need to cancel your appointment with us on {{appointmentDate}} at {{appointmentTime}}.\nPlease accept our apologies. We would be grateful if you could arrange a new appointment with us.\n\nAgain, apologies for the inconvenience caused and we very much look forward to seeing you soon.\n\nBest regards,\n{{companyName}}"),
				Tr: prisma.Str("Merhaba {{firstname}} {{salutation}},\n\nTarih: {{appointmentDate}}, saat: {{appointmentTime}}\nRandevunuzu iptal etmek zorunda kaldığımız için üzgünüz.\n\nLütfen bizimle yeni bir randevu ayarlayın.\n\nSaygılarımızla,\n{{companyName}}"),
			},
		},
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateEmailTemplate(prisma.EmailTemplateCreateInput{
		Name: "appointmentCanceledByUser",
		Title: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Terminabsage"),
				En: prisma.Str("Appointment cacelation"),
				Tr: prisma.Str("Randevu iptaliniz"),
			},
		},
		Content: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("{{salutation}} {{lastName}},\n\nihre Terminabsage ist bei uns eingegangen.\nWir bedanken uns für die Benachrichtigung.\nSie können gerne jederzeit einen neuen Termin mit uns vereinbaren.\n\nMit freundlichen Grüßen\nIhr {{companyName}}-Team"),
				En: prisma.Str("{{salutation}} {{firstName}},\n\nWe have received your appointment cancellation.\nThank you for letting us know.\nYou are welcome to arrange a new appointment with us at any time.\n\nBest regards,\n{{companyName}}"),
				Tr: prisma.Str("Merhaba {{firstname}} {{salutation}},\n\nRandevunuz iptal edilmişdir.\nBilgilendirme için teşekkürler.\nBizimle istediğiniz zaman yeni bir randevu ayarlayabilirsiniz.\n\nSaygılarımızla,\n{{companyName}}"),
			},
		},
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateEmailTemplate(prisma.EmailTemplateCreateInput{
		Name: "birthday",
		Title: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Geburtstagsgrüße"),
				En: prisma.Str("Birthday wishes"),
				Tr: prisma.Str("Doğum günü tebrikleri"),
			},
		},
		Content: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("{{salutation}} {{lastName}},\n\nwir wünschen Ihnen alles Gute und Liebe zum Geburtstag.\n\nWir freuen uns auf Ihren nächsten Besuch.\n\nMit freundlichen Grüßen\nIhr {{companyName}}-Team"),
				En: prisma.Str("{{salutation}} {{firstName}},\n\nWe wish you a very happy birthday.\n\nWe are looking forward to your next visit.\n\nBest regards,\n{{companyName}}"),
				Tr: prisma.Str("Merhaba {{firstname}} {{salutation}},\n\nSize mutlu ve saglıklı yıllar dileriz.\n\nEn kısa zamanda görüşmek dileğiyle,\n{{companyName}}"),
			},
		},
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateEmailTemplate(prisma.EmailTemplateCreateInput{
		Name: "passwordReset",
		Title: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Passwort zurücksetzen"),
				En: prisma.Str("Password reset"),
				Tr: prisma.Str("Şifreyi Yenile"),
			},
		},
		Content: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("{{salutation}} {{lastName}},\n\nunter folgendem Link können Sie Ihr Passwort zurücksetzen:\n{{customerSubdomain}}/reset-password/{{passwordToken}}\n\nWir freuen uns auf Ihren nächsten Besuch.\n\nMit freundlichen Grüßen\nIhr {{companyName}}-Team"),
				En: prisma.Str("{{salutation}} {{firstName}},\n\nPlease click here to reset your password:\n{{customerSubdomain}}/reset-password/{{passwordToken}}\n\nBest regards,\n{{companyName}}"),
				Tr: prisma.Str("Merhaba {{firstname}} {{salutation}},\n\nŞifrenizi yenilemek için lütfen tıklayın:\n{{customerSubdomain}}/reset-password/{{passwordToken}}\n\nSaygılarımızla,\n{{companyName}}"),
			},
		},
	}).Exec(ctx)
	panicIf(err)

	_, err = c.CreateEmailTemplate(prisma.EmailTemplateCreateInput{
		Name: "activationLink",
		Title: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("Account aktivieren"),
				En: prisma.Str("Activate account"),
				Tr: prisma.Str("Account aktivieren"),
			},
		},
		Content: prisma.LocalizedStringCreateOneInput{
			Create: &prisma.LocalizedStringCreateInput{
				De: prisma.Str("{{salutation}} {{lastName}},\n\nunter folgendem Link können Sie Ihren Account aktivieren:\n{{customerSubdomain}}/activate-account/{{activateToken}}\n\nWir freuen uns auf Ihren nächsten Besuch.\n\nMit freundlichen Grüßen\nIhr {{companyName}}-Team"),
				En: prisma.Str("{{salutation}} {{firstName}},\n\nPlease click here to activate your account:\n{{customerSubdomain}}/activate-account/{{activateToken}}\n\nBest regards,\n{{companyName}}"),
				Tr: prisma.Str("Merhaba {{firstname}} {{salutation}},\n\nHesabınızı şuradan etkinleştirin:\n{{customerSubdomain}}/activate-account/{{activateToken}}\n\nSaygılarımızla,\n{{companyName}}"),
			},
		},
	}).Exec(ctx)
	panicIf(err)
}
