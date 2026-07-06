package mailer

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func SendNotification(toEmail string, bookingID int) {
	host := os.Getenv("MAILHOG_HOST") 
	if host == "" {
		host = "localhost"
	}
	
	addr := fmt.Sprintf("%s:1025", host)

	// Construction de l'e-mail
	from := "noreply@micro-trips.com"
	to := []string{toEmail}
	subject := "Subject: Confirmation de votre réservation\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("Bonjour,\n\nVotre réservation n°%d a bien été prise en compte ! Bon voyage.\n", bookingID)

	msg := []byte(subject + mime + body)

	// Envoi (MailHog n'a pas besoin d'authentification)
	err := smtp.SendMail(addr, nil, from, to, msg)
	if err != nil {
		log.Printf("❌ Erreur lors de l'envoi de l'e-mail à %s: %v", toEmail, err)
		return
	}
	
	log.Printf("E-mail de confirmation envoyé avec succès à %s (via MailHog)", toEmail)
}