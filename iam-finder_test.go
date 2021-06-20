package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_getLastUsedOrCreated(t *testing.T) {
	type args struct {
		lastUsed *time.Time
		created  *time.Time
	}
	testTime := time.Now()
	tests := []struct {
		name            string
		args            args
		wantLastTouched *time.Time
		wantIsNotUsed   bool
	}{
		{
			"should return lastUsed time when lastUsed is not nil",
			struct {
				lastUsed *time.Time
				created  *time.Time
			}{lastUsed: &testTime, created: nil},
			&testTime,
			false,
		},
		{
			"should return created time when lastUsed is nil",
			struct {
				lastUsed *time.Time
				created  *time.Time
			}{lastUsed: nil, created: &testTime},
			&testTime,
			true,
		},
		{
			"should return nil when both time is nil",
			struct {
				lastUsed *time.Time
				created  *time.Time
			}{lastUsed: nil, created: nil},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLastTouched, gotIsNotUsed := getLastUsedOrCreated(tt.args.lastUsed, tt.args.created)
			if !reflect.DeepEqual(gotLastTouched, tt.wantLastTouched) {
				t.Errorf("getLastUsedOrCreated() gotLastTouched = %v, want %v", gotLastTouched, tt.wantLastTouched)
			}
			if gotIsNotUsed != tt.wantIsNotUsed {
				t.Errorf("getLastUsedOrCreated() gotIsNotUsed = %v, want %v", gotIsNotUsed, tt.wantIsNotUsed)
			}
		})
	}
}

func Test_isAccessKeyOld(t *testing.T) {
	type args struct {
		lastUsed    *time.Time
		expireHours int
	}
	now := time.Now()
	hoursAgo := []time.Time{
		now,
		now.Add(-1 * time.Hour),
		now.Add(-2 * time.Hour),
		now.Add(-3 * time.Hour),
		now.Add(-4 * time.Hour),
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "should return false when receives same our with time", args: args{
			lastUsed:    &hoursAgo[1],
			expireHours: 1,
		}, want: false},
		{name: "should return true when receives past time than expires", args: args{
			lastUsed:    &hoursAgo[2],
			expireHours: 1,
		}, want: true},
		{name: "should return false when receives least than expires", args: args{
			lastUsed:    &hoursAgo[2],
			expireHours: 3,
		}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAccessKeyOld(tt.args.lastUsed, tt.args.expireHours); got != tt.want {
				t.Errorf("isAccessKeyOld() = %v, want %v", got, tt.want)
			}
		})
	}
}
