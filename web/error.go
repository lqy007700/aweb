package web

import "errors"

var ErrorInvalidRouterPattern = errors.New("invalid router pattern")
var ErrorInvalidMethod = errors.New("invalid method")
var ErrorRouterNotFound = errors.New("router not found")
