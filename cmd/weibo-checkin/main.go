package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xwjdsh/weibo-checkin"
)

func main() {
	h := weibo.New(os.Getenv("COOKIE"))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := h.SuperTopics(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("list success")

	success := true
	for _, t := range result.Data.List {
		if ss := strings.Split(t.Oid, ":"); len(ss) == 2 {
			select {
			case <-ctx.Done():
				panic(ctx.Err())
			case <-time.After(500 * time.Microsecond):
			}
			sctx, scancel := context.WithTimeout(ctx, time.Second)
			err := h.SuperTopicSignIn(sctx, ss[1])
			scancel()
			if err != nil {
				success = false
				fmt.Println(err.Error())
			}
		}
	}

	if !success {
		os.Exit(1)
	}
}
