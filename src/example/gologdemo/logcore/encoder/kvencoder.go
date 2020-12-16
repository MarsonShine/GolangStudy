package encoder

// adapted from https://github.com/uber-go/zap/blob/master/zapcore/json_encoder.go
// and https://github.com/uber-go/zap/blob/master/zapcore/console_encoder.go

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"sync"
	"time"
	"unicode/utf8"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const _hex = "0123456789abcdef"

var bufferPool = buffer.NewPool()

var _kvPool = sync.Pool{New: func() interface{} {
	return &kvEncoder{}
}}

func getKVEncoder() *kvEncoder {
	return _kvPool.Get().(*kvEncoder)
}

func putKVEncoder(enc *kvEncoder) {
	enc.EncoderConfig = nil
	enc.buf = nil
	_kvPool.Put(enc)
}

type kvEncoder struct {
	*zapcore.EncoderConfig
	buf *buffer.Buffer
}

// NewkvEncoder creates a key=value encoder
func NewKVEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &kvEncoder{
		EncoderConfig: &cfg,
		buf:           bufferPool.Get(),
	}
}

func (enc *kvEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *kvEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *kvEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *kvEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
}

func (enc *kvEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
}

func (enc *kvEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
}

func (enc *kvEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
}

func (enc *kvEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
}

func (enc *kvEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
}

func (enc *kvEncoder) AddReflected(key string, obj interface{}) error {
	marshaled, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(marshaled)
	return err
}

func (enc *kvEncoder) OpenNamespace(key string) {
}

func (enc *kvEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
}

func (enc *kvEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
}

func (enc *kvEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
}

func (enc *kvEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	return arr.MarshalLogArray(enc)
}

func (enc *kvEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	return obj.MarshalLogObject(enc)
}

func (enc *kvEncoder) AppendBool(val bool) {
	enc.buf.AppendBool(val)
}

func (enc *kvEncoder) AppendByteString(val []byte) {
	enc.safeAddByteString(val)
}

func (enc *kvEncoder) AppendComplex128(val complex128) {
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *kvEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *kvEncoder) AppendInt64(val int64) {
	enc.buf.AppendInt(val)
}

func (enc *kvEncoder) AppendReflected(val interface{}) error {
	marshaled, err := json.Marshal(val)
	if err != nil {
		return err
	}
	_, err = enc.buf.Write(marshaled)
	return err
}

func (enc *kvEncoder) AppendString(val string) {
	enc.safeAddString(val)
}

func (enc *kvEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(val, enc)
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *kvEncoder) AppendUint64(val uint64) {
	enc.buf.AppendUint(val)
}

func (enc *kvEncoder) AddComplex64(k string, v complex64) { enc.AddComplex128(k, complex128(v)) }
func (enc *kvEncoder) AddFloat32(k string, v float32)     { enc.AddFloat64(k, float64(v)) }
func (enc *kvEncoder) AddInt(k string, v int)             { enc.AddInt64(k, int64(v)) }
func (enc *kvEncoder) AddInt32(k string, v int32)         { enc.AddInt64(k, int64(v)) }
func (enc *kvEncoder) AddInt16(k string, v int16)         { enc.AddInt64(k, int64(v)) }
func (enc *kvEncoder) AddInt8(k string, v int8)           { enc.AddInt64(k, int64(v)) }
func (enc *kvEncoder) AddUint(k string, v uint)           { enc.AddUint64(k, uint64(v)) }
func (enc *kvEncoder) AddUint32(k string, v uint32)       { enc.AddUint64(k, uint64(v)) }
func (enc *kvEncoder) AddUint16(k string, v uint16)       { enc.AddUint64(k, uint64(v)) }
func (enc *kvEncoder) AddUint8(k string, v uint8)         { enc.AddUint64(k, uint64(v)) }
func (enc *kvEncoder) AddUintptr(k string, v uintptr)     { enc.AddUint64(k, uint64(v)) }
func (enc *kvEncoder) AppendComplex64(v complex64)        { enc.AppendComplex128(complex128(v)) }
func (enc *kvEncoder) AppendFloat64(v float64)            { enc.appendFloat(v, 64) }
func (enc *kvEncoder) AppendFloat32(v float32)            { enc.appendFloat(float64(v), 32) }
func (enc *kvEncoder) AppendInt(v int)                    { enc.AppendInt64(int64(v)) }
func (enc *kvEncoder) AppendInt32(v int32)                { enc.AppendInt64(int64(v)) }
func (enc *kvEncoder) AppendInt16(v int16)                { enc.AppendInt64(int64(v)) }
func (enc *kvEncoder) AppendInt8(v int8)                  { enc.AppendInt64(int64(v)) }
func (enc *kvEncoder) AppendUint(v uint)                  { enc.AppendUint64(uint64(v)) }
func (enc *kvEncoder) AppendUint32(v uint32)              { enc.AppendUint64(uint64(v)) }
func (enc *kvEncoder) AppendUint16(v uint16)              { enc.AppendUint64(uint64(v)) }
func (enc *kvEncoder) AppendUint8(v uint8)                { enc.AppendUint64(uint64(v)) }
func (enc *kvEncoder) AppendUintptr(v uintptr)            { enc.AppendUint64(uint64(v)) }

func (enc *kvEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *kvEncoder) clone() *kvEncoder {
	clone := getKVEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.buf = bufferPool.Get()
	return clone
}

func (enc *kvEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()
	final.addElementSeparator()
	if final.LevelKey != "" {
		final.addKey(final.LevelKey)
		cur := final.buf.Len()
		final.EncodeLevel(ent.Level, final)
		if cur == final.buf.Len() {
			final.AppendString(ent.Level.String())
		}
		final.addElementSeparator()
	}
	if final.TimeKey != "" {
		final.AddTime(final.TimeKey, ent.Time)
		final.addElementSeparator()
	}
	if ent.LoggerName != "" && final.NameKey != "" {
		final.addKey(final.NameKey)
		cur := final.buf.Len()
		nameEncoder := final.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}

		nameEncoder(ent.LoggerName, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeName was a no-op. Fall back to strings to
			// keep output valid.
			final.AppendString(ent.LoggerName)
		}
		final.addElementSeparator()
	}
	if ent.Caller.Defined && final.CallerKey != "" {
		final.addKey(final.CallerKey)
		cur := final.buf.Len()
		final.EncodeCaller(ent.Caller, final)
		if cur == final.buf.Len() {
			// User-supplied EncodeCaller was a no-op. Fall back to strings to
			// keep JSON valid.
			final.AppendString(ent.Caller.String())
		}
		final.addElementSeparator()
	}
	if final.MessageKey != "" {
		final.addKey(enc.MessageKey)
		final.buf.AppendByte('"')
		final.AppendString(ent.Message)
		final.buf.AppendByte('"')
		final.addElementSeparator()
	}
	if enc.buf.Len() > 0 {
		final.buf.Write(enc.buf.Bytes())
	}
	addFields(final, final, fields)
	final.addElementSeparator()
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
		final.addElementSeparator()
	}
	if final.LineEnding != "" {
		final.buf.AppendString(final.LineEnding)
	} else {
		final.buf.AppendString(zapcore.DefaultLineEnding)
	}

	ret := final.buf
	putKVEncoder(final)
	return ret, nil
}

func (enc *kvEncoder) addKey(key string) {
	enc.buf.AppendString(key)
	enc.buf.AppendByte('=')
}

func (enc *kvEncoder) addElementSeparator() {
	enc.buf.AppendByte('#')
}

func (enc *kvEncoder) appendFloat(val float64, bitSize int) {
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *kvEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *kvEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *kvEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *kvEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func addFields(kvEnc *kvEncoder, enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
		kvEnc.buf.AppendByte('#')
	}
}
