package utility

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
)

func TestRabbitmq(t *testing.T) {
	ctx := gctx.New()
	ConsumeDeadLetter(ctx)
	Producer(ctx)

	time.Sleep(time.Hour)
}
