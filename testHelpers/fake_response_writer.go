package testHelpers

import "net/http"

type FakeResponseWriter struct {
  Body string
  Status int
  header http.Header
}

func (rw *FakeResponseWriter) Header() http.Header {
  
  return rw.header
}

func (rw *FakeResponseWriter) Write(bytes []byte) (int, error) {
  rw.Body = rw.Body + string(bytes)
  return len(bytes), nil
}

func (rw *FakeResponseWriter) WriteHeader(sc int) {
  
  rw.Status = sc
}

func NewFakeResponseWriter() *FakeResponseWriter {
  return &FakeResponseWriter{
  	header: make(http.Header),
  }
}