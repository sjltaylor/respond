package testHelpers

import "net/http"

type FakeResponseWriter struct {
  Body string
  Status int
}

func (rw *FakeResponseWriter) Header() http.Header {
  
  return nil
}

func (rw *FakeResponseWriter) Write(bytes []byte) (int, error) {
  rw.Body = rw.Body + string(bytes)
  return len(bytes), nil
}

func (rw *FakeResponseWriter) WriteHeader(sc int) {
  
  rw.Status = sc
}

func NewFakeResponseWriter() *FakeResponseWriter {
  return &FakeResponseWriter{}
}