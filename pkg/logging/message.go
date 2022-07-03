package logging

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
)

//MessageField hols key value fields
type MessageField struct {
	Key   string
	Value string
}

func (m MessageField) String() string {
	return fmt.Sprintf("%s=%s", m.Key, m.Value)
}

//Message holds fields from a log line
type Message struct {
	Timestamp *time.Time
	Level     string
	Message   string
	Fields    []MessageField
}

//LineParser exposes a Parse method to
//handle log entries
type LineParser interface {
	Parse(line []byte) (*Message, error)
}

//TODO: Make configurable
var timestampKeys = []string{"ts", "time", "timestamp", "date"}
var messageKeys = []string{"message", "msg"}

//JSONLineParser implements LineParser
type JSONLineParser struct {
}

//Parse will parse a JSON formatted log line
func (p JSONLineParser) Parse(line []byte) (*Message, error) {

	m := &Message{}
	var data map[string]interface{}

	err := json.Unmarshal(line, &data)
	if err != nil {
		return m, err
	}

	for _, key := range timestampKeys {
		//Handle Java/JavaScript unix timestamp e.g. 1656736672698
		if parseTimestampFloat(m, data, key) {
			break
		}
	}

	parseLevelString(m, data, "level")

	for _, key := range messageKeys {
		if parseMessage(m, data, key) {
			break
		}
	}

	if len(data) > 0 {
		for key, val := range data {
			value, err := json.Marshal(val)
			if err != nil {
				m.Fields = append(m.Fields, MessageField{Key: key, Value: fmt.Sprintf("%+v", val)})
			} else {
				m.Fields = append(m.Fields, MessageField{Key: key, Value: string(value)})
			}
		}
	}

	return m, nil
}

func parseMessage(m *Message, data map[string]interface{}, key string) bool {
	message, ok := data[key].(string)
	if ok {
		delete(data, key)
		m.Message = message
	}
	return ok
}

func parseLevelString(m *Message, data map[string]interface{}, key string) bool {
	level, ok := data[key].(string)
	if ok {
		delete(data, key)

		level = stripansi.Strip(level)
		level = strings.ToLower(level)
		level = strings.Trim(level, " ")
		m.Level = level
	}
	return ok
}

func parseTimestampFloat(m *Message, data map[string]interface{}, key string) bool {
	ts, ok := data[key].(float64)
	if ok {
		delete(data, key)
		sec, dec := math.Modf(ts)
		t := time.Unix(int64(sec), int64(dec*(1e9))).UTC()
		m.Timestamp = &t
	}
	return ok
}
