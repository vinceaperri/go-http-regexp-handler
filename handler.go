// Package handler provides RegexpHandler, a convient implementation of the
// http.Handler interface using regular expressions.
// 
// Before a RegexpHandler can serve requests, regular expression and function
// pairs (i.e. routes) need to be registered through Add. The handler can then
// serve a request by finding a route with a regular expression that matches
// its path and calling the route's function.
// 
// The order in which routes are registered is important. If a request's path
// matches the regular expresssions of multiple routes, the handler will only
// call the function of the route that was registered first.  This introduces a
// precedence heirachy that you may take advantage of.  Furthermore, if the
// path doesn't match any regular expression, the handler won't do anything. It
// is then useful to register a route with a catch-all regular expression (e.g.
// ".*") that could serve a 404 error.
package handler

import (
  "net/http"
  "regexp"
)

type route struct {
  re   *regexp.Regexp
  f    func(http.ResponseWriter, *http.Request, []string)
}

// RegexpHandler is an object that implements the http.Handler interface.
type RegexpHandler struct {
  routes []*route
}

// NewRegexpHandler creates a new RegexpHandler.
func NewRegexpHandler() *RegexpHandler {
  return &RegexpHandler{}
}

// Add registers a new regular expression and function pair, or route.  In
// addition to the typical parameters an http.HandlerFunc receives, the
// function will receive a slice of all submatches of the expression when
// matched with a request's path.
func (h *RegexpHandler) Add(expression string, function func(http.ResponseWriter, *http.Request, []string)) {
  re := regexp.MustCompile("^" + expression + "$")
  h.routes = append(h.routes, &route{re, function})
}

// ServeHTTP serves a request by calling the function of the first registered
// route containing an expression the request's path matches.
func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  for _, route := range h.routes {
    if matches := route.re.FindStringSubmatch(r.URL.Path); matches != nil {
      route.f(w, r, matches[1:])
      break
    }
  }
}
