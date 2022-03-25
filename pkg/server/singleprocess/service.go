package singleprocess

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"

	wphznpb "github.com/hashicorp/waypoint-hzn/pkg/pb"
	"github.com/hashicorp/waypoint/internal/server/boltdbstate"
	"github.com/hashicorp/waypoint/internal/serverconfig"
	wpoidc "github.com/hashicorp/waypoint/pkg/auth/oidc"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/server/logstream"
	"github.com/hashicorp/waypoint/pkg/serverstate"
)

// Service implements the gRPC service for the server.
type Service struct {
	pb.UnimplementedWaypointServer

	// state is the state management interface that provides functions for
	// safely mutating server state.
	state func(ctx context.Context) serverstate.Interface

	// id is our unique server ID.
	id string

	// encodeId uses the provided context to encode additional metadata
	// (if present), and returns an ID that can be decoded by DecodeId.
	encodeId func(ctx context.Context, id string) (encodedId string, err error)

	// decodeId takes a string that contains an ID (likely created
	// with EncodeId), and returns only the waypoint-relevant ID.
	decodeId func(encodedId string) (id string, err error)

	logStreamProvider logstream.Provider

	// urlConfig is not nil if the URL service is enabled. This is guaranteed
	// to have the configs set.
	urlConfig    *serverconfig.URL
	urlClientMu  sync.Mutex
	urlClientVal wphznpb.WaypointHznClient

	// urlCEB has the configuration for the entrypoint. If this is nil,
	// then the configuration is not ready. The urlCEBWatchCh can be used
	// to watch for changes. All fields protected with urlCEBMu.
	urlCEBMu      sync.RWMutex
	urlCEB        *pb.EntrypointConfig_URLService
	urlCEBWatchCh chan struct{}

	// bgCtx is used for background tasks within the service. This is
	// cancelled when Close is called.
	bgCtx       context.Context
	bgCtxCancel context.CancelFunc

	// bgWg is incremented for every background goroutine that the
	// service starts up. When Close is called, we wait on this to ensure
	// that we fully shut down before returning.
	bgWg sync.WaitGroup

	// superuser is true if all API actions should act as if a superuser
	// made them. This is used for local mode only.
	superuser bool

	// oidcCache is the cache for OIDC providers.
	oidcCache *wpoidc.ProviderCache
}

// New returns a Waypoint server implementation that uses BotlDB plus
// in-memory locks to operate safely.
func New(opts ...Option) (pb.WaypointServer, error) {
	var s Service

	// Set default config values
	cfg := config{
		idEncoder: func(ctx context.Context, id string) (encodedId string, err error) {
			return id, nil
		},
		idDecoder: func(encodedId string) (id string, err error) {
			return encodedId, nil
		},
		logStreamProvider: &singleProcessLogStreamProvider{},
	}

	for _, opt := range opts {
		if err := opt(&s, &cfg); err != nil {
			return nil, err
		}
	}

	s.encodeId = cfg.idEncoder
	s.decodeId = cfg.idDecoder

	if !cfg.oidcDisabled {
		s.oidcCache = wpoidc.NewProviderCache()
	}

	if cfg.logStreamProvider != nil {
		s.logStreamProvider = cfg.logStreamProvider
	}

	log := cfg.log
	if log == nil {
		log = hclog.L()
	}

	s.state = cfg.stateProvider

	// If we don't have a server ID, set that.

	// TODO(izaak): the serverstate interface doesn't currently support
	// the ServerId methods, but I think it probably needs to.
	state := s.state(context.Background())
	if boltdbstate, ok := state.(*boltdbstate.State); ok {
		id, err := boltdbstate.ServerIdGet()
		if err != nil {
			return nil, err
		}
		if id == "" {
			id, err = server.Id()
			if err != nil {
				return nil, err
			}

			if err := boltdbstate.ServerIdSet(id); err != nil {
				return nil, err
			}
		}
		s.id = id
	}

	// Setup our URL service config if it is enabled.
	if scfg := cfg.serverConfig; scfg != nil && scfg.URL != nil {
		// Set our config
		s.urlConfig = scfg.URL

		// Create a copy of the config that we use for initialization so
		// that we don't create races with s.urlConfig if this retries.
		cfgCopy := *scfg.URL

		// Initialize our CEB settings.
		s.urlCEBMu.Lock()
		s.urlCEB = &pb.EntrypointConfig_URLService{
			ControlAddr: cfgCopy.ControlAddress,
			Token:       cfgCopy.APIToken,
		}
		s.urlCEBWatchCh = make(chan struct{})
		s.urlCEBMu.Unlock()

		if err := s.initURLClient(
			log.Named("url_service"),
			nil,
			cfg.acceptUrlTerms,
			&cfgCopy,
		); err != nil {
			return nil, err
		}
	}

	// If we haven't initialized our server config before, do that once.
	conf, err := state.ServerConfigGet()
	if err != nil {
		return nil, err
	}
	if conf.Cookie == "" {
		if err := state.ServerConfigSet(conf); err != nil {
			return nil, err
		}
	}

	// Set specific server config for the deployment entrypoint binaries
	if scfg := cfg.serverConfig; scfg != nil && scfg.CEBConfig != nil && scfg.CEBConfig.Addr != "" {
		// only one advertise address can be configured
		addr := &pb.ServerConfig_AdvertiseAddr{
			Addr:          scfg.CEBConfig.Addr,
			Tls:           scfg.CEBConfig.TLSEnabled,
			TlsSkipVerify: scfg.CEBConfig.TLSSkipVerify,
		}

		conf := &pb.ServerConfig{
			AdvertiseAddrs: []*pb.ServerConfig_AdvertiseAddr{addr},
		}

		if err := state.ServerConfigSet(conf); err != nil {
			return nil, err
		}
	}

	// Setup background polling if not disabled
	if !cfg.pollingDisabled {
		// Setup the background context that is used for internal tasks
		s.bgCtx, s.bgCtxCancel = context.WithCancel(context.Background())

		// TODO: When more items are added, move this else where
		// pollableItems is a map of potential items Waypoint can queue a poll for.
		// Each item should implement the pollHandler interface
		pollableItems := map[string]pollHandler{
			"project":                  &projectPoll{state: state},
			"application_statusreport": &applicationPoll{state: state, workspace: "default"},
		}

		// Start our polling background goroutines.
		// See the func docs for more info.
		for pollName, pollItem := range pollableItems {
			s.bgWg.Add(1)
			go s.runPollQueuer(
				s.bgCtx, &s.bgWg, pollItem,
				log.Named("poll_queuer").Named(pollName),
			)
		}

		// Start out state pruning background goroutine. This calls
		// Prune on the state every 10 minutes.
		s.bgWg.Add(1)
		go s.runPrune(s.bgCtx, &s.bgWg, log.Named("prune"))
	}

	return &s, nil
}

// Close shuts down any background processes and resources that may
// be used by the service. This should be called after the service
// is no longer responding to requests.
func (s *Service) Close() error {
	if s.bgCtxCancel != nil {
		s.bgCtxCancel()
		s.bgWg.Wait()
	}

	if s.oidcCache != nil {
		s.oidcCache.Close()
	}
	return nil
}

type config struct {
	stateProvider func(ctx context.Context) serverstate.Interface

	idEncoder func(ctx context.Context, id string) (encodedId string, err error)
	idDecoder func(encodedId string) (id string, err error)

	serverConfig      *serverconfig.Config
	log               hclog.Logger
	superuser         bool
	oidcDisabled      bool
	pollingDisabled   bool
	logStreamProvider logstream.Provider

	acceptUrlTerms bool
}

type Option func(*Service, *config) error

// WithState sets the server to use the given state package
func WithState(state serverstate.Interface) Option {
	return func(s *Service, cfg *config) error {
		cfg.stateProvider = func(ctx context.Context) serverstate.Interface {
			return state
		}
		return nil
	}
}

// WithStateProvider sets the server to use the state provided by
// the given function
func WithStateProvider(stateProvider func(ctx context.Context) serverstate.Interface) Option {
	return func(s *Service, cfg *config) error {
		cfg.stateProvider = stateProvider
		return nil
	}
}

// WithConfig sets the server config in use with this server.
func WithConfig(scfg *serverconfig.Config) Option {
	return func(s *Service, cfg *config) error {
		cfg.serverConfig = scfg
		return nil
	}
}

// WithLogger sets the logger for use with the server.
func WithLogger(log hclog.Logger) Option {
	return func(s *Service, cfg *config) error {
		cfg.log = log
		return nil
	}
}

// WithSuperuser forces all API actions to behave as if a superuser
// made them. This is usually turned on for local mode only. There is no
// option (at the time of writing) to enable this on a network-attached server.
func WithSuperuser() Option {
	return func(s *Service, cfg *config) error {
		s.superuser = true
		return nil
	}
}

// WithAcceptURLTerms will set the config to either accept or reject the terms
// of service for using the URL service. Rejecting the TOS will disable the
// URL service. Note that the actual rejection does not occur until the
// waypoint horizon client attempts to register its guest account.
func WithAcceptURLTerms(accept bool) Option {
	return func(s *Service, cfg *config) error {
		cfg.acceptUrlTerms = accept
		return nil
	}
}

// WithIdEncoder provides the server with an optional function that
// takes a waypoint ID (user id, runner id, etc.),
// uses the provided context to encode additional metadata
// (if present), and returns an ID that can be decoded by DecodeId.
func WithIdEncoder(encoder func(ctx context.Context, id string) (encodedId string, err error)) Option {
	return func(s *Service, cfg *config) error {
		cfg.idEncoder = encoder
		return nil
	}
}

// WithIdDecoder provides the server with an optional function that
// takes a string that contains an ID (likely created
// with EncodeId), and returns only the waypoint-relevant ID.
func WithIdDecoder(decoder func(encodedId string) (id string, err error)) Option {
	return func(s *Service, cfg *config) error {
		cfg.idDecoder = decoder
		return nil
	}
}

func WithOidcDisabled(disabled bool) Option {
	return func(s *Service, cfg *config) error {
		cfg.oidcDisabled = disabled
		return nil
	}
}

func WithPollingDisabled(disabled bool) Option {
	return func(s *Service, cfg *config) error {
		cfg.pollingDisabled = disabled
		return nil
	}
}

func WithLogStreamProvider(logStreamProvider logstream.Provider) Option {
	return func(s *Service, cfg *config) error {
		cfg.logStreamProvider = logStreamProvider
		return nil
	}
}

var _ pb.WaypointServer = (*Service)(nil)
