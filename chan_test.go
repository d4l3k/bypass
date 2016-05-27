package bypass_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/d4l3k/bypass"
)

func TestChanElemsInt(t *testing.T) {
	c := make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	<-c
	c <- 4
	ch := bypass.WrapChan(c)
	out := ch.Elems().([]int)
	want := []int{2, 3, 4}
	if !reflect.DeepEqual(out, want) {
		t.Fatal("ch.Elems() = %v; not %v", out, want)
	}
}

type test struct {
	a int
}

func TestChanElemsStruct(t *testing.T) {
	c := make(chan test, 3)
	c <- test{1}
	c <- test{2}
	c <- test{3}
	<-c
	c <- test{4}
	ch := bypass.WrapChan(c)
	out := ch.Elems().([]test)
	want := []test{{2}, {3}, {4}}
	if !reflect.DeepEqual(out, want) {
		t.Fatal("ch.Elems() = %v; not %v", out, want)
	}
}

func TestChanElemsStructPtr(t *testing.T) {
	c := make(chan *test, 3)
	c <- &test{1}
	c <- &test{2}
	c <- &test{3}
	<-c
	c <- &test{4}
	ch := bypass.WrapChan(c)
	out := ch.Elems().([]*test)
	want := []*test{{2}, {3}, {4}}
	if !reflect.DeepEqual(out, want) {
		t.Fatal("ch.Elems() = %v; not %v", out, want)
	}
}

func TestChanRace(t *testing.T) {
	c := make(chan int, 10)
	cdone := make(chan struct{})
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
