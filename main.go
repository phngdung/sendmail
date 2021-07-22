package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type MailServerInfo struct {
	MailFrom string
	MailPass string
	MailHost string
	MailPort string
}

var (
	mail *MailServerInfo
)

func GetMailOption() *MailServerInfo {
	if mail != nil {
		return mail
	}
	file, err := os.Open("../config.json")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&mail)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return mail
}

func sendMail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("template.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		var email string
		var title string
		var message string

		email = r.FormValue("email")
		title = r.FormValue("title")
		message = r.FormValue("content")

		result, status := send(email, title, message)
		if result {
			fmt.Fprintf(w, "To : %s\n", email)
			fmt.Fprintf(w, "Email Sent Successfully!")
		}
		if !result {
			fmt.Fprintf(w, "Error: %s!\n", status)

		}
	}

}

func send(toMail string, subject string, body string) (bool, string) {

	mailConfig := GetMailOption()

	from := mailConfig.MailFrom
	password := mailConfig.MailPass

	to := []string{
		toMail,
	}

	smtpHost := mailConfig.MailHost
	smtpPort := mailConfig.MailPort

	message := []byte("From: Demo<" + from + ">\r\n" +
		"Subject:" + subject + "\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return false, err.Error()
	}
	fmt.Println("Email Sent Successfully!")
	return true, "OK"
}

func main() {
	http.HandleFunc("/sendMail", sendMail)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
