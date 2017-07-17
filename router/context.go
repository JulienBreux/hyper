package router

import (
	"context"
	"net/http"
	"time"

	"github.com/vaniila/hyper/message"
	"github.com/vaniila/hyper/router/cookie"
	"github.com/vaniila/hyper/tracer"
)

// HandlerFunc type
type HandlerFunc func(Context)

// HandlerFuncs type
type HandlerFuncs []HandlerFunc

// Context interface
type Context interface {
	Identity() Identity
	MachineID() string
	ProcessID() string
	Context() context.Context
	Req() *http.Request
	Res() http.ResponseWriter
	Client() Client
	Cache() CacheAdaptor
	Message() MessageAdaptor
	Cookie() Cookie
	Header() Header
	MustParam(s string) Value
	MustQuery(s string) Value
	MustBody(s string) Value
	Param(s string) (Value, error)
	Query(s string) (Value, error)
	Body(s string) (Value, error)
	File(s string) []byte
	Tracer() tracer.Tracer
	Recover() error
	Abort()
	IsAborted() bool
	Write(b []byte) Context
	Error(error) Context
	Json(o interface{}) Context
	Status(code int) Context
}

// Identity interface
type Identity interface {
	HasID() bool
	GetID() int
	SetID(int)
	HasKey() bool
	GetKey() string
	SetKey(string)
}

// CacheAdaptor interface
type CacheAdaptor interface {
	Set(key []byte, data []byte, ttl time.Duration) error
	Get(key []byte) ([]byte, error)
}

// MessageAdaptor broker interface
type MessageAdaptor interface {
	Emit([]byte, []byte) error
	Listen([]byte, message.Handler) message.Close
}

// Value interface
type Value interface {
	Key() string
	Val() []byte
	Has() bool
	MustInt() int
	MustI32() int32
	MustI64() int64
	MustU32() uint32
	MustU64() uint64
	MustF32() float32
	MustF64() float64
	MustBool() bool
	MustTime() time.Time
	String() string
}

// Header interface
type Header interface {
	Set(key string, val string)
	Get(key string) (Value, error)
	MustGet(key string) Value
}

// Cookie interface
type Cookie interface {
	Set(key string, val string, opts ...cookie.Option)
	Get(key string) (Value, error)
	MustGet(key string) Value
}

// Client interface
type Client interface {
	IP() string
	Host() string
	Device() Device
	Protocol() string
}

// Device interface
type Device interface {
	Useragent() Useragent
	Hardware() Hardware
	OS() OS
}

// Useragent interface
type Useragent interface {
	Family() string
	Major() string
	Minor() string
	Patch() string
	String() string
}

// OS interface
type OS interface {
	Family() string
	Major() string
	Minor() string
	Patch() string
}

// Hardware interface
type Hardware interface {
	Family() string
	Brand() string
	Model() string
}
