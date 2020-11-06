package util

import (
	"net/http"
	"testing"
	"time"
)

func TestRetryBackoff(t *testing.T) {
	type args struct {
		n    int
		r    *http.Request
		resp *http.Response
	}
	tests := []struct {
		name            string
		args            args
		validdateOutput func(time.Duration) bool
	}{
		{
			name: "Do not retry a non 400 error",
			args: args{
				n:    0,
				r:    &http.Request{},
				resp: &http.Response{StatusCode: http.StatusUnauthorized},
			},
			validdateOutput: func(duration time.Duration) bool {
				return duration == -1
			},
		},
		{
			name: "Retry a 400 error when the first time",
			args: args{
				n:    0,
				r:    &http.Request{},
				resp: &http.Response{StatusCode: http.StatusBadRequest},
			},
			validdateOutput: func(duration time.Duration) bool {
				return duration > 0
			},
		},
		{
			name: "Retry a 400 error when when less than 6 times",
			args: args{
				n:    5,
				r:    &http.Request{},
				resp: &http.Response{StatusCode: http.StatusBadRequest},
			},
			validdateOutput: func(duration time.Duration) bool {
				return duration > 0
			},
		},
		{
			name: "Do not retry a 400 error after 6 tries",
			args: args{
				n:    6,
				r:    &http.Request{},
				resp: &http.Response{StatusCode: http.StatusBadRequest},
			},
			validdateOutput: func(duration time.Duration) bool {
				return duration == -1
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetryBackoff(tt.args.n, tt.args.r, tt.args.resp); !tt.validdateOutput(got) {
				t.Errorf("RetryBackoff() = %v which is not valid according to the validdateOutput()", got)
			}
		})
	}
}
