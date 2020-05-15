package remon

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

func Dial(t *testing.T) *ReMon {
	r := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	/*
		if _, err := r.Ping().Result(); err != nil {
			fmt.Printf("Error: %T", err)
			t.Fatal("failed to ping redis:", err)
		}
	*/

	m, err := mongo.NewClient(
		mongooptions.Client().ApplyURI("mongodb://127.0.0.1"))
	if err != nil {
		t.Fatal("failed to new mongo client:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := m.Connect(ctx); err != nil {
		t.Fatal("failed to connect mongo server:", err)
	}

	/*
		if err := ScriptLoad(r); err != nil {
			t.Fatal("failed to script load:", err)
		}
	*/

	return New(r, m)
}

func TestReMon(t *testing.T) {
	x := Dial(t)

	func() {
		fmt.Println("Set")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		fmt.Println(x.Set(ctx, "aa:bb:cc", "world"))
		fmt.Println(x.getDataFromRedis(ctx, "aa:bb:cc"))
	}()

	func() {
		fmt.Println("Get")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		fmt.Println(x.Get(ctx, "aa:bb:cc"))
		fmt.Println(x.getDataFromRedis(ctx, "aa:bb:cc"))
	}()

}

func TestMail(t *testing.T) {
	x := Dial(t)

	func() {
		fmt.Println("PushMail")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		fmt.Println(x.PushOneMail(ctx, "aa:bb:cc", "this is a mail"))
		fmt.Println(x.getDataFromRedis(ctx, "aa:bb:cc"))
	}()

	func() {
		fmt.Println("PullMail:")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		list, err := x.ListMail(ctx, "aa:bb:cc")
		fmt.Println(err)
		fmt.Println(x.PullMail(ctx, "aa:bb:cc", []string{list[0].Id, list[len(list)-1].Id}))
		fmt.Println(x.getDataFromRedis(ctx, "aa:bb:cc"))
	}()
}

func TestEvalVar(t *testing.T) {
	x := Dial(t)

	func() {
		fmt.Println("EvalVar")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		fmt.Println(x.EvalVarJSON(ctx, "aa:bb:dd", "something=1024"))
		fmt.Println(x.getDataFromRedis(ctx, "aa:bb:dd"))
	}()

	func() {
		fmt.Println("EvalVar")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		fmt.Println(x.EvalVarJSON(ctx, "aa:bb:dd", "something=10; return -1"))
		fmt.Println(x.getDataFromRedis(ctx, "aa:bb:dd"))
	}()

	func() {
		fmt.Println("EvalVar")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		fmt.Println(x.EvalVarJSON(ctx, "aa:bb:dd", "while true do end"))
		fmt.Println(x.getDataFromRedis(ctx, "aa:bb:dd"))
	}()
}
