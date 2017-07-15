// +build ignore

package main

import (
	"fmt"
	"github.com/vektorlab/toplib"
	"github.com/vektorlab/toplib/sample"
	"github.com/vektorlab/toplib/section"
	"math/rand"
	"os"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type MockSource struct {
	samples   []*sample.Sample
	namespace sample.Namespace
}

func (m *MockSource) Collect() ([]*sample.Sample, error) {
	samples := []*sample.Sample{}
	for _, seed := range m.samples {
		// Uncomment to change iteration size
		//if rand.Intn(10) > 5 {
		//	continue
		//}
		s := sample.NewSample(seed.ID(), m.namespace)
		s.SetFloat64("CPU", float64(rand.Intn(5)))
		s.SetFloat64("MEM", float64(rand.Intn(40)))
		s.SetFloat64("DISK", float64(rand.Intn(60)))
		s.SetFloat64("GPU", float64(rand.Intn(90)))
		s.SetString("THING", RandString(50))
		s.SetString("OTHER THING", RandString(20))
		samples = append(samples, s)
	}
	return samples, nil
}

func main() {
	sources := []*MockSource{
		&MockSource{
			namespace: sample.Namespace("ns-0"),
			samples:   []*sample.Sample{},
		},
		&MockSource{
			namespace: sample.Namespace("ns-1"),
			samples:   []*sample.Sample{},
		},
		&MockSource{
			namespace: sample.Namespace("ns-2"),
			samples:   []*sample.Sample{},
		},
	}
	for i := 0; i < 35; i++ {
		sources[0].samples = append(sources[0].samples, sample.NewSample(RandString(10), sources[0].namespace))
		sources[1].samples = append(sources[1].samples, sample.NewSample(RandString(10), sources[1].namespace))
		sources[2].samples = append(sources[2].samples, sample.NewSample(RandString(10), sources[2].namespace))
	}
	sections := []toplib.Section{
		section.NewSamples(sources[0].namespace, "ID", "CPU", "MEM", "DISK", "GPU"),
		section.NewSamples(sources[1].namespace, "ID", "CPU", "MEM", "DISK", "GPU", "THING"),
		section.NewSamples(sources[2].namespace, "ID", "CPU", "MEM", "DISK", "GPU", "THING", "OTHER THING"),
		&section.Debug{
			Namespaces: []sample.Namespace{
				sources[0].namespace,
				sources[1].namespace,
				sources[2].namespace,
			},
		},
	}
	if err := toplib.Run(toplib.NewTop(sections), sources[0].Collect, sources[1].Collect, sources[2].Collect); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
