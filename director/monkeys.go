package director

import (
	"fmt"
	"github.com/eelavai/octopus/state"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// We use reflection to look up monkeys and avoid writing a bunch of boilerplate
func (d *Director) monkey(name string) func(*rand.Rand, float64) {
	v := reflect.ValueOf(d)
	name = strings.Title(name)
	method := v.MethodByName(fmt.Sprintf("%sMonkey", name))
	return func(rng *rand.Rand, intensity float64) {
		args := []reflect.Value{
			reflect.ValueOf(rng),
			reflect.ValueOf(intensity),
		}
		method.Call(args)
	}
}

func (d *Director) spawn(name string) {
	rng := state.NewRand("monkey " + name)

	c := d.config[name]
	monkey := d.monkey(name)

	if c.frequency == 0 {
		return
	}

	time.Sleep(c.offset)
	start := time.Now()

	for {
		time.Sleep(time.Duration(rng.ExpFloat64() * float64(c.frequency)))
		dt := time.Now().Sub(start)
		if c.ramp == 0 {
			monkey(rng, 0)
		} else {
			monkey(rng, float64(dt)/float64(c.ramp))
		}
	}
}

func (d *Director) LatencyMonkey(rng *rand.Rand, intensity float64) {
	target := d.randomLink(rng)
	latency := d.makeLatency(rng, intensity)
	log.Printf("[monkey] Setting latency for %v to %v", target, latency)
	target.SetLatency(latency)
}
func (d *Director) JitterMonkey(rng *rand.Rand, intensity float64) {
	target := d.randomLink(rng)
	jitter := d.makeJitter(rng, intensity)
	log.Printf("[monkey] Setting jitter for %v to %v", target, jitter)
	target.SetJitter(jitter)
}
func (d *Director) LagsplitMonkey(rng *rand.Rand, intensity float64) {
	targets := d.randomPartition(rng)
	// TODO: different latency function?
	latency := d.makeLatency(rng, intensity)
	duration := d.makeDuration(rng, 1000, intensity)
	for _, target := range targets {
		go target.Lag(latency, duration)
	}
}
func (d *Director) LinkMonkey(rng *rand.Rand, intensity float64) {
	target := d.randomLink(rng)
	duration := d.makeDuration(rng, 1000, intensity)
	log.Printf("[monkey] Killing %v for %v", target, duration)
	go target.Kill(duration)
}
func (d *Director) NetsplitMonkey(rng *rand.Rand, intensity float64) {
	targets := d.randomPartition(rng)
	duration := d.makeDuration(rng, 1000, intensity)
	log.Printf("[monkey] Killing links in partition %v for %v", targets, duration)
	for _, target := range targets {
		go target.Kill(duration)
	}
}
func (d *Director) FreezeMonkey(rng *rand.Rand, intensity float64) {
	target := d.randomAgent(rng)
	duration := d.makeDuration(rng, 1000, intensity)
	log.Printf("[monkey] Freezing %v for %v", target, duration)
	go target.Stop(duration)
}
func (d *Director) MurderMonkey(rng *rand.Rand, intensity float64) {
	target := d.randomAgent(rng)
	duration := d.makeDuration(rng, 1000, intensity)
	log.Printf("[monkey] Murdering %v for %v", target, duration)
	go target.Kill(duration)
}
