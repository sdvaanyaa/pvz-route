package packaging

import (
	"fmt"
)

type WeightLimitError struct {
	Weight  float64
	Max     float64
	Package string
}

func (e *WeightLimitError) Error() string {
	return fmt.Sprintf(
		"invalid weight %.2fkg for %s: must be less than %.2fkg",
		e.Weight,
		e.Package,
		e.Max,
	)
}
