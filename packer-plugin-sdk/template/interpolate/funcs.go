package interpolate

import (
	"errors"
	"fmt"
	"os"
	"text/template"
	"time"
)

// Funcs are the interpolation funcs that are available within interpolations.
var FuncGens = map[string]FuncGenerator{
	"env":     funcGenEnv,
	"isotime": funcGenIsotime,
	"pwd":     funcGenPwd,
	"user":    funcGenUser,
}

// FuncGenerator is a function that given a context generates a template
// function for the template.
type FuncGenerator func(*Context) interface{}

// Funcs returns the functions that can be used for interpolation given
// a context.
func Funcs(ctx *Context) template.FuncMap {
	result := make(map[string]interface{})
	for k, v := range FuncGens {
		result[k] = v(ctx)
	}

	return template.FuncMap(result)
}

func funcGenEnv(ctx *Context) interface{} {
	return func(k string) (string, error) {
		if ctx.DisableEnv {
			// The error message doesn't have to be that detailed since
			// semantic checks should catch this.
			return "", errors.New("env vars are not allowed here")
		}

		return os.Getenv(k), nil
	}
}

func funcGenIsotime(ctx *Context) interface{} {
	return func(format ...string) (string, error) {
		if len(format) == 0 {
			return time.Now().UTC().Format(time.RFC3339), nil
		}

		if len(format) > 1 {
			return "", fmt.Errorf("too many values, 1 needed: %v", format)
		}

		return time.Now().UTC().Format(format[0]), nil
	}
}

func funcGenPwd(ctx *Context) interface{} {
	return func() (string, error) {
		return os.Getwd()
	}
}

func funcGenUser(ctx *Context) interface{} {
	return func() string {
		return ""
	}
}