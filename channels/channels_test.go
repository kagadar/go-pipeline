package channels

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAwait(t *testing.T) {
	type want struct {
		i  int
		ok bool
	}
	for _, tt := range []struct {
		name string
		fn   func(chan int)
		want want
	}{
		{
			name: "yield",
			fn:   func(c chan int) { c <- 1 },
			want: want{i: 1, ok: true},
		},
		{
			name: "close",
			fn:   func(c chan int) { close(c) },
			want: want{i: 0, ok: false},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := make(chan int)
			go tt.fn(c)
			got, ok, err := Await(context.Background(), c)
			if err != nil {
				t.Fatalf("Await() unexpected error: %v", err)
			}
			if got != tt.want.i || ok != tt.want.ok {
				t.Errorf("Await() got {%d %t} want %v", got, ok, tt.want)
			}
		})
	}
}

func TestAwait_Err(t *testing.T) {
	for _, tt := range []struct {
		name string
		want error
		ctx  func() context.Context
	}{
		{
			name: "cancel",
			want: context.Canceled,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
		},
		{
			name: "deadline",
			want: context.DeadlineExceeded,
			ctx: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now())
				<-ctx.Done()
				cancel()
				return ctx
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := make(chan int)
			if _, _, err := Await(tt.ctx(), c); err != tt.want {
				t.Errorf("Await() unexpected error: got %v want %v", err, tt.want)
			}
		})
	}
}

func TestCollect(t *testing.T) {
	c := make(chan int)
	go func() {
		for i := 1; i < 4; i++ {
			c <- i
		}
		close(c)
	}()
	got, err := Collect(context.Background(), c)
	if err != nil {
		t.Fatalf("Collect() unexpected error: %v", err)
	}
	if diff := cmp.Diff(got, []int{1, 2, 3}); diff != "" {
		t.Errorf("Collect() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestCollect_Err(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan int)
	go func() {
		c <- 1
		cancel()
	}()
	got, err := Collect(ctx, c)
	if err != context.Canceled {
		t.Errorf("Collect() unexpected error: got %v want %v", err, context.Canceled)
	}
	if diff := cmp.Diff(got, []int{1}); diff != "" {
		t.Errorf("Collect() unexpected diff (-got +want):\n%s", diff)
	}
}
