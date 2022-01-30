package servicecontext

// Info type holds useful info about the currently-running service
type Info struct {
	Name        string
	Environment string
}

var service *Info

// Set gives the singleton pointer a value
func Set(name, env string) {
	service = &Info{
		Name:        name,
		Environment: env,
	}
}

// Get returns the value of the singleton instance pointer
func Get() Info {
	if service == nil {
		panic("cannot get service context before it is set")
	}

	return *service
}

// IsSet returns true if the singleton has been initialised
func IsSet() bool {
	return service != nil
}
