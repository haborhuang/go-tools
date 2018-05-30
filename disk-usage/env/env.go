package env

import (
	"github.com/caarlos0/env"
	"github.com/haborhuang/go-tools/disk-usage/types"
	"fmt"
)

type envDPConfig struct {
	StatDisks []string `env:"STAT_DISKS" envSeparator:";"`
}

func ParseDisksPathsOrDie(defaultEnv string) types.DisksPaths {
	var config envDPConfig
	env.Parse(&config)

	if len(config.StatDisks) < 1 {
		config.StatDisks = append(config.StatDisks, defaultEnv)
	}

	dp, err := types.SliceToDisksPaths(config.StatDisks, "")
	if nil != err {
		panic(fmt.Errorf("Convert env to disks paths error: %v", err))
	}

	return dp
}

type envMetricsConfig struct {
	MetricsNS     string `env:"METRICS_NS" envDefault:"whispircn_v1"`
	MetricsSubSys string `env:"METRICS_SUBSYSTEM"`
}

func ParseMetricsConfig() types.MetricsConfig {
	var config envMetricsConfig
	env.Parse(&config)

	return types.MetricsConfig{
		Namespace: config.MetricsNS,
		SubSys: config.MetricsSubSys,
	}
}