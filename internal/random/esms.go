package random

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

const nGauss = 1000

type EsmsRandomiser struct {
	rnd          *rand.Rand
	gaussianVars []float64
}

var _ Randomiser = &EsmsRandomiser{}

func NewEsmsRandomiser() *EsmsRandomiser {
	return NewEsmsRandomiserWithSeed(time.Now().UnixMicro())
}

func NewEsmsRandomiserWithSeed(seed int64) *EsmsRandomiser {
	rnd := rand.New(rand.NewSource(seed))
	r := &EsmsRandomiser{rnd: rnd}
	r.fillGaussianVarsArr(nGauss)
	return r
}

// Return a pseudo-random integer uniformly distributed
// between 0 and max
func (r *EsmsRandomiser) UniformRandom(max int) int {
	return r.rnd.Intn(max + 1)
}

func (r *EsmsRandomiser) AveragedRandomPartDev(average, div int) int {
	return r.AveragedRandom(average, average/div)
}

func (r *EsmsRandomiser) AveragedRandom(average int, maxDeviation int) int {
	randGaussian := r.gaussianVars[r.UniformRandom(nGauss-1)]
	deviation := float64(maxDeviation) * randGaussian

	return average + int(deviation)
}

// Given a string with comma separated values (like "a,cd,k")
// returns a random value.
func (r *EsmsRandomiser) RandElem(csv string) string {
	elems := strings.Split(csv, ",")
	return elems[r.UniformRandom(len(elems)-1)]
}

// Throws a bet with probability prob of success. Returns
// true iff succeeded.
func (r *EsmsRandomiser) ThrowWithProb(prob int) bool {
	aThrow := 1 + r.UniformRandom(99)
	return prob >= aThrow
}

func (r *EsmsRandomiser) fillGaussianVarsArr(amount uint) {
	r.gaussianVars = make([]float64, amount)

	for i := uint(0); i < amount; i++ {
		var (
			s  float64
			v1 float64
			v2 float64
			x  float64
		)

		for {
			for {
				u1 := r.rnd.Float64()
				u2 := r.rnd.Float64()

				v1 = 2*u1 - 1
				v2 = 2*u2 - 1

				s = v1*v1 + v2*v2

				if s < 1.0 {
					break
				}
			}

			x = v1 * math.Sqrt(-2*math.Log(s)/s)

			if !(x >= 1.0 || x <= -1.0) {
				break
			}
		}

		r.gaussianVars[i] = x
	}
}
