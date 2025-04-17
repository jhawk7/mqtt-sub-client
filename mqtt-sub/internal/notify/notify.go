package notify

import (
	"context"
	"fmt"
	"mqtt-sub/internal/common"
	"mqtt-sub/internal/dataparser"
	"net/smtp"
	"strings"

	"github.com/redis/go-redis/v9"
)

type INotifier interface {
	Notify(dataparser.IDataParser, string) error
}

type Notifier struct {
	email  string
	pass   string
	server string
	port   string
	to     []string
	rdb    *redis.Client
}

func InitNotifier(c *common.Config) INotifier {
	rdb := initRDB(c)

	return &Notifier{
		email:  c.ClientFrom,
		pass:   c.SMTPPass,
		server: c.SMTPServer,
		port:   c.SMTPPort,
		to:     c.ClientsTo,
		rdb:    rdb,
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

	if len(alertmsg) == 0 {
		common.LogInfo("notifier message is empty; dropping")
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
	key := parser.GetMeterName()
	ttl, tErr := n.rdb.TTL(context.Background(), key).Result()
	if tErr != nil {
		common.LogError(fmt.Errorf("failed to get ttl for key %v; %v", key, tErr), false)
		return false
	}

	// key is expired or doesn't exist per TTL documentation
	if ttl == -2 {
		//create and inc
		if _, incErr := n.rdb.Incr(context.Background(), key).Result(); incErr != nil {
			common.LogError(fmt.Errorf("failed to increment value for parser key %v; %v", key, incErr), false)
			return false
		}

		//set expiration
		if _, expErr := n.rdb.Expire(context.Background(), key, parser.NotificationRate()).Result(); expErr != nil {
			common.LogError(fmt.Errorf("failed to set expiration for key %v; %v", key, expErr), false)
			return false
		}

		return true
	}

	common.LogInfo(fmt.Sprintf("key %v has not yet expired and is therefore rate limited", key))
	return false
}
