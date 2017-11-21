package errstack

import (
	"bytes"
	"fmt"
	"strconv"
)

type chain []interface{}
type errmap map[string]interface{}

const builderSep = "|"

// Builder is a type to incrementally build set of errors under common key structure
// Builder intentionally doesn't implement standard Error interface. You have to explicitly
// convert it into an Error (using ToReqErr) once all checks are done.
// Basic idea of builder is to easily combine request errors.
// Example:
//
//	var errb = NewBuilder()
//	if len(obj.Name) < 3 {
//		errb.Put("first_name", "name is too short")
//	}
//	if strings.Contains(obj.Name, "%") {
//		errb.Put("first_name", "name contains invalid characters")
//	}
//	...
//	return errb.ToReqErr()
type Builder interface {
	// Fork creates a new builder which shares the same space but all new added errors
	// will be assigned to keys prefixed with `prefix`
	Fork(prefix string) Builder
	// ForkIdx is a handy function to call Fork with an `int` (eg: index keys / rows)
	ForkIdx(idx int) Builder

	// Putter returns a Putter which abstract error setting from error key.
	Putter(key string) Putter

	// Puts new error under the key. You can put multiple errors under the same key
	// and they will be agregated
	Put(key string, value interface{})
	// Get returns errors under `key`
	Get(key string) interface{}

	// NotNil checks if there are any errors in the builder.
	NotNil() bool
	// Converts the Builder into a request error.
	ToReqErr() E
}

// Putter is an interface which provides a way to set an error abstracting from
// it's identifier (key). Please refer to the Builder documentation to see the
// example below without using Putter.
// Example:
//
//	validateName(obj.FirstName, errb.Putter("first_name"))
//	validateName(obj.LastName, errb.Putter("last_name"))
//
//	func validateName(name string, errp.Putter) {
//		if len(name) < 3 {
//			errp.Put(name is too short")
//		}
//		if strings.Contains(name, "%") {
//			errp.Put(name contains invalid characters")
//		}
//	}
type Putter interface {
	Put(interface{})
	Fork(prefix string) Putter
	ForkIdx(idx int) Putter
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

var errmapSep = []byte(": ")

// Error implements error interface
func (em errmap) Error() string {
	var buffer bytes.Buffer
	for k, v := range em {
		buffer.WriteString(k)
		buffer.Write(errmapSep)
		fmt.Fprintln(&buffer, v)
	}
	return buffer.String()
}

type builder struct {
	m      errmap
	prefix string
}

func (b builder) Fork(prefix string) Builder {
	prefix = prefix + builderSep
	if b.prefix != "" {
		prefix = b.prefix + prefix
	}
	return builder{b.m, prefix}
}

func (b builder) ForkIdx(idx int) Builder {
	return b.Fork(strconv.Itoa(idx))
}

func (b builder) Put(key string, value interface{}) {
	if value != nil {
		b.m.Append(b.prefix+key, value)
	}
}

func (b builder) Get(key string) interface{} {
	return b.m[key]
}

func (b builder) NotNil() bool {
	return len(b.m) > 0
}

func (b builder) ToReqErr() E {
	if b.NotNil() {
		return newRequest(b.m, 1)
	}
	return nil
}

func (b builder) Putter(key string) Putter {
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

func (bs builderSetter) ForkIdx(key int) Putter {
	return bs.Fork(strconv.Itoa(key))
}
