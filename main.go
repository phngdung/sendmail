package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "net/smtp"
)

func sendMail(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, _ := template.ParseFiles("form.gtpl")
        t.Execute(w, nil)
    } else {
        r.ParseForm()

        var email string
        var title string
        var message string
        
        email= r.FormValue("email")
        title= r.FormValue("title")
        message= r.FormValue("message")

        send(email,title,message)
    }
}

func main() {
    http.HandleFunc("/sendMail", sendMail)
    err := http.ListenAndServe(":9090", nil) 
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func send(toMail string,subject string, body string) {

    // Sender data.
    from := "noreply@daugia.io"
    password := "abcD123$"
  
    // Receiver email address.
    to := []string{
      toMail,
    }
  
    // smtp server configuration.
    smtpHost := "smtp.yandex.com"
    smtpPort := "587"
  
    // Message.
    message := []byte("From: Demo<noreply@daugia.io>\r\n" +
          "Subject:"+ subject+"\r\n" +
          "\r\n" +
           body)
          
    // Authentication.
    auth := smtp.PlainAuth("", from, password, smtpHost)
    
    // Sending email.
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
    if err != nil {
      fmt.Println(err)
      return
    }
    fmt.Println("Email Sent Successfully!")
  }