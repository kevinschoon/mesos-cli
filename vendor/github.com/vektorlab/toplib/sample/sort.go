package sample

import "sort"

type By func(*Sample, *Sample) bool

func (by By) Sort(samples []*Sample, reverse bool) {
	s := &sorter{
		samples: samples,
		by:      by,
	}
	if reverse {
		sort.Sort(sort.Reverse(s))
	} else {
		sort.Sort(s)
	}
}

type sorter struct {
	samples []*Sample
	by      func(s1, s2 *Sample) bool
}

func (s *sorter) Len() int           { return len(s.samples) }
func (s *sorter) Swap(i, j int)      { s.samples[i], s.samples[j] = s.samples[j], s.samples[i] }
func (s *sorter) Less(i, j int) bool { return s.by(s.samples[i], s.samples[j]) }

func Sort(n string, samples []*Sample) {
	var (
		fn      By
		reverse bool
	)
	if len(samples) == 0 {
		return
	}
	if _, ok := samples[0].strings[n]; ok {
		fn = func(i, j *Sample) bool {
			return i.GetString(n) < j.GetString(n)
		}
	}
	if _, ok := samples[0].floats[n]; ok {
		reverse = true
		fn = func(i, j *Sample) bool {
			return i.GetFloat64(n) < j.GetFloat64(n)
		}
	}
	By(fn).Sort(samples, reverse)
}
