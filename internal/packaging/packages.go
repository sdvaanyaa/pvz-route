package packaging

const (
	MaxWeightBag = 10.0
	MaxWeightBox = 30.0

	PriceBag  = 5.0
	PriceBox  = 20.0
	PriceFilm = 1.0
)

// CompositePackage represents a packaging strategy composed of multiple sub-packages.
// It validates and calculates price by delegating to each component.
type CompositePackage struct {
	components []PackageStrategy
}

// NewCompositePackage creates a new CompositePackage from given PackageStrategy components.
func NewCompositePackage(parts ...PackageStrategy) *CompositePackage {
	return &CompositePackage{components: parts}
}

// ValidateWeight checks if the weight is valid for all components.
// Returns an error if any component rejects the weight.
func (c *CompositePackage) ValidateWeight(weight float64) error {
	for _, part := range c.components {
		if err := part.ValidateWeight(weight); err != nil {
			return err
		}
	}
	return nil
}

// CalculatePrice sums up the price adjustments from all components.
func (c *CompositePackage) CalculatePrice(basePrice float64) float64 {
	price := basePrice
	for _, part := range c.components {
		price = part.CalculatePrice(price)
	}
	return price
}

type NonePackage struct{}

func (p *NonePackage) ValidateWeight(_ float64) error {
	return nil
}

func (p *NonePackage) CalculatePrice(basePrice float64) float64 {
	return basePrice
}

type BagPackage struct{}

func (p *BagPackage) ValidateWeight(weight float64) error {
	if weight >= MaxWeightBag {
		return &WeightLimitError{
			Weight:  weight,
			Max:     MaxWeightBag,
			Package: PackageBag,
		}
	}
	return nil
}

func (p *BagPackage) CalculatePrice(basePrice float64) float64 {
	return basePrice + PriceBag
}

type BoxPackage struct{}

func (p *BoxPackage) ValidateWeight(weight float64) error {
	if weight >= MaxWeightBox {
		return &WeightLimitError{
			Weight:  weight,
			Max:     MaxWeightBox,
			Package: PackageBox,
		}
	}
	return nil
}

func (p *BoxPackage) CalculatePrice(basePrice float64) float64 {
	return basePrice + PriceBox
}

type FilmPackage struct{}

func (p *FilmPackage) ValidateWeight(_ float64) error {
	return nil
}

func (p *FilmPackage) CalculatePrice(basePrice float64) float64 {
	return basePrice + PriceFilm
}
