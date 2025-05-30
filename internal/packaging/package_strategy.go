package packaging

const (
	PackageNone    = "none"
	PackageBag     = "bag"
	PackageBox     = "box"
	PackageFilm    = "film"
	PackageBagFilm = "bag+film"
	PackageBoxFilm = "box+film"
)

// PackageStrategy defines the behavior for different packaging strategies.
// It includes methods to validate if a package's weight is acceptable
// and to calculate the final price based on a base price.
type PackageStrategy interface {
	// ValidateWeight checks if the given weight complies with the packaging rules.
	// Returns an error if the weight is invalid.
	ValidateWeight(weight float64) error

	// CalculatePrice returns the final price for the package based on the base price,
	// adjusted according to the packaging strategy.
	CalculatePrice(basePrice float64) float64
}
