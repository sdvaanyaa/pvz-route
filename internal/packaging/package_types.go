package packaging

import (
	"errors"
	"fmt"
)

var ErrUnknownPackageType = errors.New("unknown package type")

var packageStrategies = map[string]PackageStrategy{
	PackageNone:    &NonePackage{},
	PackageBag:     &BagPackage{},
	PackageBox:     &BoxPackage{},
	PackageFilm:    &FilmPackage{},
	PackageBagFilm: NewCompositePackage(&BagPackage{}, &FilmPackage{}),
	PackageBoxFilm: NewCompositePackage(&BoxPackage{}, &FilmPackage{}),
}

// GetPackageStrategy returns the PackageStrategy corresponding to the given packageType string.
// If the packageType is not recognized, it returns an error wrapped with ErrUnknownPackageType.
func GetPackageStrategy(packageType string) (PackageStrategy, error) {
	strategy, ok := packageStrategies[packageType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownPackageType, packageType)
	}
	return strategy, nil
}
