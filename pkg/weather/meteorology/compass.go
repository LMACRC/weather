package meteorology

type Direction string

var directions = []Direction{
	"N", "NNE", "NE", "ENE",
	"E", "ESE", "SE", "SSE",
	"S", "SSW", "SW", "WSW",
	"W", "WNW", "NW", "NNW"}

func CardinalDirection(degrees float64) Direction {
	i := int((degrees + 11.25) / 22.5)
	return directions[i%16]
}
