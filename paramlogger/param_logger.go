package paramlogger

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/gobuffalo/buffalo"
)

// ParameterExclusionList is the list of parameter names that will be filtered
// from the application logs (see maskSecrets).
// Important: this list will be used in case insensitive.
var ParameterExclusionList = []string{
	"Password",
	"PasswordConfirmation",
	"CreditCard",
	"CVC",
}

var filteredIndicator = []string{"[FILTERED]"}

// ParameterLogger logs form and parameter values to the logger
type parameterLogger struct {
	excluded []string
}

// ParameterLogger logs form and parameter values to the loggers
func ParameterLogger(next buffalo.Handler) buffalo.Handler {
	pl := parameterLogger{
		excluded: ParameterExclusionList,
	}

	return pl.Middleware(next)
}

// Middleware is a buffalo middleware function to connect this parameter filterer with buffalo
func (pl parameterLogger) Middleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		defer func() {
			req := c.Request()
			if req.Method != "GET" {
				if err := pl.logForm(c); err != nil {
					c.Logger().Error(err)
				}
			}

			params, ok := c.Params().(url.Values)
			if ok {
				params = pl.maskSecrets(params)
			}

			b, err := json.Marshal(params)
			if err != nil {
				c.Logger().Error(err)
			}

			c.LogField("params", string(b))
		}()

		return next(c)
	}
}

func (pl parameterLogger) logForm(c buffalo.Context) error {
	req := c.Request()
	mp := req.MultipartForm
	if mp != nil {
		return pl.multipartParamLogger(mp, c)
	}

	if err := pl.addFormFieldTo(c, req.Form); err != nil {
		return fmt.Errorf("unable to add form field %v: %w", req.Form, err)
	}

	return nil
}

func (pl parameterLogger) multipartParamLogger(mp *multipart.Form, c buffalo.Context) error {
	uv := url.Values{}
	for k, v := range mp.Value {
		for _, vv := range v {
			uv.Add(k, vv)
		}
	}
	for k, v := range mp.File {
		for _, vv := range v {
			uv.Add(k, vv.Filename)
		}
	}

	if err := pl.addFormFieldTo(c, uv); err != nil {
		return fmt.Errorf("unable to add form field %v: %w", uv, err)
	}
	return nil
}

func (pl parameterLogger) addFormFieldTo(c buffalo.Context, form url.Values) error {
	maskedForm := pl.maskSecrets(form)
	b, err := json.Marshal(maskedForm)

	if err != nil {
		return err
	}

	c.LogField("form", string(b))
	return nil
}

// maskSecrets matches ParameterExclusionList against parameters passed in the
// request, and returns a copy of the request parameters replacing excluded params
// with [FILTERED].
func (pl parameterLogger) maskSecrets(form url.Values) url.Values {
	if len(pl.excluded) == 0 {
		pl.excluded = ParameterExclusionList
	}

	copy := url.Values{}
	for key, values := range form {
	blcheck:
		for _, excluded := range pl.excluded {
			copy[key] = values
			if strings.ToUpper(key) == strings.ToUpper(excluded) {
				copy[key] = filteredIndicator
				break blcheck
			}

		}
	}
	return copy
}
