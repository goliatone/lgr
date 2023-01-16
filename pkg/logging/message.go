package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
)

var numericLevelMap = map[int]string{
	10: "trace",
	20: "debug",
	30: "info",
	40: "warn",
	50: "error",
	60: "fatal",
	70: "panic",
}

var levelNumberMap = map[string]int{
	"trace":   10,
	"debug":   20,
	"info":    30,
	"warn":    40,
	"warning": 40,
	"error":   50,
	"err":     50,
	"fatal":   60,
	"panic":   70,
}

type FieldFormatter = func(format string, a ...any) string

// MessageData encodes the log payload
type MessageData = map[string]interface{}

// MessageField hols key value fields
type MessageField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	tpl   string
}

func (m *MessageField) WithTemplate(s string) {
	m.tpl = s
}

func (m MessageField) String() string {
	return fmt.Sprintf(m.tpl, m.Key, m.Value)
}

type sortableFields []*MessageField

func (s sortableFields) Len() int           { return len(s) }
func (s sortableFields) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortableFields) Less(i, j int) bool { return s[i].Key < s[j].Key }

// Message holds fields from a log line
type Message struct {
	Timestamp *time.Time      `json:"timestamp"`
	Level     string          `json:"label"`
	Incidence int             `json:"level"`
	Message   string          `json:"message"`
	Caller    string          `json:"caller"`
	Fields    []*MessageField `json:"fields"`
	Line      int             `json:"line"`
}

func (m Message) WithFieldTemplate(t string) {
	for _, field := range m.Fields {
		field.WithTemplate(t)
	}
}

// HasFields return true if there are extra fields
func (m Message) HasFields() bool {
	return len(m.Fields) > 0
}

// GetTimestampOrNow return given timestamp or now
func (m Message) GetTimestampOrNow() *time.Time {
	if m.Timestamp == nil {
		t := time.Now()
		return &t
	}
	return m.Timestamp
}

// LineParser exposes a Parse method to
// handle log entries
type LineParser interface {
	Parse(line []byte) (*Message, error)
}

// TODO: Make configurable
var timestampKeys = []string{"ts", "time", "timestamp", "date", "@timestamp"}
var messageKeys = []string{"message", "msg"}
var levelKeys = []string{"level", "log.level", "severity"}
var callerKeys = []string{"caller", "logger"}

// JSONLineParser implements LineParser
type JSONLineParser struct {
}

func dropNonJSON(b []byte) []byte {
	s := bytes.IndexRune(b, '{')
	if s == -1 {
		return b
	}
	return b[s:]
}

// Parse will parse a JSON formatted log line
func (p JSONLineParser) Parse(line []byte) (*Message, error) {

	line = dropNonJSON(line)

	m := &Message{}
	var data MessageData

	err := json.Unmarshal(line, &data)
	if err != nil {
		return m, err
	}

	for _, key := range timestampKeys {
		//Handle Java/JavaScript unix timestamp e.g. 1656736672698
		if parseTimestampFloat(m, data, key) {
			break
		}
		//Handle iso date e.g. 2022-07-02T04:42:57.952Z
		if parseTimestampString(m, data, key) {
			break
		}
	}

	for _, key := range levelKeys {
		if !parseLevelString(m, data, key) {
			parseLevelInt(m, data, key)
		}
		if i, ok := levelNumberMap[m.Level]; ok {
			m.Incidence = i
		}
	}

	for _, key := range messageKeys {
		if parseMessage(m, data, key) {
			break
		}
	}

	for _, key := range callerKeys {
		if parseCaller(m, data, key) {
			break
		}
	}

	if len(data) > 0 {
		for key, val := range data {
			value, err := json.Marshal(val)
			if err != nil {
				m.Fields = append(m.Fields, &MessageField{
					Key:   key,
					Value: fmt.Sprintf("%+v", val),
				})
			} else {
				m.Fields = append(m.Fields, &MessageField{
					Key:   key,
					Value: string(value),
				})
			}
		}
		sort.Sort(sortableFields(m.Fields))
	}

	return m, nil
}

func parseMessage(m *Message, data MessageData, key string) bool {
	message, ok := data[key].(string)
	if ok {
		delete(data, key)
		m.Message = message
	}
	return ok
}

func parseCaller(m *Message, data MessageData, key string) bool {
	caller, ok := data[key].(string)
	if ok {
		delete(data, key)
		m.Caller = caller
	}
	return ok
}

func parseLevelString(m *Message, data MessageData, key string) bool {
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

func parseLevelInt(m *Message, data MessageData, key string) bool {
	level, ok := data[key].(float64)
	if !ok {
		return false
	}

	label, ok := numericLevelMap[int(level)]
	if !ok {
		return false
	}
	delete(data, key)
	m.Level = label
	return true
}

func parseTimestampFloat(m *Message, data MessageData, key string) bool {
	ts, ok := data[key].(float64)
	if ok {
		delete(data, key)
		sec, dec := math.Modf(ts)
		t := time.Unix(int64(sec), int64(dec*(1e9))).UTC()
		m.Timestamp = &t
	}
	return ok
}

var tsFormats = []string{
	time.RFC3339,
	time.UnixDate,
	time.Layout,
	time.ANSIC,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339Nano,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

func parseTimestampString(m *Message, data MessageData, key string) bool {
	ts, ok := data[key].(string)
	if ok {
		//TODO: we could iterate over all the formats
		for _, format := range tsFormats {
			if t, err := time.Parse(format, ts); err == nil {
				delete(data, key)
				m.Timestamp = &t
				return true
			}
		}
	}

	return false
}
