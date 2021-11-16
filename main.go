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
	sentineladdrs := os.Getenv("SENTINELS")
	return strings.Split(sentineladdrs, ";")
}

func main() {

	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		Password:     "str0ng_passw0rd",
		DialTimeout:  time.Second * 1,
		ReadTimeout:  time.Second * 1,
		WriteTimeout: time.Second * 1,
		DB:           0,

		MasterName:    "mymaster",
		SentinelAddrs: getSentinelAddrs(),
	})
	ctx := context.Background()

	go func() {
		x := 3
		for {
			x++
			s := &Something{
				Added:  time.Now(),
				Number: x,
			}

			if err := rdb.Set(ctx, "mykey", s, time.Duration(0)).Err(); err != nil {
				fmt.Printf("fail to set key: %s\n", err)
			}

			time.Sleep(time.Millisecond * 200)
		}
	}()

	go func() {
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
	}()

	forever := make(chan bool)

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
