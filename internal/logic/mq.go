package logic

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/rabbitmq/amqp091-go"
)

func Connect() {
	var (
		ctx = gctx.New()
		url = g.Cfg().MustGet(ctx, "rabbitmq.url").String()
	)
	conn, err := amqp091.Dial(url)
	fail(ctx, err)
	defer conn.Close()

	ch, err := conn.Channel()
	fail(ctx, err)
	defer ch.Close()

	ch.ExchangeDeclare()
	ch.QueueDeclare()
	ch.QueueBind()
	ch.ExchangeBind()
}

func fail(ctx context.Context, err error) {
	if err != nil {
		g.Log().Panicf(ctx, "Dial err: %s", err)
	}
}
