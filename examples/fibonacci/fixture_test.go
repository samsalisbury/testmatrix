package fibonacci

import (
	"testing"

	"github.com/samsalisbury/testmatrix"
)

// fixtureConfig defines our desired initial conditions.
type fixtureConfig struct {
	PrewarmValues []int
}

func (fc *fixtureConfig) Init(t *testing.T, scenario testmatrix.Scenario) *fixture {
	s := unwrapScenario(scenario)
	for _, n := range fc.PrewarmValues {
		s.Fib(n)
	}
	return &fixture{
		s.fibProvider,
	}
}

// fixture represents a fully realised fixture, generated from the fixtureConfig
// and injected scenario.
type fixture struct {
	fibProvider
}

// defaultFixture returns the default fixture config.
func defaultFixture() *fixtureConfig {
	return &fixtureConfig{}
}

func prewarmed(n ...int) *fixtureConfig {
	return &fixtureConfig{PrewarmValues: n}
}
