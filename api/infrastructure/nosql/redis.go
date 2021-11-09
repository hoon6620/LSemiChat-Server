package nosql

import (
	"app/api/llog"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
)

var client *redis.Client

type NosqlHandler interface {
	CreateAuth(userid string, rTime time.Time, token string) error
	DeleteAuth(givenUuid string) (int64, error)
}

func New() {
	dsn := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	if len(dsn) == 0 {
		dsn = "http://localhost:6379"
	}
	llog.Info(dsn)
	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		llog.Fatal(err)
	}
}

func CreateAuth(userid string, rTime time.Time, token string) error {
	now := time.Now()

	errAccess := client.Set(token, userid, rTime.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	return nil
}

func CheckValidToken(token string) error {
	_, err := client.Get(token).Result()
	if err != nil {
		return err
	}
	return nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
