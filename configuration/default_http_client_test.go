package configuration

import (
	"reflect"
	"testing"
	"time"
)

func TestDefaultClient_ShouldReturnDefaultTransport(t *testing.T) {
	want := "*http.Transport"
	subject := NewDefaultHttpClient(4*time.Second, false)
	got := reflect.TypeOf(subject.Transport).String()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}

func TestDefaultClient_ShouldReturnLoggingRoundTripper(t *testing.T) {
	want := "*configuration.loggingRoundTripper"
	subject := NewDefaultHttpClient(4*time.Second, true)
	got := reflect.TypeOf(subject.Transport).String()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}

func TestDefaultClient_ShouldReturnCustomTimeout(t *testing.T) {
	want := 2 * time.Second
	subject := NewDefaultHttpClient(want, false)
	got := subject.Timeout

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}
