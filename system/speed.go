package system

// Speed is a type that specifies one of the defined system update rates.
type Speed int

// Slow, Normal and Fast enumerate the three possible system update rates, while
// NumSpeeds tracks how many there are.
const (
	Slow Speed = iota
	Normal
	Fast
	NumSpeeds int = iota
)
