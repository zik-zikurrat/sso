package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

const (
	defaultlLevel     = slog.LevelDebug
	defaultTimeFormat = time.StampMilli
)

type (
	PrettyJSONHandler struct {
		writer io.Writer
		opts   *Options
		mu     *sync.Mutex
		buf    *sync.Pool
	}
	Options struct {
		AddSource   bool
		Level       slog.Leveler
		ReplaceAttr func(groups []string, attr slog.Attr) slog.Attr
		TimeFormat  string
		NoColor     bool
	}
)

func (o *Options) setDefault() {
	if o.Level == nil {
		o.Level = defaultlLevel
	}
	if o.TimeFormat == "" {
		o.TimeFormat = defaultTimeFormat
	}
}

func NewPrettyJSONHandler(w io.Writer, opts *Options) *PrettyJSONHandler {
	if opts == nil {
		opts = &Options{}
	}
	opts.setDefault()
	return &PrettyJSONHandler{
		writer: w,
		opts:   opts,
		mu:     &sync.Mutex{},
		buf: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

func (h *PrettyJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := h.buf.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		h.buf.Put(buf)
	}()

	buf.WriteString("{\n")

	first := true

	writeField := func(key string, quotedValue string) {
		if !first {
			buf.WriteString(",\n")
		}
		first = false
		buf.WriteString("  ")
		buf.WriteString(strconv.Quote(key))
		buf.WriteString(": ")
		buf.WriteString(quotedValue)
	}

	writeField("time", strconv.Quote(r.Time.Format(h.opts.TimeFormat)))
	writeField("level", strconv.Quote(r.Level.String()))
	writeField("msg", strconv.Quote(r.Message))

	r.Attrs(func(a slog.Attr) bool {
		writeField(a.Key, h.formatValue(a.Value))
		return true
	})

	buf.WriteString("\n}\n")

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.writer.Write(buf.Bytes())
	return err
}

func (h *PrettyJSONHandler) formatValue(v slog.Value) string {
	v = v.Resolve()
	switch v.Kind() {
	case slog.KindString:
		return strconv.Quote(v.String())
	case slog.KindInt64:
		return strconv.FormatInt(v.Int64(), 10)
	case slog.KindUint64:
		return strconv.FormatUint(v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.FormatFloat(v.Float64(), 'f', -1, 64)
	case slog.KindBool:
		return strconv.FormatBool(v.Bool())
	case slog.KindDuration:
		return strconv.Quote(v.Duration().String())
	case slog.KindTime:
		return strconv.Quote(v.Time().Format(h.opts.TimeFormat))
	case slog.KindAny:
		data, err := json.Marshal(v.Any())
		if err != nil {
			return strconv.Quote(fmt.Sprintf("%v", v.Any()))
		}
		return string(data)
	case slog.KindGroup:
		return `"<group>"`
	default:
		return strconv.Quote(fmt.Sprintf("%v", v.Any()))
	}
}

func (h *PrettyJSONHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *PrettyJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyJSONHandler) WithGroup(name string) slog.Handler {
	return h
}
