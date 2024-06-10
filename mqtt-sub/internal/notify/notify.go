package notify

import (
	"fmt"
	"mqtt-sub/internal/common"
	"net/smtp"
	"strings"
)

type INotifier interface {
	Notify(string) error
}

type Notifier struct {
	email string
	pass  string
	host  string
	to    []string
}

func InitNotifier(c *common.Config) INotifier {
	return &Notifier{
		email: c.ClientFrom,
		pass:  c.SMTPPass,
		host:  c.SMTPServer,
		to:    c.ClientsTo,
	}
}

func (n *Notifier) Notify(alertmsg string) error {
	auth := smtp.PlainAuth("", n.email, n.pass, n.host)
	if sErr := smtp.SendMail(n.host, auth, n.email, n.to, []byte(alertmsg)); sErr != nil {
		return fmt.Errorf("failed to send msg via smtp; %v", sErr)
	}
	common.LogInfo(fmt.Sprintf("successfully notified recipients %v", strings.Join(n.to, ",")))
	return nil
}
