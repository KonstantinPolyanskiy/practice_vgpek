package logger

import "go.uber.org/zap"

type Config struct {
	Level         string
	HasCaller     bool
	HasStacktrace bool
	Encoding      string
}

func New(cfg Config) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return zap.NewProduction()
	}

	zapCfg := zap.NewProductionConfig()

	zapCfg.Level = lvl
	zapCfg.Encoding = cfg.Encoding
	zapCfg.DisableCaller = !cfg.HasCaller
	zapCfg.DisableStacktrace = !cfg.HasStacktrace

	return zapCfg.Build()
}
