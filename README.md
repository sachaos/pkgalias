pkgalias
===

A linter for Go to fix alias of package.

# Install

```shell
$ go install github.com/sachaos/pkgalias
```

# Quickstart

## 1. Prepare config `.pkgalias.yaml`

```yaml
settings:
- alias: validatorv10
  fullpath: github.com/go-playground/validator/v10
```

## 2. Example source code

```go
package main

import "github.com/go-playground/validator/v10"

type User struct {
	Name string `validate:"required"`
}

func main() {
	user := &User{
		Name: "sachaos",
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		panic(err)
	}
}
```

## 3. Run pkgalias

```shell
$ pkgalias ./
./main.go:3:8: invalid alias: package name should be "validatorv10", insert this.
./main.go:14:14: invalid alias: use "validatorv10" instead of "validator"
```

## 4. Fix automatically with `-fix` option

```shell
$ pkgalias -fix ./
./main.go:3:8: invalid alias: package name should be "validatorv10", insert this.
./main.go:14:14: invalid alias: use "validatorv10" instead of "validator"
```

```go
package main

import validatorv10 "github.com/go-playground/validator/v10"

type User struct {
	Name string `validate:"required"`
}

func main() {
	user := &User{
		Name: "sachaos",
	}

	validate := validatorv10.New()
	if err := validate.Struct(user); err != nil {
		panic(err)
	}
}
```
