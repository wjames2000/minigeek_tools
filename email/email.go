// @Time : 2020/5/14 4:36 下午
// @Author : minigeek
package email

import (
	"crypto/tls"
	"math/rand"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/wjames2000/minigeek_tools/log"
	"github.com/wjames2000/minigeek_tools/uuid"
)

// 随机数种子
var Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// 获取随机数
func getRandNumber(n int) int {
	return Rnd.Intn(n)
}

var eLog *log.Log

func init() {
	eLog = log.Init("20060102.email")
}

var (
	emailListTLS = [][]string{
		{"a@hxyfavorite.top", "Admin1234", "smtp.exmail.qq.com", "smtp.exmail.qq.com:465"},
		{"b@hxyfavorite.top", "Admin12345", "smtp.exmail.qq.com", "smtp.exmail.qq.com:465"},
		{"c@hxyfavorite.top", "Admin123", "smtp.exmail.qq.com", "smtp.exmail.qq.com:465"},
		{"d@hxyfavorite.top", "Admin1234", "smtp.exmail.qq.com", "smtp.exmail.qq.com:465"},
	}
	emailList = [][]string{
		{"a@hxyfavorite.top", "Admin1234", "smtp.exmail.qq.com", "smtp.exmail.qq.com:25"},
		{"b@hxyfavorite.top", "Admin12345", "smtp.exmail.qq.com", "smtp.exmail.qq.com:25"},
		{"c@hxyfavorite.top", "Admin123", "smtp.exmail.qq.com", "smtp.exmail.qq.com:25"},
		{"d@hxyfavorite.top", "Admin1234", "smtp.exmail.qq.com", "smtp.exmail.qq.com:25"},
	}
)

func SendEmail(title, content, touser, projectName string, i ...int) error {
	var index int
	if len(i) > 0 {
		index = i[0]
	}

	newTitle := "[" + projectName + "]" + title + "--" + uuid.NewUUID().HexToUpper()
	auth := smtp.PlainAuth("", emailList[index][0], emailList[index][1], emailList[index][2])
	host := emailList[index][3]
	touser = strings.Replace(touser, " ", "", -1)
	to := strings.Split(touser, ",") // 收件人  ;号隔开
	content_type := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + touser + "\r\nFrom: " + emailList[index][0] + ">\r\nSubject:" + newTitle + "\r\n" + content_type + "\r\n\r\n" + content)
	err := smtp.SendMail(host, auth, emailList[index][0], to, msg)
	if err != nil {
		if index < len(emailList)-1 {
			return SendEmail(title, content, touser, projectName, index+1)
		}
		return err
	}
	eLog.Println("email.SendEmail-->projectName:" + projectName + ",title:" + title + ",content:" + content + ",发件邮箱:" + emailList[index][0] + ",收件邮箱:" + touser)
	return nil
}
func SendEmailTLS(title, content, touser, projectName string, i ...int) (err error) {
	index := getRandNumber(len(emailListTLS))
	if len(i) > 0 {
		index = i[0]
	}
	newTitle := "[" + projectName + "]" + title + "--" + uuid.NewUUID().HexToUpper()
	auth := smtp.PlainAuth("", emailListTLS[index][0], emailListTLS[index][1], emailListTLS[index][2])
	addr := emailListTLS[index][3]
	user := emailListTLS[index][0]

	content_type := "Content-Type: text/html; charset=UTF-8"
	touser = strings.Replace(touser, " ", "", -1)
	toUser := strings.Split(touser, ",") // 收件人逗号隔开
	msg := []byte("To: " + touser + "\r\nFrom: " + emailListTLS[index][0] + "\r\nSubject:" + newTitle + "\r\n" + content_type + "\r\n\r\n" + content)

	host, _, _ := net.SplitHostPort(addr)
	tlsconfig := &tls.Config{InsecureSkipVerify: true, ServerName: host}
	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return
	}
	defer client.Close()

	if ok, _ := client.Extension("AUTH"); ok {
		if err = client.Auth(auth); err != nil {
			return
		}
	}
	if err = client.Mail(user); err != nil {
		return
	}
	for _, addr := range toUser {
		if err = client.Rcpt(addr); err != nil {
			eLog.Println("email.SendEmailTLS-->收件邮箱:" + addr + ",err:" + err.Error())
			continue
		}
	}
	w, err := client.Data()
	if err != nil {
		return
	}
	_, err = w.Write(msg)
	if err != nil {
		return
	}
	err = w.Close()
	if err != nil {
		return
	}
	eLog.Println("email.SendEmailTLS-->projectName:" + projectName + ",title:" + title + ",content:" + content + ",发件邮箱:" + emailListTLS[index][0] + ",收件邮箱:" + touser)
	return client.Quit()
}
