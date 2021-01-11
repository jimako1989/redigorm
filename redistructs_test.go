package redistructs

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"testing"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/ory/dockertest/v3"
	"github.com/tk42/redistructs/types"
)

var (
	dockerRes *dockertest.Resource
	redisPool *redigo.Pool
	TTL       = int64(2)
)

type Post struct {
	ID        uint64 `redis:"id"`
	UserID    uint64 `redis:"user_id"`
	Title     string `redis:"title"`
	Body      string `redis:"body"`
	CreatedAt int64  `redis:"created_at"`
}

func (p *Post) StoreType() types.StoreType {
	return types.Serialized
}

func (p *Post) PrimaryKey() string {
	return fmt.Sprint(p.ID)
}

func (p *Post) KeyDelimiter() string {
	return "/"
}

func (p *Post) ScoreMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     p.ID,
		"recent": p.CreatedAt,
	}
}

func (p *Post) Expire() interface{} {
	return time.Duration(time.Second)
}

func (p *Post) DatabaseIdx() int {
	return 0
}

// Serialized implements the types.Model interface
func (p *Post) Serialized() []byte {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(&p)
	if err != nil {
		panic("Failed to Serialized")
	}
	return buf.Bytes()
}

// Deserialized implements the types.Model interface
func (p *Post) Deserialized(b []byte) {
	err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(&p)
	if err != nil {
		panic("Failed to Deserialized. " + err.Error())
	}
}

func TestMain(m *testing.M) {
	setup()

	defer dockerRes.Close()

	m.Run()
}

func setup() {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("could not connect to docker, " + err.Error())
	}
	dockerRes, err = dockerPool.Run("redis", "5.0", nil)
	if err != nil {
		log.Fatal("could not start resource, " + err.Error())
	}

	redisPool = &redigo.Pool{
		MaxIdle:     5,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(fmt.Sprintf("redis://localhost:%s", dockerRes.GetPort("6379/tcp")))
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func TestPut(t *testing.T) {
	conn := redisPool.Get()
	defer conn.Close()

	now := time.Now()
	postStore := New(redisPool, *types.CreateConfig(), &Post{})

	err := postStore.Put(context.TODO(), []*Post{
		{
			ID:        1,
			UserID:    1,
			Title:     "post 1",
			Body:      "This is a post 1",
			CreatedAt: now.UnixNano(),
		},
		{
			ID:        2,
			UserID:    2,
			Title:     "post 2",
			Body:      "This is a post 2",
			CreatedAt: now.Add(-24 * 60 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        3,
			UserID:    1,
			Title:     "post 3",
			Body:      "This is a post 3",
			CreatedAt: now.Add(24 * 60 * 60 * time.Second).UnixNano(),
		},
		{
			ID:        4,
			UserID:    1,
			Title:     "post 4",
			Body:      "This is a post 4",
			CreatedAt: now.Add(-24 * 60 * 60 * time.Second).UnixNano(),
		},
	})
	if err != nil {
		panic(err)
	}

	// TODO
	keys, _ := redigo.Strings(conn.Do("KEYS"))
	fmt.Printf("%v", keys)
}
