package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func InitLogger(debug bool) (zerolog.Logger, error) {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	return logger, nil
}

type OpenObserverHook struct {
	Url      string
	Email    string
	Password string
}

type OpenObserverData struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

func (d OpenObserverData) String() string {
	bytes, _ := json.Marshal(d)
	return string(bytes)
}

func (o OpenObserverHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	go func() {
		data := OpenObserverData{
			Level: level.String(),
			Msg:   msg,
		}
		req, err := http.NewRequest("POST", o.Url, strings.NewReader(data.String()))
		if err != nil {
			fmt.Println("error sending log hook. error creating request: ", err)
		}
		req.SetBasicAuth(o.Email, o.Password)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("error sending log hook. error sending request: ", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			fmt.Println("error sending log hook. response error: ", string(body))
		}
	}()
}

type Logger struct {
	namespace string
	logger    *zerolog.Logger
}

func NewLogger(logger *zerolog.Logger, namespace string) *Logger {

	return &Logger{
		namespace: namespace,
		logger:    logger,
	}
}

func (l *Logger) Debugf(msg string, keyvals ...interface{}) {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "") // Add a default value for any missing keyvals
	}

	kvPairs := make([]interface{}, 0, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i]
		value := keyvals[i+1]
		kvPairs = append(kvPairs, key, value)
	}

	msgf := fmt.Sprintf(msg, kvPairs...)
	l.logger.Debug().
		Msgf("%s: %s", l.namespace, msgf)
}

func (l *Logger) Infof(msg string, keyvals ...interface{}) {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "") // Add a default value for any missing keyvals
	}

	kvPairs := make([]interface{}, 0, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i]
		value := keyvals[i+1]
		kvPairs = append(kvPairs, key, value)
	}

	msgf := fmt.Sprintf(msg, kvPairs...)
	l.logger.Info().
		Msgf("%s: %s", l.namespace, msgf)
}

func (l *Logger) Warnf(msg string, keyvals ...interface{}) {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "") // Add a default value for any missing keyvals
	}

	kvPairs := make([]interface{}, 0, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i]
		value := keyvals[i+1]
		kvPairs = append(kvPairs, key, value)
	}

	msgf := fmt.Sprintf(msg, kvPairs...)
	l.logger.Warn().
		Msgf("%s: %s", l.namespace, msgf)
}

func (l *Logger) Errorf(msg string, keyvals ...interface{}) {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "") // Add a default value for any missing keyvals
	}

	kvPairs := make([]interface{}, 0, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		key := keyvals[i]
		value := keyvals[i+1]
		kvPairs = append(kvPairs, key, value)
	}

	msgf := fmt.Sprintf(msg, kvPairs...)
	l.logger.Error().
		Msgf("%s: %s", l.namespace, msgf)
}
