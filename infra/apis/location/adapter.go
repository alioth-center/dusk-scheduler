package location

import "context"

type Address struct {
	Region    string
	City      string
	Longitude float64
	Latitude  float64
}

type Location interface {
	DetectIP(ctx context.Context, ip string) (result *Address, err error)
}
