package discovery

type Device interface {
}

type DiscoveryService interface {
	Discover() ([]Device, error)
}
