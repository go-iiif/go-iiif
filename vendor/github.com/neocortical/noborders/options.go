package noborders

// DefaultEntropyThreshold is the default value of the entropy threshold.
const DefaultEntropyThreshold = 0.08

// DefaultVarianceThreshold is the default value of the variance threshold.
const DefaultVarianceThreshold = 10000000

// Options defines custom options
type Options interface {
	SetMultiPass(bool) Options
	SetEntropy(float64) Options
	SetVariance(float64) Options
	MultiPass() bool
	Entropy() float64
	Variance() float64
}

type options struct {
	multiplePasses    bool
	entropyThreshold  float64
	varianceThreshold float64
}

// Opts initializes a new Options object with defaults.
func Opts() Options {
	return &options{
		multiplePasses:    false,
		entropyThreshold:  DefaultEntropyThreshold,
		varianceThreshold: DefaultVarianceThreshold,
	}
}

// SetMultiPass controls whether the algorithm will be applied multiple times.
// This is useful for images that have been screenshotted multiple times but
// may result in increased loss of empty background space. The default is false.
// If true, the algorithm will be repeated until quiescent.
func (o *options) SetMultiPass(v bool) Options {
	o.multiplePasses = v
	return o
}

// SetEntropy sets the entropy threshold.
func (o *options) SetEntropy(v float64) Options {
	o.entropyThreshold = v
	return o
}

// SetVariance sets the variance threshold.
func (o *options) SetVariance(v float64) Options {
	o.varianceThreshold = v
	return o
}

// MultiPass returns the multiple passes setting.
func (o *options) MultiPass() bool {
	return o.multiplePasses
}

// Entropy returns the entropy setting.
func (o *options) Entropy() float64 {
	return o.entropyThreshold
}

// Variance returns the variance setting.
func (o *options) Variance() float64 {
	return o.varianceThreshold
}
