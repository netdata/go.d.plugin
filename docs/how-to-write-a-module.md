# How to write a module

Let's assume you want to write `example` module.

Steps are:

- add the source code to the [`modules/example/`](https://github.com/netdata/go.d.plugin/tree/master/modules).
  - [module interface](#module-interface).
  - [suggested module layout](#module-layout).
  - [helper packages](#helper-packages).
- add the configuration to the [`config/go.d/exmaple.conf`](https://github.com/netdata/go.d.plugin/tree/master/config/go.d).
- add the module to the [`config/go.d.conf`](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).
- import the module in [`modules/init.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/init.go).
- update the [`available modules list`](https://github.com/netdata/go.d.plugin#available-modules).
 
> :exclamation: If you prefer reading the source code, then check 
> [the implementation](https://github.com/netdata/go.d.plugin/tree/master/modules/example) of the `example` module,
> it should give you an idea of  how things work. 

## Module Interface

Every module should implement the following interface:

```go
type Module interface {
    Init() bool
    Check() bool
    Charts() *Charts
    Collect() map[string]int64
    Cleanup()
}
```

### Init method

-   `Init` does module initialization.
-   If it returns `false`, the job will be disabled.

We propose to use the following template:

```go
// example.go

func (e *Example) Init() bool {
    err := e.validateConfig()
    if err != nil {
        e.Errorf("config validation: %v", err)
        return false
    }

    someValue, err := e.initSomeValue()
    if err != nil {
        e.Errorf("someValue init: %v", err)
        return false
    }
    e.someValue = someValue

    // ...
    return true 
}
```

Move specific initialization methods into the `init.go` file. See [suggested module layout](#module-Layout).

### Check method

-   `Check` returns whether the job is able to collect metrics.
-   Called after `Init` and only if `Init` returned `true`.
-   If it returns `false`, the job will be disabled.

The simplest way to implement `Check` is to see if we are getting any metrics from `Collect`.
A lot of modules use such approach.

```go
// example.go

func (e *Example) Check() bool {
    return len(e.Collect()) > 0
}
```

### Charts method

Netdata module produces [`charts`](https://learn.netdata.cloud/docs/agent/collectors/plugins.d#chart), not raw metrics.

Use [`agent/module`](https://github.com/netdata/go.d.plugin/blob/master/agent/module/charts.go) package to create them.

-   `Charts` returns the charts' definition.
-   Called after `Check` and only if `Check` returned `true`.
-   If it returns `nil`, the job will be disabled
-   :warning: Make sure not to share returned value between module instances (jobs).

Usually charts initialized in `Init` and `Chart` method just returns the charts instance:

```go
// example.go

func (e *Example) Charts() *Charts {
    return e.charts
}
```

### Collect method

-   `Collect` collects metrics.
-    Called only if `Check` returned `true`.
-    Called every `update_every` seconds.
-    `map[string]int64` keys are charts dimensions ids'.

We propose to use the following template:
 
```go
// example.go

func (e *Example) Collect() map[string]int64 {
    ms, err := e.collect()
    if err != nil {
        e.Error(err)
    }

    if len(ms) == 0 {
        return nil
    }
    return ms
}
```

Move metrics collection logic into the `collect.go` file. See [suggested module layout](#module-Layout).

### Cleanup method

-   `Cleanup` performs the job cleanup/teardown.
-    Called if `Init` or `Check` fails, or we want to stop the job after `Collect`.


If you have nothing to clean up:

```go
// example.go

func (Example) Cleanup() {}
```

## Module Layout

The general idea - do not put everything in a file.

We recommend using a file per a logical area. This approach makes it easier to maintain the module.

Suggested minimal layout:

| Filename                                          | Contains                                               |
| ------------------------------------------------- |------------------------------------------------------- |
| [`module_name.go`](#file-module_namego)           | Module configuration, implementation and registration. |
| [`charts.go`](#file-chartsgo)                     | Charts, charts templates and constructor functions.    |
| [`init.go`](#file-initgo)                         | Initialization methods.                                |
| [`collect.go`](#file-collectgo)                   | Metrics collection implementation.                     |
| [`module_name_test.go`](#file-module_name_testgo) | Public methods/functions tests.                        |
| [`testdata/`](#file-module_name_testgo)           | Tests fixtures.                                        |

### File `module_name.go`

> :exclamation: See the example [`examlpe.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/example/example.go).

Don't overload this file with the implementation details.

Usually it contains only:
 
-   module registration.
-   module configuration.
-   [module interface implementation](#module-interface).

### File `charts.go`

> :exclamation: See the example: [`charts.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/example/charts.go).

Put charts, charts templates and charts constructor functions in this file. 

### File `init.go`

> :exclamation: See the example: [`init.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/example/init.go).

All the module initialization details should go in this file.

-   make a function for each value that needs to be initialized.
-   a function should return a value(s), not implicitly set/change any values in the main struct.

```go
// init.go

// Prefer this approach.
func (e Example) initSomeValue() (someValue, error) {
    // ...
    return someValue, nil 
}

// This approach is ok too, but we recommend to not use it.
func (e *Example) initSomeValue() error {
    // ...
    m.someValue = someValue
    return nil
}
```     

### File `collect.go`

> :exclamation: See the example: [`collect.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/example/collect.go).

This file is the entry point for the metrics collection.

Feel free to split it into several files if you think it makes the code more readable.

Use `collect_` prefix for the filenames: `collect_this.go`, `collect_that.go`, etc.

```go
// collect.go

func (e *Example) collect() (map[string]int64, error) {
    collected := make(map[string])int64
    // ...
    // ...
    // ...
    return collected, nil
}
```

### File `module_name_test.go`

> :exclamation: See the example: [`example_test.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/example/example_test.go).

Testing is mandatory.

-   test only public functions and methods (`New`, `Init`, `Check`, `Charts`, `Cleanup`, `Collect`).
-   do not create a test function per a case, use [table driven tests](https://github.com/golang/go/wiki/TableDrivenTests).
Prefer `map[string]struct{ ... }` over `[]struct{ ... }`.
-   use helper functions _to prepare_ test cases to keep them clean and readable.  

### Directory `testdata/`

The directory contains tests fixtures.
Its name should be [`testdata`](https://golang.org/cmd/go/#hdr-Package_lists_and_patterns).

> Directory and file names that begin with "." or "_" are ignored by the go tool, as are directories named "testdata".

## Helper packages

There are [some helper packages](https://github.com/netdata/go.d.plugin/tree/master/pkg) for writing a module.

-   if you need IP ranges consider to use [`iprange`](https://github.com/netdata/go.d.plugin/tree/master/pkg/iprange#iprange).
-   if you parse an application log files, then [`log`](https://github.com/netdata/go.d.plugin/tree/master/pkg/logs) is handy.
-   if you need filtering check [`matcher`](https://github.com/netdata/go.d.plugin/tree/master/pkg/matcher#supported-format).
-   if you collect metrics from an HTTP endpoint use [`web`](https://github.com/netdata/go.d.plugin/tree/master/pkg/web).
-   if you collect metrics from a prometheus endpoint, then [`prometheus`](https://github.com/netdata/go.d.plugin/tree/master/pkg/prometheus)
and [`web`](https://github.com/netdata/go.d.plugin/tree/master/pkg/web) is what you need.
-   [`tlscfg`](https://github.com/netdata/go.d.plugin/tree/master/pkg/tlscfg) provides TLS support.
-   [`stm`](https://github.com/netdata/go.d.plugin/tree/master/pkg/stm) helps you to convert any struct to a `map[string]int64`.
