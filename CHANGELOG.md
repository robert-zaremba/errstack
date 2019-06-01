# v1

Extended E interface:

+ Added `StatusCode`
+ Added `Kind` and Kind values, which mostly follow the `upspin.io/errors`
+ Fixed `MarshalJSON` methods
+ Added support to the `fmt.Formatter` interface. All `New*` and `Wrap*` functions return a value which can have extended printing. In essence: `fmt.Fprinf("%+v", e)` will print the layered message with a stacktrace.
+ Added new `New` method which takes a new values as an argument: `Kind`.
+ Added new `Details` and `AddDetail` method - this is used to add more information to the error

Added new module functions:

+ `IsKind` - again, follows the `upspin.io/errors.Is` implementation
+ `RootError` - which comes in place of module `Cause` function

Breaking changes:

+ module `Cause` function was renamed to `RootError`
+ updated `Error` string format
+ Change `HasUnderlying` interface to pkg/errors.Causer

Internal changes:

+ Use `stack.Stack` instead of `stack.Multi`
