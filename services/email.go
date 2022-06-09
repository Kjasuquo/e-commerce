package services

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"log"
	"time"
)

type Service struct{}

// SendMail  METHOD THAT  WILL BE USED TO SEND EMAILS TO USERS
func (s *Service) SendMail(subject, body, recipient, Private, Domain string) error {
	privateAPIKey := Private
	yourDomain := Domain

	mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Create a new template

	//err := mg.CreateTemplate(ctx, &mailgun.Template{
	//	Name: "i_template",
	//	Version: mailgun.TemplateVersion{
	//		Template: `'<div class="entry"> <h1>{{.title}}</h1> <div class="body"> {{.body}} </div> </div>'`,
	//		Engine:   mailgun.TemplateEngineGo,
	//		Tag:      "v1",
	//	},
	//})
	//if err != nil {
	//	return err
	//}

	// Create a new message with template
	m := mg.NewMessage("Oja Ecommerce <Oja@Decadev.gon>", subject, "")
	m.SetTemplate("i_template")

	// Add recipients
	err := m.AddRecipient(recipient)
	if err != nil {
		return err
	}

	// Add the variables recipient be used by the template
	err = m.AddVariable("title", subject)
	if err != nil {
		return err
	}
	err = m.AddVariable("body", body)
	if err != nil {
		return err
	}

	// Send the message with a 10 second timeout
	_, _, err = mg.Send(ctx, m)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}
