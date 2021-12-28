package scaler

const (
	limitRate    = 10
	puUpperLimit = 1000
	puLowerLimit = 100
	ceilBase     = 10
)

type PUCalculator struct {
	dbCount int
	buffer  int
}

func NewPUCalculator(dbCount, buffer int) *PUCalculator {
	return &PUCalculator{
		dbCount: dbCount,
		buffer:  buffer,
	}
}

// calc desired pu by dbCount
// for example.
// buffer=5
// dbcount -> PU
//   0 -> 100
//   4 -> 100
//   5 -> 200
//  14 -> 200
//  15 -> 300
//  84 -> 900
//  85 -> 1000
// 100 -> 1000
func (c *PUCalculator) DesiredPU() int {
	ceiledPUWithBuf := ceil(c.dbCount+c.buffer) * limitRate
	return min(ceiledPUWithBuf+puLowerLimit, puUpperLimit)
}

func (c *PUCalculator) IsUpperLimit(pu int) bool {
	return pu == puUpperLimit
}

func (c *PUCalculator) IsLowerLimit(pu int) bool {
	return pu == puLowerLimit
}

func ceil(x int) int {
	return (x / ceilBase) * ceilBase
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
