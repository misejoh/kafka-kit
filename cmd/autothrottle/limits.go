package main

import (
	"errors"
	"math"

	"github.com/DataDog/topicmappr/kafkametrics"
)

// Limits is a map of instance-type
// to network bandwidth limits.
type Limits map[string]float64

// NewLimitsConfig is used to initialize
// a Limits.
type NewLimitsConfig struct {
	// Min throttle rate in MB/s.
	Minimum float64
	// Max throttle rate as a portion of capacity.
	Maximum float64
	// Map of instance-type to total network capacity in MB/s.
	CapacityMap map[string]float64
}

// NewLimits takes a minimum float64 and a map
// of instance-type to float64 network capacity values (in MB/s).
func NewLimits(c NewLimitsConfig) Limits {
	lim := Limits{
		// Min. config.
		"minimum": c.Minimum,
		"maximum": c.Maximum,
	}

	// Update with provided
	// capacity map.
	for k, v := range c.CapacityMap {
		lim[k] = v
	}

	return lim
}

// headroom takes a *kafkametrics.Broker and last set
// throttle rate and returns the headroom based on utilization
// vs capacity. Headroom is determined by subtracting the current
// throttle rate from the current outbound network utilization.
// This yields a crude approximation of how much non-replication
// throughput is currently being demanded. The non-replication
// throughput is then subtracted from the total network capacity available.
// This value suggests what headroom is available for replication.
// We then use the greater of:
// - this value * the configured portion of free bandwidth eligible for replication
// - the configured minimum replication rate in MB/s
func (l Limits) headroom(b *kafkametrics.Broker, t float64) (float64, error) {
	if b == nil {
		return l["minimum"], errors.New("Nil broker provided")
	}

	if k, exists := l[b.InstanceType]; exists {
		nonThrottleUtil := math.Max(b.NetTX-t, 0.00)
		return math.Max((k-nonThrottleUtil)*(l["maximum"]/100), l["minimum"]), nil
	}

	return l["minimum"], errors.New("Unknown instance type")
}