# Middleware

[![Standard Test](https://github.com/gobuffalo/middleware/actions/workflows/standard-go-test.yml/badge.svg)](https://github.com/gobuffalo/middleware/actions/workflows/standard-go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gobuffalo/middleware)](https://goreportcard.com/report/github.com/gobuffalo/middleware)
[![Go Reference](https://pkg.go.dev/badge/github.com/gobuffalo/middleware.svg)](https://pkg.go.dev/github.com/gobuffalo/middleware)

Middleware is a fundamental module that provides the default middleware for
buffalo applications.


## Installation

### with Buffalo

The default middleware is setup by default on a new Buffalo app.

```console
$ buffalo new myapp
```

Buffalo CLI will generate application scafflod and `actions/app.go` will be
configured to use default middlewares by default.

### Standalone

```console
$ go get github.com/gobuffalo/middleware
```


## Middlewares

### csrf

`csrf` middleware provides
[CSRF](https://en.wikipedia.org/wiki/Cross-site_request_forgery)
protection for Buffalo apps.
This middleware will be enabled by default in your app when you generate a new
application scaffold with `buffalo new` command.


### forcessl

`forcessl` middleware provides a feature that automatically redirects requests
that is not use HTTPS.
This middleware will be enabled by default in your app when you generate a new
application scaffold with `buffalo new` command.
It is configured to enforce the redirection in the `production` mode only. You
could customize it in `actions/app.go` if you need a different behavior.

### i18n

`i18n` middleware provides internationalization support in your application:

* User language detection from configurable sources
* Translation helper using locales bundles from github.com/nicksnyder/go-i18n
* Localized views

See <https://gobuffalo.io/documentation/guides/localization/> for further
information about Buffalo translation features and configuration.

### paramlogger

`paramlogger` middleware provides the request parameter logging feature.
This middleware will be enabled by default in your app when you generate a new
application scaffold with `buffalo new` command.

By default, it filters out pre-defined sensitive parameters and they will be
printed as `[FILTERED]` in the log. (keywords are case insensitive)

Currently, the pre-defined parameters are:

* password
* passwordconfirmation
* creditcard
* cvc

`paramlogger` also allows to exclude user-defined parameters by specifying
those in the `ParameterExclusionList` before using the ParamegerLogger
middleware in the app.

```go
    paramlogger.ParameterExclusionList = append(paramlogger.ParameterExclusionList, "secret", "phone_number")
    app.Use(paramlogger.ParameterLogger)
```

or

```go
    paramlogger.ParameterExclusionList = []string{
        "password",
        "passwordconfirmation",
        "secret",
        "phone_number",
    }
    app.Use(paramlogger.ParameterLogger)
```

## Documentation

* <https://gobuffalo.io/documentation/request_handling/middleware/>
* <https://gobuffalo.io/documentation/guides/localization/>

