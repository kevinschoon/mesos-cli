package sample

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortSamples(t *testing.T) {
	samples := []*Sample{
		NewSample("S1", Namespace("")),
		NewSample("S2", Namespace("")),
		NewSample("S3", Namespace("")),
	}
	samples[0].SetFloat64("flt", float64(1.0))
	samples[1].SetFloat64("flt", float64(2.0))
	samples[2].SetFloat64("flt", float64(3.0))
	samples[0].SetString("str", "aaa")
	samples[1].SetString("str", "bbb")
	samples[2].SetString("str", "ccc")
	Sort("flt", samples)
	assert.Equal(t, float64(3.0), samples[0].GetFloat64("flt"))
	assert.Equal(t, float64(2.0), samples[1].GetFloat64("flt"))
	assert.Equal(t, float64(1.0), samples[2].GetFloat64("flt"))
	Sort("str", samples)
	assert.Equal(t, "aaa", samples[0].GetString("str"))
	assert.Equal(t, "bbb", samples[1].GetString("str"))
	assert.Equal(t, "ccc", samples[2].GetString("str"))
}
