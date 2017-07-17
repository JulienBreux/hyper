package router

import (
	"fmt"
	"log"
)

type router struct {
	pat           string
	aliases       []string
	name          string
	method        string
	routes        []Route
	namespace     bool
	ws            bool
	http          bool
	documentation string
	summary       string
	params        []Param
	memory        int64
	middleware    HandlerFuncs
	handler       HandlerFunc
	catch         HandlerFunc
	models        []Model
}

// response struct
type model struct {
	code      int
	structure interface{}
}

func (v *router) add(pat, method string) Route {
	// pass down params
	var params []Param
	if l := len(v.params); l > 0 {
		for _, param := range v.params {
			params = append(params, param)
		}
	}
	// pass down middleware
	var middleware HandlerFuncs
	if l := len(v.middleware); l > 0 {
		for _, handler := range v.middleware {
			middleware = append(middleware, handler)
		}
	}
	r := &router{
		pat:        pat,
		method:     method,
		ws:         true,
		http:       true,
		params:     params,
		memory:     v.memory,
		middleware: middleware,
	}
	v.routes = append(v.routes, r)
	return r
}

func (v *router) Get(pat string) Route {
	return v.add(pat, "GET")
}

func (v *router) Head(pat string) Route {
	return v.add(pat, "HEAD")
}

func (v *router) Options(pat string) Route {
	return v.add(pat, "OPTIONS")
}

func (v *router) Post(pat string) Route {
	return v.add(pat, "POST")
}

func (v *router) Put(pat string) Route {
	return v.add(pat, "PUT")
}

func (v *router) Patch(pat string) Route {
	return v.add(pat, "PATCH")
}

func (v *router) Delete(pat string) Route {
	return v.add(pat, "DELETE")
}

func (v *router) Namespace(pat string) Route {
	if !v.namespace {
		log.Fatalf("Route %s is not a namespace, you are only allowed to attach route(s) to namespaces.", v.pat)
	}
	// pass down params
	var params []Param
	if l := len(v.params); l > 0 {
		for _, param := range v.params {
			params = append(params, param)
		}
	}
	// pass down middleware
	var middleware HandlerFuncs
	if l := len(v.middleware); l > 0 {
		for _, handler := range v.middleware {
			middleware = append(middleware, handler)
		}
	}
	r := &router{
		pat:        pat,
		namespace:  true,
		params:     params,
		memory:     v.memory,
		middleware: middleware,
	}
	v.routes = append(v.routes, r)
	return r
}

func (v *router) Alias(aliases ...string) Route {
	for _, alias := range aliases {
		v.aliases = append(v.aliases, alias)
	}
	return v
}

func (v *router) Name(s string) Route {
	v.name = s
	return v
}

func (v *router) Summary(s string) Route {
	v.summary = s
	return v
}

func (v *router) Doc(s string) Route {
	v.documentation = s
	return v
}

func (v *router) Params(ps ...Param) Route {
	var params = make(map[string]int)
	for i, param := range v.params {
		params[fmt.Sprintf("%v#%v", param.Config().Type(), param.Config().Name())] = i
	}
	for _, param := range ps {
		if param != nil {
			switch v.method {
			case "GET", "HEAD", "DELETE":
				if param.Config().Type() == ParamBody {
					log.Fatalf("Route [%s] %s does not accept any request body parameter [%v]", v.method, v.pat, param.Config().Name())
				}
			}
			switch i, ok := params[fmt.Sprintf("%v#%v", param.Config().Type(), param.Config().Name())]; {
			case ok:
				v.params[i] = param
			default:
				v.params = append(v.params, param)
			}
		}
	}
	return v
}

func (v *router) MaxMemory(m int64) Route {
	v.memory = m
	return v
}

func (v *router) Handle(f HandlerFunc) Route {
	if v.handler != nil {
		log.Fatalf("Route %s can only have one response handler", v.pat)
	}
	v.handler = f
	return v
}

func (v *router) Catch(f HandlerFunc) Route {
	if v.catch != nil {
		log.Fatalf("Route %s can only have one catch handler", v.pat)
	}
	v.catch = f
	return v
}

func (v *router) Middleware(fs ...HandlerFunc) Route {
	for _, f := range fs {
		if f != nil {
			v.middleware = append(v.middleware, f)
		}
	}
	return v
}

func (v *router) Websocket(b bool) Route {
	v.ws = b
	return v
}

func (v *router) HTTP(b bool) Route {
	v.http = b
	return v
}

func (v *router) Models(ms ...Model) Route {
	return v
}

func (v *router) Config() RouteConfig {
	return &config{
		pat:           v.pat,
		aliases:       v.aliases,
		name:          v.name,
		method:        v.method,
		routes:        v.routes,
		namespace:     v.namespace,
		ws:            v.ws,
		http:          v.http,
		documentation: v.documentation,
		summary:       v.summary,
		params:        v.params,
		memory:        v.memory,
		middleware:    v.middleware,
		handler:       v.handler,
		catch:         v.catch,
		models:        v.models,
	}
}
