package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type Something struct {
	Added  time.Time `json:"added"`
	Number int       `json:"number"`
}

func (s *Something) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Something) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return nil
}

func getSentinelAddrs() []string {
	x := os.Getenv("SENTINELS")
	return strings.Split(x, ";")
}

func main() {
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		Password:      "str0ng_passw0rd",
		MasterName:    "mymaster",
		DB:            0,
		SentinelAddrs: getSentinelAddrs(),
	})

	ctx := context.Background()
	go setkey(rdb, ctx)
	go getkey(rdb, ctx)

	forever := make(chan bool)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func getkey(rdb *redis.Client, ctx context.Context) {
	for {
		time.Sleep(time.Millisecond * 200)

		res := rdb.Get(ctx, "mykey")
		if err := res.Err(); err != nil {
			fmt.Printf("fail to get value for key, %s\n", err)
			continue
		}

		s := Something{}
		if err := res.Scan(&s); err != nil {
			fmt.Printf("fail to unmarshal %s\n", err)
			continue
		}
		fmt.Printf("%+v\n", s)
	}
}

func setkey(rdb *redis.Client, ctx context.Context) {
	x := 3

	for {
		time.Sleep(time.Millisecond * 200)
		s := &Something{
			Added:  time.Now(),
			Number: x,
		}
		expire := time.Duration(0)
		key := "mykey"

		if err := rdb.Set(ctx, key, s, expire).Err(); err != nil {
			fmt.Printf("fail to set key:%s %s\n", key, err)
		}
		x++
	}
}
