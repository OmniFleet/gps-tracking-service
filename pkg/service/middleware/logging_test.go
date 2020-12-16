package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

type LoggingFormat struct {
	Level     string    `json:"level"`
	Time      time.Time `json:"time"`
	Duration  float64   `json:"duration"`
	BytesIn   float64   `json:"bytesIn"`
	BytesOut  float64   `json:"bytesOUt"`
	Method    string    `json:"method"`
	Proto     string    `json:"proto"`
	RemoteIP  string    `json:"remoteIP"`
	RequestID string    `json:"requestId"`
	Status    int       `json:"status"`
	URL       string    `json:"url"`
	UserAgent string    `json:"userAgent"`
	Message   string    `json:"message"`
}

func newLogger(buffer *bytes.Buffer) *zerolog.Logger {
	log := zerolog.New(buffer)
	return &log
}

func UnmarshalBytesToStruct(b []byte) (*LoggingFormat, error) {
	var data LoggingFormat
	err := json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func TestFieldsLogged(t *testing.T) {
	buf := &bytes.Buffer{}
	log := newLogger(buf)
	handler := ZeroLog(log)

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/some/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to the SecurityHeaders
	// middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	handler(next).ServeHTTP(rr, r)

	d, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}

	data, err := UnmarshalBytesToStruct(d)
	if err != nil {
		t.Errorf("log could not be unmarshaled: %s", err.Error())
	}

	t.Run("LogLevel", func(t *testing.T) {
		if data.Level != "info" {
			t.Errorf("want: %q; got: %q", "info", data.Level)
		}
	})
	t.Run("RequestTime", func(t *testing.T) {
		now := time.Now()
		logTime := data.Time.Unix()
		if logTime > now.Unix() || logTime < (now.Unix()-10) {
			t.Errorf("got %q; want %q; now %q; logtime %q", logTime, now.Unix(), now, data.Time)
		}
	})
	t.Run("Duration", func(t *testing.T) {
		if !(data.Duration > 0) {
			t.Errorf("no duration calculated: %f", data.Duration)
		}
	})
	t.Run("BytesOut", func(t *testing.T) {
		if !(data.BytesOut > 0) {
			t.Errorf("bytes out not recorded: %f", data.BytesOut)
		}
	})
	t.Run("Method", func(t *testing.T) {
		if data.Method != http.MethodGet {
			t.Errorf("got %q; want %q", data.Method, http.MethodGet)
		}
	})
	t.Run("Status", func(t *testing.T) {
		if data.Status != http.StatusOK {
			t.Errorf("got %q; want %q", data.Status, http.StatusOK)
		}
	})
	t.Run("URL", func(t *testing.T) {
		if data.URL != "/some/path" {
			t.Errorf("got %q; want %q", data.URL, "/some/path")
		}
	})
}

func TestLoggingMiddlewareCallsNextHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to the SecurityHeaders
	// middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	buf := &bytes.Buffer{}
	handler := ZeroLog(newLogger(buf))
	handler(next).ServeHTTP(rr, r)

	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}
	defer rs.Body.Close()

	// check if the next handler is called
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("next handler returns OK; got %s", string(body))
	}
}
