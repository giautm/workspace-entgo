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
	entClient     *ent.Client
	apollotracing bool
	introspection bool
	opencensus    bool
	playground    bool

	queryCacheSize          int
	persistedQueryCacheSize int
}

type Option func(*config) error

func WithPlayground() Option {
	return func(c *config) error {
		c.playground = true
		return nil
	}
}

func WithEntTransaction(client *ent.Client) Option {
	return func(c *config) error {
		c.entClient = client
		return nil
	}
}

func WithEnableIntrospection() Option {
	return func(c *config) error {
		c.introspection = true
		return nil
	}
}

func WithOpencensus() Option {
	return func(c *config) error {
		c.opencensus = true
		return nil
	}
}

func WithApolloTracing() Option {
	return func(c *config) error {
		c.apollotracing = true
		return nil
	}
}

func WithQueryCache(size int) Option {
	return func(c *config) error {
		c.queryCacheSize = size
		return nil
	}
}

func WithAutomaticPersistedQuery(size int) Option {
	return func(c *config) error {
		c.persistedQueryCacheSize = size
		return nil
	}
}

var (
	ProductionOptions = []Option{
		WithAutomaticPersistedQuery(1000),
		WithOpencensus(),
		WithQueryCache(1000),
	}
	DevelopmentOptions = append(ProductionOptions,
		WithApolloTracing(),
		WithAutomaticPersistedQuery(0),
		WithEnableIntrospection(),
		WithPlayground(),
		WithQueryCache(0),
	)
)

type Server struct {
	cfg    *config
	schema graphql.ExecutableSchema
}

func NewServer(
	resolvers *resolver.Resolver,
	opts ...Option,
) (*Server, error) {
	cfg := &config{
		entClient:  nil,
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
		schema: schema,
	}, nil
}

// Routes defines and returns the routes for this server.
func (s *Server) MountRoutes(r chi.Router) {
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
	if s.cfg.entClient != nil {
		srv.Use(entgql.Transactioner{
			TxOpener: s.cfg.entClient,
		})
	}

	return srv
}
