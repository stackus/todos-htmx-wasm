package log

import (
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/stackus/errors"
)

// WebLogger is a http.Handler middleware that logs HTTP requests using zerolog.Logger.
func WebLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

			start := time.Now()

			defer func() {
				var err error
				var logFn func() *zerolog.Event
				var stack []byte
				p := recover()

				switch {
				case p != nil:
					stack = debug.Stack()
					logFn = logger.Error
					// ensure the status code reflects this panic
					if ww.Status() < 500 {
						ww.WriteHeader(http.StatusInternalServerError)
					}
					err = errors.ErrInternalServerError.Msgf("%s", p)
				case ww.Status() < 400:
					logFn = logger.Info
				case ww.Status() < 500:
					logFn = logger.Warn
				default:
					logFn = logger.Error
				}
				log := logFn()
				if err != nil {
					log = log.Err(err)
				}
				if stack != nil {
					log = log.Strs("Stack", strings.Split(string(stack), "\n"))
				}

				log = log.Str("RemoteAddr", request.RemoteAddr).
					Int("ContentLength", ww.BytesWritten()).
					Dur("ResponseTime", time.Since(start))

				log.Msgf("[%d] %s %s", ww.Status(), request.Method, request.RequestURI)
			}()

			next.ServeHTTP(ww, request)
		})
	}
}
