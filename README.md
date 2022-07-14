# lipence/errors

Package errors provide stacktrace, chain and data field support to errors in go.

This is particularly useful when you want to understand the state of execution when an error was returned unexpectedly.

It provides the type \*Error which implements the standard golang error interface, so you can use this library interchangably with code that is expecting a normal error return.

## Interface

**Message** is the basic message interface, all error objects implement this interface.

```go
type Message interface {
	Code() string
	Message() string
}
```

## Usage

Create and use an usual error object:

```go
// an error object which implements `Message`
var Err0001 = errors.New(
	"e0001",          // code
	"internal error", // message
)

func doSomething() error {
    // something goes wrong
	if _, ok := someData["key"]; !ok {
	    return Err0001 // just return error
	}
	// ......
}
```

Wrap an error (explain why error happens):

error object (which implement golang `error` interface) could be wrapped with `errors.Note` and `errors.Because`.

`errors.Because` accepts current error object as cause, and receives `Message` as explaination.

`errors.Note` simply receives current error object to trace stack.

```go
const FileMinLength = 10
var Err0001 = errors.New("e0001", "Can't read file")

func readSomething(path string) ([]byte, error) {
	if content, err := ioutil.ReadFile(path); err != nil {
		// return nil, errors.Note(Err0001)       // <= Just return Err0001
	    return nil, errors.Because(Err0001, err)  // <= Err0001's cause is err
	} else {
		return content, nil
	}
}
```

Annotate an error (collecting releated data):

`errors.Note` and `errors.Because` accepts dataFields as optional params. DataFields implements `errors.Field`, which is alias of zapcore.Field. Usage refers to [field.go](field.go)

```go
var Err0002 = errors.New("e0002", "File is empty")

func readSomethingElse(path string) ([]byte, error) {
	if content, err := ioutil.ReadFile(path); err != nil {
		// reason and dataField
	    return nil, errors.Because(Err0001, err, errors.String("path", path))
	} else if len(content) < FileMinLength {
		// only dataFields
		return nil, errors.Note(Err0002, errors.String("path", path),
			errors.Int("length", len(content)),
			errors.Int("minLength", FileMinLength))
	} else {
		return content, nil
	}
}
```

## Output Example

Console:

```log
2021-12-17T01:52:37.357024+0800 ERROR   main.go:81  
open /tmp/a.txt: no such file or directory
e0001: Can't read file: {"path":"/tmp/a.txt"}
  [4] main.readSomethingElse
    example.com/test/err-test/main.go:50
  [3] main.processSomething
    example.com/test/err-test/main.go:62
e0003: Cant Process File
  [3] main.processSomething
    example.com/test/err-test/main.go:63
  [2] main.main
    example.com/test/err-test/main.go:69
  [1] runtime.main
    runtime/proc.go:225
  [0] runtime.goexit
    runtime/asm_arm64.s:1130
```

JSON:

```json
[
    {
        "underlying": "open /tmp/a.txt: no such file or directory"
    },
    {
        "underlying": "e0001: Can't read file",
        "data": {
            "path": "/tmp/a.txt"
        },
        "stackTrace": [
            {
                "func": "[5] main.readSomethingElse",
                "line": "example.com/test/err-test/main.go:51"
            },
            {
                "func": "[4] main.processSomething",
                "line": "example.com/test/err-test/main.go:63"
            }
        ]
    },
    {
        "underlying": "e0003: Cant Process File",
        "stackTrace": [
            {
                "func": "[4] main.processSomething",
                "line": "example.com/test/err-test/main.go:64"
            },
            {
                "func": "[3] main.bootUp",
                "line": "example.com/test/err-test/main.go:82"
            },
            {
                "func": "[2] main.main",
                "line": "example.com/test/err-test/main.go:73"
            },
            {
                "func": "[1] runtime.main",
                "line": "runtime/proc.go:225"
            },
            {
                "func": "[0] runtime.goexit",
                "line": "runtime/asm_arm64.s:1130"
            }
        ]
    }
]
```

## Changelog

- v0.1.1 Built-in data fields Types and Factory
- v0.1.0 First Release

## Contact

Kenta Lee ( [kenta.li@cardinfolink.com](mailto:kenta@cardinfolink.com) )

## License

`lipence/errors` source code is available under the Apache-2.0 [License](/LICENSE)
