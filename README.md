# Eval

Eval is a library for the [Go Programming Language][go]. It provides some
operations to help testing operations.

## Status

[![Build Status](https://img.shields.io/travis/raiqub/eval/master.svg?style=flat&label=linux%20build)](https://travis-ci.org/raiqub/eval)
[![AppVeyor Build](https://img.shields.io/appveyor/ci/skarllot/eval/master.svg?style=flat&label=windows%20build)](https://ci.appveyor.com/project/skarllot/eval)
[![Coverage Status](https://coveralls.io/repos/raiqub/eval/badge.svg?branch=master&service=github)](https://coveralls.io/github/raiqub/eval?branch=master)
[![GoDoc](https://godoc.org/github.com/raiqub/eval?status.svg)](http://godoc.org/github.com/raiqub/eval)

## Features

 * **Environment** type which provides a Docker environment for testing purposes.
 * **MongoDBEnvironment** type which represents a MongoDB's Environment.

## Installation

To install raiqub/eval library run the following command:

```bash
go get gopkg.in/raiqub/eval.v0
```

To import this package, add the following line to your code:

```bash
import "gopkg.in/raiqub/eval.v0"
```

## Examples

Examples can be found on [library documentation][doc].

## Running tests

The tests can be run via the provided Bash script:

```bash
./test.sh
```

## License

raiqub/eval is made available under the [Apache Version 2.0 License][license].


[go]: http://golang.org/
[doc]: http://godoc.org/github.com/raiqub/eval
[license]: http://www.apache.org/licenses/LICENSE-2.0
