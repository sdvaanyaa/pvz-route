package packaging

const (
	MaxWeightBag = 10.0
	MaxWeightBox = 30.0

	PriceBag  = 5.0
	PriceBox  = 20.0
	PriceFilm = 1.0
)

type CompositePackage struct {
	components []PackageStrategy
}

func NewCompositePackage(parts ...PackageStrategy) *CompositePackage {
	return &CompositePackage{components: parts}
}

func (c *CompositePackage) ValidateWeight(weight float64) error {
	for _, part := range c.components {
		if err := part.ValidateWeight(weight); err != nil {
			return err
		}
	}
	return nil
}

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
