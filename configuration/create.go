package configuration

import (
	"github.com/RKO-solver/rko-go/metaheuristc/ga"
	"github.com/RKO-solver/rko-go/metaheuristc/ils"
	"github.com/RKO-solver/rko-go/metaheuristc/lns"
	"github.com/RKO-solver/rko-go/metaheuristc/multistart"
	"github.com/RKO-solver/rko-go/metaheuristc/sa"
	"github.com/RKO-solver/rko-go/metaheuristc/vns"
)

func DefaultConfiguration() *MetaheuristicsConfiguration {
	return &MetaheuristicsConfiguration{
		MultiStart: multistart.DefaulConfigurationtMultiStart(),
		BRKGA:      ga.DefaultConfigurationBRKGA(),
		GA:         ga.DefaultConfigurationGA(),
		ILS:        ils.DefaultConfigurationILS(),
		SA:         sa.DefaultConfigurationSA(),
		VNS:        vns.DefaultConfigurationVNS(),
		LNS:        lns.DefaultConfigurationVNS(),
	}
}
