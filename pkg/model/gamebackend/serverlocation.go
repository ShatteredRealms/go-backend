package gamebackend

type ServerLocation string

var (
	ServerLocations = map[string]struct{}{}

	USCentral1ServerLocation = addServerLocation("us-central-1")
)

func addServerLocation(location string) ServerLocation {
	ServerLocations[location] = struct{}{}
	return ServerLocation(location)
}
