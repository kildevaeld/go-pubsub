package pubsub

import (
	"github.com/googollee/go-assert"
	"testing"
	"time"
)

func TestPubsubNoblock(t *testing.T) {
	c := make(chan interface{})
	ps := New()
	ps.Subscribe("nonblock", c)
	quit := make(chan int)
	go func() {
		select {
		case <-quit:
			return
		case <-time.After(time.Second):
			panic("blocked")
		}
	}()
	ps.Publish("nonblock", "msg")
	quit <- 1
}

func TestPubsub(t *testing.T) {
	quit := make(chan int)
	count := 0
	ps := New()

	abc := make(chan interface{}, 1)
	go func() {
		i := <-abc
		assert.Equal(t, i, "abc")
		quit <- 1
	}()
	count++
	ps.Subscribe("abc", abc)

	ab_any := make(chan interface{}, 2)
	go func() {
		i := <-ab_any
		assert.Equal(t, i, "abc")
		i = <-ab_any
		assert.Equal(t, i, "abd")
		quit <- 1
	}()
	count++
	ps.PSubscribe("ab*", ab_any)

	cde := make(chan interface{}, 1)
	go func() {
		i := <-cde
		assert.Equal(t, i, "cde")
		quit <- 1
	}()
	count++
	ps.Subscribe("cde", cde)

	cd_any1 := make(chan interface{}, 2)
	go func() {
		i := <-cd_any1
		assert.Equal(t, i, "cde")
		quit <- 1
	}()
	count++
	ps.PSubscribe("cd?", cd_any1)

	cd_any2 := make(chan interface{}, 2)
	go func() {
		i := <-cd_any2
		assert.Equal(t, i, "cde")
		quit <- 1
	}()
	count++
	ps.PSubscribe("cd[e]", cd_any2)

	ps.Publish("abc", "abc")
	ps.Publish("abd", "abd")
	ps.Publish("cde", "cde")

	for i := 0; i < count; i++ {
		<-quit
	}
}
