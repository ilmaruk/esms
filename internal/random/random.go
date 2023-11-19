package random

type Randomiser interface {
	// Return a pseudo-random integer uniformly distributed
	// between 0 and max
	UniformRandom(max int) int

	AveragedRandomPartDev(average, div int) int

	AveragedRandom(average int, maxDeviation int) int

	// Given a string with comma separated values (like "a,cd,k")
	// returns a random value.
	RandElem(csv string) string

	// Throws a bet with probability prob of success. Returns
	// true iff succeeded.
	ThrowWithProb(prob int) bool
}
