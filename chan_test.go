package bypass_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/d4l3k/bypass"
)

func TestChanBypass(t *testing.T) {
	c := make(chan int, 10)
	c <- 1
	c <- 2
	c <- 3
	ch := bypass.WrapChan(c)
	out := ch.Elems().([]int)
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(out, want) {
		t.Fatal("ch.Elems() = %v; not %v", out, want)
	}
}

func TestChanRace(t *testing.T) {
	c := make(chan int, 10)
	cdone := make(chan struct{}, 2)
	wrapped := bypass.WrapChan(c)
	for i := 0; i < 2; i++ {
		go func() {
			for {
				select {
				case <-cdone:
					return
				case c <- 1:
				}
			}
		}()
	}
	go func() {
		for {
			_, open := <-c
			if !open {
				return
			}
		}
	}()
	go func() {
		for {
			wrapped.Elems()
		}
	}()
	time.Sleep(100 * time.Millisecond)
	cdone <- struct{}{}
	cdone <- struct{}{}
	close(c)

}
