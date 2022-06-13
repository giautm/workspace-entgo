package graphql

import (
	"net/http"

	"entgo.io/contrib/entgql"
	"giautm.dev/awesome/ent"
	"giautm.dev/awesome/graphql/directive"
	"giautm.dev/awesome/graphql/generated"
	"giautm.dev/awesome/graphql/resolver"
	"github.com/99designs/gqlgen-contrib/gqlopencensus"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

type config struct {
	apollotracing bool
	introspection bool
	opencensus    bool
	playground    bool

	queryCacheSize          int
	persistedQueryCacheSize int
}

// Option is the GQL server option.
type Option func(*config) error

// WithEntTransaction return the option for enable playground endpoint
func WithPlayground() Option {
	return func(c *config) error {
		c.playground = true
		return nil
	}
}

// WithEnableIntrospection return the option for enable introspection query
func WithEnableIntrospection() Option {
	return func(c *config) error {
		c.introspection = true
		return nil
	}
}

// WithOpencensus return the option for enable Opencensus integration
func WithOpencensus() Option {
	return func(c *config) error {
		c.opencensus = true
		return nil
	}
}

// WithApolloTracing return the option for enable Apollo Tracing
func WithApolloTracing() Option {
	return func(c *config) error {
		c.apollotracing = true
		return nil
	}
}

// WithQueryCache return the option for set the query cache size
func WithQueryCache(size int) Option {
	return func(c *config) error {
		c.queryCacheSize = size
		return nil
	}
}

// WithAutomaticPersistedQuery return the option for set the persisted query cache size
func WithAutomaticPersistedQuery(size int) Option {
	return func(c *config) error {
		c.persistedQueryCacheSize = size
		return nil
	}
}

var (
	// ProductionOptions is options for production mode
	ProductionOptions = []Option{
		WithAutomaticPersistedQuery(1000),
		WithOpencensus(),
		WithQueryCache(1000),
	}
	// DevelopmentOptions is options for development mode
	DevelopmentOptions = append(ProductionOptions,
		WithApolloTracing(),
		WithAutomaticPersistedQuery(0),
		WithEnableIntrospection(),
		WithPlayground(),
		WithQueryCache(0),
	)
)

// Server is the GraphQL server.
type Server struct {
	cfg    *config
	ent    *ent.Client
	schema graphql.ExecutableSchema
}

// NewServer returns a new GraphQL server.
func NewServer(
	resolvers *resolver.Resolver,
	entClient *ent.Client,
	opts ...Option,
) (*Server, error) {
	cfg := &config{
		playground: false,
	}
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: resolvers,
		Directives: generated.DirectiveRoot{
			Auth: directive.Auth,
		},
	})

	return &Server{
		cfg:    cfg,
		ent:    entClient,
		schema: schema,
	}, nil
}

// Routes defines and returns the routes for this server.
func (s *Server) MountRoutes(r chi.Router) {
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = ent.NewContext(ctx, s.ent)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Handle("/query", s.handleQuery())
	if s.cfg.playground {
		r.Method(http.MethodGet, "/playground", playground.Handler("Playground", "/query"))
	}
}

func (s *Server) handleQuery() http.Handler {
	srv := handler.New(s.schema)
	srv.SetRecoverFunc(recoverFunc)
	srv.SetErrorPresenter(errorPresenter)
	// srv.AddTransport(transport.Websocket{
	// 	KeepAlivePingInterval: 10 * time.Second,
	// 	Upgrader: websocket.Upgrader{
	// 		CheckOrigin: func(r *http.Request) bool {
	// 			return true
	// 		},
	// 	},
	// 	InitFunc: WebsocketInit,
	// })
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	if s.cfg.introspection {
		srv.Use(extension.Introspection{})
	}

	if s.cfg.queryCacheSize > 0 {
		srv.SetQueryCache(lru.New(s.cfg.queryCacheSize))
	}
	if s.cfg.persistedQueryCacheSize > 0 {
		srv.Use(extension.AutomaticPersistedQuery{
			Cache: lru.New(s.cfg.persistedQueryCacheSize),
		})
	}
	if s.cfg.opencensus {
		srv.Use(gqlopencensus.Tracer{})
	}
	if s.cfg.apollotracing {
		srv.Use(apollotracing.Tracer{})
	}

	srv.Use(entgql.Transactioner{
		TxOpener: entgql.TxOpenerFunc(ent.OpenTxFromContext),
	})
	return srv
}
