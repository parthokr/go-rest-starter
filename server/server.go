package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Middleware func(next http.Handler) http.HandlerFunc
type apiFunc func(w http.ResponseWriter, r *http.Request) error

type Server struct {
	addr        string
	routes      map[string]apiFunc
	middlewares map[string][]Middleware
}

type Router struct {
	isMounted         bool
	isRegistered      bool
	routes            map[string]apiFunc
	localMiddlewares  map[string][]Middleware // middleware for individual routes
	globalMiddlewares []Middleware            // middleware for all routes
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, routes: make(map[string]apiFunc), middlewares: make(map[string][]Middleware)}
}

func NewRouter() *Router {
	return &Router{
		isRegistered:      false,
		isMounted:         false,
		routes:            make(map[string]apiFunc),
		localMiddlewares:  make(map[string][]Middleware),
		globalMiddlewares: make([]Middleware, 0),
	}
}

func (r *Router) Get(route string, middlewares []Middleware, f apiFunc) {
	key := fmt.Sprintf("GET %s", route)
	r.routes[key] = f
	r.localMiddlewares[key] = middlewares
}

func (r *Router) Post(route string, middlewares []Middleware, f apiFunc) {
	key := fmt.Sprintf("POST %s", route)
	r.routes[key] = f
	r.localMiddlewares[key] = middlewares
}

func (r *Router) Use(middlewares ...Middleware) {
	r.globalMiddlewares = middlewares
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHttpHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error
			// this is not necessarily for consumer of the API
			// this is for the developer
			log.Printf("Error: %s", err.Error())
		}
	}
}

func (s *Server) Mount(subroute string, r *Router) {
	// Mount routes with a subroute
	// You can not mount a router if it is already registered
	if r.isRegistered {
		panic("Registered router can not be mounted")
	}

	for route, f := range r.routes {
		method, path := func(route string) (string, string) {
			x := strings.SplitN(route, " ", 2)
			return x[0], fmt.Sprintf("%s%s", subroute, x[1])
		}(route)
		newKey := fmt.Sprintf("%s %s", method, path)
		s.routes[newKey] = f
		s.middlewares[newKey] = r.localMiddlewares[route]
		delete(r.routes, route)
		delete(r.localMiddlewares, route)
	}
	// prepend global middlewares
	for route, middlewares := range s.middlewares {
		s.middlewares[route] = append(r.globalMiddlewares, middlewares...)
	}
	// clear global middlewares
	r.globalMiddlewares = make([]Middleware, 0)
	r.isMounted = true
}

func (s *Server) Register(r *Router) {
	// Register all routes of the router
	// You can not register a router if it is already mounted
	if r.isMounted {
		panic("Mounted router can not be registered")
	}
	for route, f := range r.routes {
		s.routes[route] = f
		s.middlewares[route] = r.localMiddlewares[route]
		delete(r.routes, route)
		delete(r.localMiddlewares, route)
	}
	// prepend global middlewares
	for route, middlewares := range s.middlewares {
		s.middlewares[route] = append(r.globalMiddlewares, middlewares...)
	}
	// clear global middlewares
	r.globalMiddlewares = make([]Middleware, 0)
	r.isRegistered = true
}

func chainMiddlewares(middlewares []Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	for route, f := range s.routes {
		mux.HandleFunc(route, chainMiddlewares(s.middlewares[route])(makeHttpHandlerFunc(f)))
	}
	server := http.Server{
		Addr:    s.addr,
		Handler: mux,
	}
	log.Printf("Server running at %s", s.addr)
	err := server.ListenAndServe()
	return err
}
