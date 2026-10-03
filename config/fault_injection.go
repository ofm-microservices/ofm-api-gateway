package config

import "time"

// FaultInjectionConfig controls test-only failures at the gateway HTTP
// boundary. It is deliberately disabled unless explicitly enabled in a
// non-production environment.
type FaultInjectionConfig struct {
	Enabled     bool          `env:"FAULT_INJECTION_ENABLED" envDefault:"false"`
	Profile     string        `env:"FAULT_INJECTION_PROFILE" envDefault:"none"`
	Target      string        `env:"FAULT_INJECTION_TARGET" envDefault:"*"`
	Rate        float64       `env:"FAULT_INJECTION_RATE" envDefault:"0"`
	Delay       time.Duration `env:"FAULT_INJECTION_DELAY" envDefault:"0s"`
	MaxFailures int64         `env:"FAULT_INJECTION_MAX_FAILURES" envDefault:"0"`
	TTL         time.Duration `env:"FAULT_INJECTION_TTL" envDefault:"0s"`
	TestToken   string        `env:"FAULT_INJECTION_TEST_TOKEN" envDefault:""`
}
