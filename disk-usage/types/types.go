package types

import (
	"fmt"
	"strings"
)

func SliceToDisksPaths(s []string, delimiter string) (DisksPaths, error) {
	if len(s) < 1 {
		return nil, nil
	}

	if delimiter == "" {
		delimiter = ":"
	}

	ps := make(DisksPaths)
	for _, kv := range s {
		dp := strings.Split(kv, delimiter)
		if len(dp) != 2 {
			return nil, fmt.Errorf("Cannot parse '%s'", kv)
		}

		ps[dp[0]] = dp[1]
	}

	return ps, nil
}

type DisksPaths map[string]string

type MetricsConfig struct {
	Namespace string
	SubSys string
}