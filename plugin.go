package redis

import (
	"github.com/go-redis/redis/v8"
	endure "github.com/roadrunner-server/endure/pkg/container"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	"go.uber.org/zap"
)

const PluginName = "redis"

type Plugin struct {
	cfg config.Configurer
	log *zap.Logger
}

func (p *Plugin) Init(cfg config.Configurer, log *zap.Logger) error {
	p.cfg = cfg
	p.log = log
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

// Provides declares factory methods.
func (p *Plugin) Provides() []any {
	return []any{
		p.ServiceUniversalClient,
	}
}

func (p *Plugin) ServiceUniversalClient(n endure.Named) (redis.UniversalClient, error) {
	const op = errors.Op("redis plugin service universal client")

	key := n.Name() + "." + PluginName

	var cfg Config

	if p.cfg.Has(key) {
		if err := p.cfg.UnmarshalKey(key, &cfg); err != nil {
			return nil, errors.E(op, errors.Providers, err)
		}
	}

	p.log.Debug("collect redis config", zap.String("key", key), zap.Any("config", cfg))

	cfg.InitDefaults()

	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:              cfg.Addrs,
		DB:                 cfg.DB,
		Username:           cfg.Username,
		Password:           cfg.Password,
		SentinelPassword:   cfg.SentinelPassword,
		MaxRetries:         cfg.MaxRetries,
		MinRetryBackoff:    cfg.MaxRetryBackoff,
		MaxRetryBackoff:    cfg.MaxRetryBackoff,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		PoolSize:           cfg.PoolSize,
		MinIdleConns:       cfg.MinIdleConns,
		MaxConnAge:         cfg.MaxConnAge,
		PoolTimeout:        cfg.PoolTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFreq,
		ReadOnly:           cfg.ReadOnly,
		RouteByLatency:     cfg.RouteByLatency,
		RouteRandomly:      cfg.RouteRandomly,
		MasterName:         cfg.MasterName,
	}), nil
}
