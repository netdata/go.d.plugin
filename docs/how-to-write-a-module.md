<!--
title: "How to write a Netdata collector in Go"
description: "This guide will walk you through the technical implementation of writing a new Netdata collector in Golang, with tips on interfaces, structure, configuration files, and more."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/docs/how-to-write-a-module.md"
sidebar_label: "How to write a Netdata collector in Go"
learn_status: "Published"
learn_topic_type: "Tasks"
learn_rel_path: "Developers"
sidebar_position: 20
-->

# How to write a Netdata collector in Go

Let's assume you want to write a collector named `example`.

Steps are:

- Add the source code to [`modules/example/`](https://github.com/netdata/go.d.plugin/tree/master/modules).
    - [module interface](#module-interface).
    - [suggested module layout](#module-layout).
    - [helper packages](#helper-packages).
- Add the configuration to [`config/go.d/example.conf`](https://github.com/netdata/go.d.plugin/tree/master/config/go.d).
- Add the module to [`config/go.d.conf`](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).
- Import the module in [`modules/init.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/init.go).
- Update the [`available modules list`](https://github.com/netdata/go.d.plugin#available-modules).

> :exclamation: If you prefer reading the source code, then check
> [the implementation](https://github.com/netdata/go.d.plugin/tree/master/modules/example) of the `example` module,
> it should give you an idea of  how things work.

## Module Interface

Every module should implement the following interface:

```
type Module interface {
    Init() bool
    Check() bool
    Charts() *Charts
    Collect() map[string]int64
    Cleanup()
}
```

### Init method

- `Init` does module initialization.
- If it returns `false`, the job will be disabled.

We propose to use the following template:

```
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

- `Check` returns whether the job is able to collect metrics.
- Called after `Init` and only if `Init` returned `true`.
- If it returns `false`, the job will be disabled.

The simplest way to implement `Check` is to see if we are getting any metrics from `Collect`. A lot of modules use such
approach.

```
// example.go

func (e *Example) Check() bool {
    return len(e.Collect()) > 0
}
```

### Charts method

:exclamation: Netdata module produces [`charts`](https://learn.netdata.cloud/docs/agent/collectors/plugins.d#chart), not
raw metrics.

Use [`agent/module`](https://github.com/netdata/go.d.plugin/blob/master/agent/module/charts.go) package to create them,
it contains charts and dimensions structs.

- `Charts` returns the [charts](https://learn.netdata.cloud/docs/agent/collectors/plugins.d#chart1) (`*module.Charts`).
- Called after `Check` and only if `Check` returned `true`.
- If it returns `nil`, the job will be disabled
- :warning: Make sure not to share returned value between module instances (jobs).

Usually charts initialized in `Init` and `Chart` method just returns the charts instance:

```
// example.go

func (e *Example) Charts() *Charts {
    return e.charts
}
```

### Collect method

- `Collect` collects metrics.
- Called only if `Check` returned `true`.
- Called every `update_every` seconds.
- `map[string]int64` keys are charts dimensions ids'.

We propose to use the following template:

```
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

- `Cleanup` performs the job cleanup/teardown.
- Called if `Init` or `Check` fails, or we want to stop the job after `Collect`.

If you have nothing to clean up:

```
// example.go

func (Example) Cleanup() {}
```

## Module Layout

The general idea is to not put everything in a single file.

We recommend using one file per logical area. This approach makes it easier to maintain the module.

Suggested minimal layout:

| Filename                                          | Contains                                               |
|---------------------------------------------------|--------------------------------------------------------|
| [`module_name.go`](#file-module_namego)           | Module configuration, implementation and registration. |
| [`charts.go`](#file-chartsgo)                     | Charts, charts templates and constructor functions.    |
| [`init.go`](#file-initgo)                         | Initialization methods.                                |
| [`collect.go`](#file-collectgo)                   | Metrics collection implementation.                     |
| [`module_name_test.go`](#file-module_name_testgo) | Public methods/functions tests.                        |
| [`testdata/`](#file-module_name_testgo)           | Files containing sample data.                          |

### File `module_name.go`

> :exclamation: See the example [`example.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/example/example.go).

Don't overload this file with the implementation details.

Usually it contains only:

- module registration.
- module configuration.
- [module interface implementation](#module-interface).

### File `charts.go`

> :exclamation: See the example: [`charts.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/example/charts.go).

Put charts, charts templates and charts constructor functions in this file.

### File `init.go`

> :exclamation: See the example: [`init.go`](https://github.com/netdata/go.d.plugin/blob/master/modules/example/init.go).

All the module initialization details should go in this file.

- make a function for each value that needs to be initialized.
- a function should return a value(s), not implicitly set/change any values in the main struct.

```
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

```
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

> if you have no experience in testing we recommend starting with [testing package documentation](https://golang.org/pkg/testing/).

> we use `assert` and `require` packages from [github.com/stretchr/testify](https://github.com/stretchr/testify) library,
> check [their documentation](https://pkg.go.dev/github.com/stretchr/testify).

Testing is mandatory.

- test only public functions and methods (`New`, `Init`, `Check`, `Charts`, `Cleanup`, `Collect`).
- do not create a test function per a case, use [table driven tests](https://github.com/golang/go/wiki/TableDrivenTests)
  . Prefer `map[string]struct{ ... }` over `[]struct{ ... }`.
- use helper functions _to prepare_ test cases to keep them clean and readable.

### Directory `testdata/`

Put files with sample data in this directory if you need any. Its name should
be [`testdata`](https://golang.org/cmd/go/#hdr-Package_lists_and_patterns).

> Directory and file names that begin with "." or "_" are ignored by the go tool, as are directories named "testdata".

## Helper packages

There are [some helper packages](https://github.com/netdata/go.d.plugin/tree/master/pkg) for writing a module.
