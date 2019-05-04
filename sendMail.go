package main

import (
	"net/smtp"
	"strings"
)

func mySendMail(toEmail string, content string) error {
	auth := smtp.PlainAuth("", "1335783803@qq.com", "peynpmapjgezhdgd", "smtp.qq.com")
	to := []string{toEmail}
	contentType := "Content-Type: text/html; charset=utf-8"
	nickName := "From:" + "biao" + "<1335783803@qq.com>"
	subject := "Subject:" + "用户账号激活"
	toUsers := "To: " + strings.Join(to, ",")
	body := "\n\n\n" + content
	msg := strings.Join([]string{nickName, subject, toUsers, contentType, body}, "\r\n")
	err := smtp.SendMail("smtp.qq.com:587", auth, "1335783803@qq.com", to, []byte(msg))
	return err
}
