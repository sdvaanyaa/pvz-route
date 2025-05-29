package packaging

const (
	PackageNone    = "none"
	PackageBag     = "bag"
	PackageBox     = "box"
	PackageFilm    = "film"
	PackageBagFilm = "bag+film"
	PackageBoxFilm = "box+film"
)

type PackageStrategy interface {
	ValidateWeight(weight float64) error
	CalculatePrice(basePrice float64) float64
}
