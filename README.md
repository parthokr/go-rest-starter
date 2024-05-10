## A starter template for a RESTful API using golang
### No frameworks, just the standard library

### Todo
- [x] Implement subrouters
- [x] Implement middleware at global and subrouter level
- [x] Implement middleware chaining
- [ ] Validate DTOs
- [ ] Implement a simple logger
- [ ] Create CRUD operations for a simple entity
- [ ] Implement a simple authentication mechanism
- [ ] Implement a simple authorization mechanism
- [ ] Role based access control

### Spin up the server
```go
package main

import (
	"fmt"
	"go-rest-starter/server"
	"net/http"
)

func main() {
	apiServer := server.NewServer(":8080")
	router := server.NewRouter()
	router.Get("/", []server.Middleware{}, func(w http.ResponseWriter, r *http.Request) error {
        fmt.Fprint(w, "Hello, World!")
        return nil
    })
	apiServer.Register(router)
	apiServer.Run()
}
```

#### You can make use of server.Middleware to add middleware to your routes
```go
...
router.Get("/", []server.Middleware{middleware1, middleware2}, func(w http.ResponseWriter, r *http.Request) error {
    fmt.Fprint(w, "Hello, World!")
    return nil
})
...
```
### or to the entire router
```go
...
router.Use(middleware1, middleware2, middleware3)
...
```
### So what is a middleware?
#### server.Middleware is an alias for `func(next http.Handler) http.HandlerFunc`
A middleware to log the request method and path
```go
...
func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method %s path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
...
```