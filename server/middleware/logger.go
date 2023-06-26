package middleware

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPLogger logs a gin HTTP request in JSON format. Allows to set the
// logger for testing purposes.
func HTTPLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		// Log using the params
		data := make(map[string]string)
		data["client_id"] = param.ClientIP
		data["method"] = param.Method
		data["status_code"] = strconv.Itoa(param.StatusCode)
		data["body_size"] = strconv.Itoa(param.BodySize)
		data["path"] = param.Path
		data["latency"] = param.Latency.String()
		data["error"] = param.ErrorMessage

		if c.Writer.Status() >= 500 {
			// response content only for errors (to get error messages in logs)
			data["response"] = blw.body.String()
			log.Printf("error: %s, err=%v, data=%+v", c.Errors.String(), c.Err(), data)
		} else {
			if c.Errors.String() != "" {
				log.Printf("error: %s, %+v", c.Errors.String(), data)
				log.Println(param.Path)
			} else {
				log.Printf("%s: %+v", param.Path, data)
			}
		}
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
