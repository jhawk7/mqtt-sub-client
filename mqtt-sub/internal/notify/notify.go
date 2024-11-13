package notify

import (
	"context"
	"fmt"
	"mqtt-sub/internal/common"
	"mqtt-sub/internal/dataparser"
	"net/smtp"
	"strings"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type INotifier interface {
	Notify(dataparser.IDataParser, string) error
}

type Notifier struct {
	email   string
	pass    string
	server  string
	port    string
	to      []string
	rdb     *redis.Client
	limiter *redis_rate.Limiter
}

func InitNotifier(c *common.Config) INotifier {
	rdb := initRDB(c)
	limiter := redis_rate.NewLimiter(rdb)

	return &Notifier{
		email:   c.ClientFrom,
		pass:    c.SMTPPass,
		server:  c.SMTPServer,
		port:    c.SMTPPort,
		to:      c.ClientsTo,
		rdb:     rdb,
		limiter: limiter,
	}
}

func initRDB(config *common.Config) *redis.Client {
	opts := redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPass,
		DB:       0,
	}

	svc := redis.NewClient(&opts)

	if _, pingErr := svc.Ping(context.Background()).Result(); pingErr != nil {
		err := fmt.Errorf("failed to connect to redis instance; %v", pingErr)
		common.LogError(err, true)
	}

	common.LogInfo("established connection to redis instance")
	return svc
}

func (n *Notifier) Notify(parser dataparser.IDataParser, alertmsg string) error {
	//check rate limiter
	if allowed := n.checkRateLimit(parser); !allowed {
		common.LogInfo(fmt.Sprintf("topic notifications for %v exceeded limit; dropping", parser.GetMeterName()))
		return nil
	}

	auth := smtp.PlainAuth("", n.email, n.pass, n.server)

	host := fmt.Sprintf("%v:%v", n.server, n.port)
	if sErr := smtp.SendMail(host, auth, n.email, n.to, []byte(alertmsg)); sErr != nil {
		return fmt.Errorf("failed to send msg via smtp; %v", sErr)
	}
	common.LogInfo(fmt.Sprintf("successfully notified recipients %v", strings.Join(n.to, ",")))
	return nil
}

func (n *Notifier) checkRateLimit(parser dataparser.IDataParser) bool {
	limit := redis_rate.Limit{
		Rate:   1,
		Burst:  1,
		Period: parser.NotificationRate(),
	}

	res, limitErr := n.limiter.Allow(context.Background(), parser.GetMeterName(), limit)
	if limitErr != nil {
		err := fmt.Errorf("failed to retrieve notify limit for event %v; %v", parser.GetMeterName(), limitErr)
		common.LogError(err, false)
		return true
	}

	return res.Allowed > 0
}
