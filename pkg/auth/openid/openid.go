package openid

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/backoff/v2"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/rancher/opni/pkg/auth"
	"github.com/rancher/opni/pkg/config/v1beta1"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/rbac"
	"github.com/rancher/opni/pkg/util"
	"golang.org/x/sync/singleflight"
)

var (
	ErrNoSigningKeyFound = fmt.Errorf("no signing key found in the JWK set")
	sfGroup              singleflight.Group
)

const (
	TokenKey = "token"
)

type OpenidMiddleware struct {
	keyRefresher *jwk.AutoRefresh
	conf         *OpenidConfig
	logger       *slog.Logger

	wellKnownConfig *WellKnownConfiguration
	lock            sync.Mutex

	cache    *UserInfoCache
	configId string
}

var _ auth.Middleware = (*OpenidMiddleware)(nil)

func New(ctx context.Context, config v1beta1.AuthProviderSpec) (*OpenidMiddleware, error) {
	conf, err := util.DecodeStruct[OpenidConfig](config.Options)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(util.Must(json.Marshal(conf)))

	m := &OpenidMiddleware{
		keyRefresher: jwk.NewAutoRefresh(ctx),
		conf:         conf,
		logger:       logger.New().WithGroup("openid"),
		configId:     string(sum[:]),
	}

	if m.conf.IdentifyingClaim == "" {
		m.conf.IdentifyingClaim = "sub"
	}

	go m.tryConfigureKeyRefresher(ctx)
	return m, nil
}

func (m *OpenidMiddleware) Handle(c *gin.Context) {
	lg := m.logger
	m.lock.Lock()
	if m.wellKnownConfig == nil {
		m.lock.Unlock()
		lg.Debug("error handling request: auth provider is not ready")
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	m.lock.Unlock()

	lg.Debug("handling auth request")
	// Some providers serve their JWKS URI at `/.well-known/jwks.json`, which is
	// not a registered well-known URI. openid-configuration is, however.
	ctx, ca := context.WithTimeout(c.Request.Context(), time.Second*5)
	defer ca()
	set, err := m.keyRefresher.Fetch(ctx, m.wellKnownConfig.JwksUri)
	if err != nil {
		lg.Error("failed to fetch JWK set:", logger.Err(err))
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		lg.Error("no authorization header in request")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	bearerToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
	var userID string
	switch GetTokenType(bearerToken) {
	case IDToken:
		idt, err := ValidateIDToken(bearerToken, set)
		if err != nil {
			lg.Error("failed to validate ID token:", logger.Err(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claim, ok := idt.Get(m.conf.IdentifyingClaim)
		if !ok {
			lg.Error("identifying claim not found in ID token", "claim", m.conf.IdentifyingClaim)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID = fmt.Sprint(claim)
	case Opaque:
		userInfo, err := m.cache.Get(bearerToken)
		if err != nil {
			lg.Error("failed to get user info:", logger.Err(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		uid, err := userInfo.UserID()
		if err != nil {
			lg.Error("failed to get user id:", logger.Err(err))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID = uid
	default:
		lg.Error("could not determine token type")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Header("Authorization", "")
	c.Set(rbac.UserIDKey, userID)
}

func (m *OpenidMiddleware) tryConfigureKeyRefresher(ctx context.Context) {
	lg := m.logger
	result, _, _ := sfGroup.Do(m.configId, func() (interface{}, error) {
		p := backoff.Exponential(
			backoff.WithMaxRetries(0),
			backoff.WithMinInterval(50*time.Millisecond),
			backoff.WithMaxInterval(time.Minute),
			backoff.WithMultiplier(2),
			backoff.WithJitterFactor(0.05),
		)
		b := p.Start(ctx)
		for backoff.Continue(b) {
			wellKnownCfg, err := m.conf.GetWellKnownConfiguration()
			if err != nil {
				if isDiscoveryErrFatal(err) {
					panic("fatal error fetching openid configuration")
				} else {
					lg.Warn("failed to fetch openid configuration (will retry)", logger.Err(err))
				}
				continue
			}
			return wellKnownCfg, nil
		}
		panic("unreachable")
	})

	wellKnownCfg := result.(*WellKnownConfiguration)
	lg.Info("successfully fetched openid configuration", "issuer", wellKnownCfg.Issuer)
	m.lock.Lock()
	defer m.lock.Unlock()
	m.wellKnownConfig = wellKnownCfg
	httpClient := http.DefaultClient
	if m.conf.Discovery != nil && m.conf.Discovery.CACert != nil {
		lg.Info("using custom CA cert for openid discovery", "filename", m.conf.Discovery.CACert)
		certPool := x509.NewCertPool()
		data, err := os.ReadFile(*m.conf.Discovery.CACert)
		if err != nil {
			panic("openid discovery: failed to read CA cert")
		}
		if !certPool.AppendCertsFromPEM(data) {
			panic("openid discovery: invalid ca cert")
		}
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: certPool,
				},
			},
		}
	}
	m.keyRefresher.Configure(wellKnownCfg.JwksUri,
		jwk.WithHTTPClient(httpClient),
	)
	var err error
	m.cache, err = NewUserInfoCache(m.conf, m.logger, WithHTTPClient(httpClient))
	if err != nil {
		panic("failed to create user info cache")
	}
}
