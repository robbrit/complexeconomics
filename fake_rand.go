package econerra

// fakeRand is an implementation of rand.Source that produces a fixed sequence
// of integers.
type fakeRand struct {
	i   int
	seq []int64
}

func NewFakeRand(seq []int64) *fakeRand {
	return &fakeRand{0, seq}
}

func (fr *fakeRand) Int63() int64 {
	i := fr.i
	fr.i++
	return fr.seq[i%len(fr.seq)]
}

func (fr *fakeRand) Seed(seed int64) {}
