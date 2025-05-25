package middleware

import (
	"log"
	"net/http"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

// getStatusColor returns appropriate color based on HTTP status code
func getStatusColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return ColorGreen
	case statusCode >= 300 && statusCode < 400:
		return ColorYellow
	case statusCode >= 400 && statusCode < 500:
		return ColorRed
	case statusCode >= 500:
		return ColorPurple
	default:
		return ColorWhite
	}
}

// getMethodColor returns appropriate color for HTTP method
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return ColorBlue
	case "POST":
		return ColorGreen
	case "PUT", "PATCH":
		return ColorYellow
	case "DELETE":
		return ColorRed
	default:
		return ColorCyan
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

// Write ensures status code is captured even if WriteHeader wasn't called explicitly
func (rw *responseWriter) Write(data []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(data)
}

// LogMiddleware logs HTTP requests and errors based on environment
func LogMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Wrap the response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				written:        false,
			}

			// Call the next handler
			next.ServeHTTP(wrapped, r)

			// Format timestamp
			timestamp := time.Now().Format("15:04:05")

			// Get colors for method and status
			methodColor := getMethodColor(r.Method)
			statusColor := getStatusColor(wrapped.statusCode)

			// Format the main log entry with colors
			log.Printf("%s[%s]%s %s%s%-7s%s %s%s%s %s%d",
				ColorDim, timestamp, ColorReset, // [timestamp]
				methodColor, ColorBold, r.Method, ColorReset, // METHOD
				ColorCyan, r.URL.Path, ColorReset, // /path
				statusColor, wrapped.statusCode, // status code
			)

			// Log errors with additional detail
			if wrapped.statusCode >= 400 {
				log.Printf("%s[ERROR]%s %s%s %s%s - %sStatus: %s%d%s - %sUA: %s%s",
					ColorRed+ColorBold, ColorReset, // [ERROR]
					methodColor+ColorBold, r.Method, ColorReset, // METHOD
					ColorCyan, r.URL.Path, ColorReset, // /path
					statusColor, wrapped.statusCode, ColorReset, // Status: code
					ColorDim, r.UserAgent(), ColorReset, // UA: user-agent
				)
			}

		})
	}
}
