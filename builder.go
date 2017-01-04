package errstack

import (
	"fmt"
)

type chain []interface{}
type errmap map[string]interface{}

const builderSep = "|"

// Builder is a type to incrementally build set of errors under common key structure
type Builder interface {
	Fork(subkey string) Builder
	Put(key string, value interface{})
	Get(key string) interface{}
	NotNil() bool
	ToReqErr() E
	Setter(key string) Putter
}

// Putter is an interface which provides a way to set an error abstracting from
// it's identifier (key)
type Putter interface {
	Put(interface{})
	Fork(prefix string) Putter
}

// Append chains errors under the same key
func (em errmap) Append(key string, value interface{}) {
	x, ok := em[key]
	if !ok {
		em[key] = value
		return
	}
	if ls, ok := x.(chain); ok || x == nil {
		em[key] = append(ls, value)
	} else {
		em[key] = chain{x, value}
	}
}

// Error implements error interface
func (em errmap) Error() string {
	return fmt.Sprint(map[string]interface{}(em))
}

type builder struct {
	m      errmap
	prefix string
}

// Fork creates a new builder which shares the same space but all new errors
// added will be assigned to keys prefixed with `prefix`
func (b builder) Fork(prefix string) Builder {
	prefix = prefix + builderSep
	if b.prefix != "" {
		prefix = b.prefix + prefix
	}
	return builder{b.m, prefix}
}

// Put adds new error. If err already exists under the same key,
// then it is appended.
func (b builder) Put(key string, value interface{}) {
	if value != nil {
		b.m.Append(b.prefix+key, value)
	}
}

// Get returns error under `key`
func (b builder) Get(key string) interface{} {
	return b.m[key]
}

// NotNil check if there are any errors in builder.
func (b builder) NotNil() bool {
	return len(b.m) > 0
}

// ToReqErr returns underlying error object
func (b builder) ToReqErr() E {
	if b.NotNil() {
		return newRequest(b.m, 1)
	}
	return nil
}

// Setter returns a Setter which abstract error setting from error key.
func (b builder) Setter(key string) Putter {
	return builderSetter{key, b}
}

// NewBuilder creates new builders with given prefix being appended to each error keys.
// This structure is not thread safe.
func NewBuilder() Builder {
	return builder{map[string]interface{}{}, ""}
}

type builderSetter struct {
	key string
	b   builder
}

func (bs builderSetter) Put(err interface{}) {
	bs.b.Put(bs.key, err)
}

func (bs builderSetter) Fork(key string) Putter {
	if bs.key != "" {
		key = bs.key + ":" + key
	}
	return builderSetter{key, bs.b}
}
