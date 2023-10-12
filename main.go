package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	ctx := context.WithValue(context.Background(), "userId", "umizu")
	username, err := fetchUsername(ctx, 550)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("time: %v\nresponse: %+v", time.Since(start), username)
}

func fetchUsername(ctx context.Context, waitTime int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(waitTime))
	defer cancel()

	userId := ctx.Value("userId").(string)

	type result struct {
		username string
		err      error
	}
	resultch := make(chan result, 1)

	go func(userId string) {
		res, err := thirdPartyHttpCall(userId)
		resultch <- result{
			username: res,
			err:      err}
	}(userId)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-resultch:
		return res.username, res.err
	}
}

func thirdPartyHttpCall(userId string) (string, error) {
	time.Sleep(time.Millisecond * 500)
	return fmt.Sprintf("some username from userId '%s'", userId), nil
}
