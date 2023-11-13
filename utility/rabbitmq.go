package utility

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/rabbitmq/amqp091-go"
)

const (
	dlxName = "my.dlx"
	dlkName = "my.dlk"

	normalEx = "nex"
	normalQ  = "nq"
	normalRK = "nrk"
)

var (
	RabbitmqConn    *amqp091.Connection
	RabbitmqChannel *amqp091.Channel
)

func init() {
	var (
		ctx = gctx.New()
		url = g.Cfg().MustGet(ctx, "rabbitmq.url").String()
		err error
	)
	RabbitmqConn, err = amqp091.Dial(url)
	if err != nil {
		g.Log().Panicf(ctx, "Dial err: %s", err)
	}
	RabbitmqChannel, err = RabbitmqConn.Channel()
	if err != nil {
		g.Log().Panicf(ctx, "Channel err: %s", err)
	}
}

func ConsumeDeadLetter(ctx context.Context) {
	err := RabbitmqChannel.ExchangeDeclare(dlxName, "direct", true, false,
		false, false, nil)
	if err != nil {
		g.Log().Panicf(ctx, "ExchangeDeclare err: %s", err)
	}

	q, err := RabbitmqChannel.QueueDeclare("my.dlq", true, false, false,
		false, nil)
	if err != nil {
		g.Log().Panicf(ctx, "QueueDeclare err: %s", err)
	}

	err = RabbitmqChannel.QueueBind(q.Name, dlkName, dlxName, false, nil)
	if err != nil {
		g.Log().Panicf(ctx, "QueueBind err: %s", err)
	}

	deliveries, err := RabbitmqChannel.Consume(q.Name, "my.dlk.c", true,
		false, false, false, nil)
	if err != nil {
		g.Log().Panicf(ctx, "Consume err: %s", err)
	}
	_ = grpool.Add(ctx, func(ctx context.Context) {
		for delivery := range deliveries {
			g.Log().Infof(ctx, "%s: %s", delivery.MessageId, delivery.Body)
		}
	})
}

func Producer(ctx context.Context) {
	err := RabbitmqChannel.ExchangeDeclare(normalEx, "direct", false, true,
		false, false, nil)
	if err != nil {
		g.Log().Panicf(ctx, "ExchangeDeclare err: %s", err)
	}

	q, err := RabbitmqChannel.QueueDeclare(normalQ, false, true, false,
		false, amqp091.Table{
			"x-dead-letter-exchange":    dlxName,
			"x-message-ttl":             10000, // 10s
			"x-dead-letter-routing-key": dlkName,
		})
	if err != nil {
		g.Log().Panicf(ctx, "QueueDeclare err: %s", err)
	}

	err = RabbitmqChannel.QueueBind(q.Name, normalRK, normalEx, false, nil)
	if err != nil {
		g.Log().Panicf(ctx, "QueueBind err: %s", err)
	}

	_ = grpool.Add(ctx, func(ctx context.Context) {
		var i int
		for {
			i++
			err := RabbitmqChannel.PublishWithContext(ctx, normalEx, normalRK, false,
				false, amqp091.Publishing{
					Body: []byte(fmt.Sprintf("%d", i)),
				})
			if err != nil {
				g.Log().Errorf(ctx, "PublishWithContext err: %s", err)
				break
			}
			time.Sleep(time.Second)
		}
	})
}
