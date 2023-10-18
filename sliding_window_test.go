package ratelimiter

import (
	"context"
	"github.com/althk/ratelimiter/store"
	"reflect"
	"testing"
	"time"
)

func TestSlidingWindow_Allow(t *testing.T) {
	type fields struct {
		d time.Duration
		l int64
	}
	tests := []struct {
		name   string
		fields fields
		want   []bool
	}{
		{name: "test1", fields: fields{d: 2 * time.Second, l: 2}, want: []bool{true, true, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := NewSWL(tt.fields.d, tt.fields.l, store.NewTesting())
			got := make([]bool, 0)
			for i := range tt.want {
				got = append(got, sw.Allow(context.Background(), "id1"))
				if i == len(tt.want)-2 {
					time.Sleep(3 * time.Second)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Allow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlidingWindow_AllowN(t *testing.T) {
	type fields struct {
		d time.Duration
		l int64
	}
	type args struct {
		tokens int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []bool
	}{
		{name: "test1", fields: fields{d: 2 * time.Second, l: 2}, want: []bool{true, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			swl := NewSWL(tt.fields.d, tt.fields.l, store.NewTesting())
			got := make([]bool, 0)
			for i := range tt.want {
				got = append(got, swl.AllowN(context.Background(), "id1", time.Now(), 2))
				if i == len(tt.want)-2 {
					time.Sleep(4 * time.Second) // allow 1 refill
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllowN() = %v, want %v", got, tt.want)
			}
		})
	}
}
