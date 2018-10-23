# go-mockgen

[![GoDoc](https://godoc.org/github.com/efritz/go-mockgen?status.svg)](https://godoc.org/github.com/efritz/go-mockgen)
[![Build Status](https://secure.travis-ci.org/efritz/go-mockgen.png)](http://travis-ci.org/efritz/go-mockgen)
[![Maintainability](https://api.codeclimate.com/v1/badges/8546037d609e215de82d/maintainability)](https://codeclimate.com/github/efritz/go-mockgen/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/8546037d609e215de82d/test_coverage)](https://codeclimate.com/github/efritz/go-mockgen/test_coverage)

A mock interface code generator.

## Testing with Mocks

More usage documentation coming soon.

## Generating Mocks

Install with `go get -u github.com/efritz/go-mockgen/...`.

More usage documentation coming soon.

### Flags

The following flags are defined by the binary.

| Name       | Short Flag | Description  |
| ---------- | ---------- | ------------ |
| package    | p          | The name of the generated package. Is the name of target directory if dirname or filename is supplied by default. |
| prefix     |            | A prefix used in the name of each mock struct. Should be TitleCase by convention. |
| interfaces | i          | A whitelist of interfaces to generate given the import paths. |
| filename   | o          | The target output file. All mocks are writen to this file. |
| dirname    | d          | The target output directory. Each mock will be written to a unique file. |
| force      | f          | Do not abort if a write to disk would overwrite an existing file. |

If neither dirname nor filename are supplied, then the generated code is printed to standard out.

## License

Copyright (c) 2018 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
