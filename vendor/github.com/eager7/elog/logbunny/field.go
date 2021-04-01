package logbunny

import (
	"fmt"
	"math"
	"runtime"
	"strings"
	"sync"
	"time"
)

// define the stack skip about AddCaller

var (
	CallerSkip  = 1
	_callerSkip = 2
)

// Field type defined the field type matching
type fieldType int

const (
	unknownType fieldType = iota
	boolType
	floatType
	intType
	int64Type
	uintType
	uint64Type
	uintptrType
	stringType
	timeType
	durationType
	objectType
	stringerType
	errorType
	skipType
)

// Field different filed may not fill in specific type
// Why use Field but with the lower-case field? I just followed what zap dose!
type Field struct {
	key     string
	valType fieldType

	vInt    int64       //int uint int64 uint64 uintprt time duration
	vString string      //string  error error bytes
	vObj    interface{} //object marshaler
}

// fieldPool is used to reuse the filed without realloc to save the time
var _fieldPool = sync.Pool{
	New: func() interface{} {
		return new(Field)
	},
}

// Recycle will put back the object into the pool. Make sure the object you wanna recycle
// has noone hold. By default you should not call this function only if you
// actually know the effection
func Recycle(obj *Field) {
	_fieldPool.Put(obj)
}

// Skip constructs a no-op Field.
func Skip() *Field {
	f := _fieldPool.Get().(*Field)
	f.valType = skipType
	return f
}

// Bool constructs a Field with the given key and value. Bools are marshaled
// lazily.
func Bool(key string, val bool) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = boolType
	f.vInt = 0
	if val {
		f.vInt = 1
	}
	return f
}

// Float64 constructs a Field with the given key and value. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float64(key string, val float64) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = floatType
	f.vInt = int64(math.Float64bits(val))
	return f
}

// Int constructs a Field with the given key and value.
func Int(key string, val int) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = intType
	f.vInt = int64(val)
	return f
}

// Int64 constructs a Field with the given key and value. Like ints
func Int64(key string, val int64) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = int64Type
	f.vInt = val
	return f
}

// Uint constructs a Field with the given key and value.
func Uint(key string, val uint) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = uintType
	f.vInt = int64(val)
	return f
}

// Uint64 constructs a Field with the given key and value.
func Uint64(key string, val uint64) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = uint64Type
	f.vInt = int64(val)
	return f
}

// Uintptr constructs a Field with the given key and value.
func Uintptr(key string, val uintptr) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = uintptrType
	f.vInt = int64(val)
	return f
}

// String constructs a Field with the given key and value.
func String(key string, val string) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = stringType
	f.vString = val
	return f
}

// Stringer constructs a Field with the given key and the output of the value's String method
func Stringer(key string, val fmt.Stringer) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = stringerType
	f.vString = ""
	if val != nil {
		f.vString = val.String()
	}
	return f
}

// Time constructs a Field with the given key and value. It represents a
// time.Time as a floating-point number of seconds since epoch. Conversion to a
// float64 happens eagerly.
func Time(key string, val time.Time) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = timeType
	f.vObj = val
	return f
}

// Err constructs a Field that lazily stores err.Error() under the key "error".
func Err(err error) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = "error"
	f.valType = errorType
	f.vString = ""
	if err != nil {
		f.vString = err.Error()
	}
	return f
}

// Duration constructs a Field with the given key and value. It represents
// durations as an integer number of nanoseconds.
func Duration(key string, val time.Duration) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = durationType
	f.vInt = int64(val)
	return f
}

// Object constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy.
func Object(key string, val interface{}) *Field {
	f := _fieldPool.Get().(*Field)
	f.key = key
	f.valType = objectType
	f.vObj = val
	return f
}

// Caller will add the line & file name to logger. The caller info is refer to goroutine
// stack information which can be get by runtime.Caller(_skip). _callerSkip is a index
// defined the layer about the stack info which contains the file name and line number.
// This could be unstable while the call stack changed. But mostly we got a stable value in it.
func Caller() *Field {
	_, filename, line, ok := runtime.Caller(CallerSkip)
	if !ok {
		return nil
	}
	fileNames := strings.Split(filename, "/")
	filename = fileNames[len(fileNames)-1]

	f := _fieldPool.Get().(*Field)
	f.key = "caller_info"
	f.valType = stringType
	f.vString = fmt.Sprintf("%s:%d", filename, line)
	return f
}

func caller() *Field {
	_, filename, line, ok := runtime.Caller(_callerSkip)
	if !ok {
		return nil
	}
	fileNames := strings.Split(filename, "/")
	filename = fileNames[len(fileNames)-1]

	f := _fieldPool.Get().(*Field)
	f.key = "caller_info"
	f.valType = stringType
	f.vString = fmt.Sprintf("%s:%d", filename, line)
	return f
}

func SetCallerSkip(newSkip int) {
	if newSkip < 0 {
		return
	}
	_callerSkip = newSkip
}
