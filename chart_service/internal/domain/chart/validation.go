package chart

var allowedPositionCounts = map[int]struct{}{
	5:  {},
	10: {},
	20: {},
	25: {},
	30: {},
	40: {},
	50: {},
}

func IsValidPositionCount(v int) bool {
	_, ok := allowedPositionCounts[v]
	return ok
}
