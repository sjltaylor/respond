package respond

import "net/http"

type fakeResponseWriter struct {
  body string
  status int
}

func (rw *fakeResponseWriter) Header() http.Header {
  
  return nil
}

func (rw *fakeResponseWriter) Write(bytes []byte) (int, error) {
  rw.body = rw.body + string(bytes)
  return len(bytes), nil
}

func (rw *fakeResponseWriter) WriteHeader(sc int) {
  
  rw.status = sc
}

func newFakeResponseWriter() *fakeResponseWriter {
  return &fakeResponseWriter{}
}