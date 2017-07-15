package sample

import (
	"fmt"
)

// Namespace represents a group of samples
type Namespace string

// SampleFunc returns samples when invoked or an error
type SampleFunc func() ([]*Sample, error)

//Sample is a mixed-type data structure
//containing values for displaying in toplib.
type Sample struct {
	floats    map[string]float64
	strings   map[string]string
	namespace Namespace
}

func NewSample(id string, namespace Namespace) *Sample {
	sample := &Sample{
		floats:    make(map[string]float64),
		strings:   make(map[string]string),
		namespace: namespace,
	}
	sample.SetString("ID", id)
	return sample
}

func (s *Sample) ID() string           { return s.GetString("ID") }
func (s *Sample) Namespace() Namespace { return s.namespace }

func (s *Sample) SetFloat64(n string, v float64) {
	s.floats[n] = v
}

func (s *Sample) GetFloat64(n string) float64 {
	var value float64
	if v, ok := s.floats[n]; ok {
		value = v
	}
	return value
}

func (s *Sample) SetString(n, v string) {
	s.strings[n] = v
}

func (s *Sample) GetString(n string) string {
	var value string
	if v, ok := s.strings[n]; ok {
		value = v
	}
	return value
}

func (s *Sample) String(n string) string {
	var value string
	if v, ok := s.strings[n]; ok {
		value = v
	}
	if v, ok := s.floats[n]; ok {
		value = fmt.Sprintf("%.2f", v)
	}
	return value
}

func (s *Sample) Strings(fields []string) []string {
	strings := make([]string, len(fields))
	for i, field := range fields {
		strings[i] = s.String(field)
	}
	return strings
}
