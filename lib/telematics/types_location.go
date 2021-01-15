package telematics

import (
	"math"
)

// Distance in meters
type Distance float64

// Speed in meters-per-second
type Speed float64

type Location struct {
	Date      Timestamp `bson:"time" json:"date"`
	Latitude  float64   `bson:"lat" json:"lat"`
	Longitude float64   `bson:"lon" json:"long"`
}

type Locations []*Location

func (a Locations) Len() int           { return len(a) }
func (a Locations) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Locations) Less(i, j int) bool { return a[i].Date.Before(a[j].Date.Time) }

func (a *Location) DistanceFrom(b *Location) Distance {
	return distance(a.Latitude, a.Longitude, b.Latitude, b.Longitude)
}

func (a *Location) SpeedFrom(b *Location) Speed {
	dst := float64(a.DistanceFrom(b))
	dur := a.Date.Sub(b.Date.Time)

	return Speed(math.Abs(dst / dur.Seconds()))
}
