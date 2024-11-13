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
	email  string
	pass   string
	server string
	port   string
	to     []string
}

func InitNotifier(c *common.Config) INotifier {
	return &Notifier{
		email:  c.ClientFrom,
		pass:   c.SMTPPass,
		server: c.SMTPServer,
		port:   c.SMTPPort,
		to:     c.ClientsTo,
	}
}

//use redis to store topic as id and # of tokens/ message that can be sent in a given time period (also store)
// topicid, tokens, last ts, refresh rate

func (n *Notifier) Notify(alertmsg string) error {
	auth := smtp.PlainAuth("", n.email, n.pass, n.server)

	host := fmt.Sprintf("%v:%v", n.server, n.port)
	if sErr := smtp.SendMail(host, auth, n.email, n.to, []byte(alertmsg)); sErr != nil {
		return fmt.Errorf("failed to send msg via smtp; %v", sErr)
	}
	common.LogInfo(fmt.Sprintf("successfully notified recipients %v", strings.Join(n.to, ",")))
	return nil
}
