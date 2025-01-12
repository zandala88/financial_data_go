package server

import (
	"crypto/tls"
	"errors"
	"financia/config"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"net"
)

func SendEmail(target, subject, content string) {
	server := config.Configs.Email.Server
	port := config.Configs.Email.Port
	usr := config.Configs.Email.User
	pwd := config.Configs.Email.Password

	zap.S().Info("SMTP服务器:", server, "端口:", port, "用户名:", usr, "密码:", pwd)
	// 创建一个新的邮件发送器
	d := gomail.NewDialer(server, port, usr, pwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// 尝试连接并获取 SendCloser
	sendCloser, err := d.Dial()
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Println("连接超时:", err)
		}
		return
	}
	defer sendCloser.Close() // 确保在函数结束时关闭连接

	fmt.Println("成功连接到SMTP服务器")

	// 创建邮件
	m := gomail.NewMessage()
	// 设置邮件头
	m.SetHeader("From", usr)
	// 设置收件人
	m.SetHeader("To", target)
	// 设置抄送人为自己
	m.SetAddressHeader("Cc", usr, "admin")
	// 设置邮件主题
	m.SetHeader("Subject", subject)
	// 设置邮件内容，支持html格式
	m.SetBody("text/html", content)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("邮件发送错误:", err)
		return
	}
	fmt.Println("邮件发送成功！")
}
