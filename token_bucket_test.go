package ratelimiter

import (
	"context"
	"github.com/althk/ratelimiter/store"
	"reflect"
	"testing"
	"time"
)

func TestTokenBucket_Allow(t *testing.T) {
	type fields struct {
		refillRate float32
		cap        int32
	}
	tests := []struct {
		name   string
		fields fields
		want   []bool
	}{
		{name: "test1", fields: fields{refillRate: 1.0, cap: 3}, want: []bool{true, true, true, false, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewTBL(tt.fields.refillRate, tt.fields.cap, store.NewTesting())
			got := make([]bool, 0)
			time.Sleep(time.Second * 4) // allow 4 refills
			for i := range tt.want {
				got = append(got, tb.Allow(context.Background(), "id1"))
				if i == len(tt.want)-2 {
					time.Sleep(time.Second) // allow 1 refill
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Allow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenBucket_AllowN(t *testing.T) {
	type fields struct {
		refillRate float32
		cap        int32
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
		{name: "test1", fields: fields{refillRate: 2.0, cap: 4}, want: []bool{true, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewTBL(tt.fields.refillRate, tt.fields.cap, store.NewTesting())
			got := make([]bool, 0)
			time.Sleep(time.Second * 1) // allow 1 refills
			for i := range tt.want {
				got = append(got, tb.AllowN(context.Background(), "id1", time.Now(), 2))
				if i == len(tt.want)-2 {
					time.Sleep(time.Second) // allow 1 refill
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllowN() = %v, want %v", got, tt.want)
			}
		})
	}
}
