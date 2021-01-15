package telematics

type Segment struct {
	StartDate           Timestamp        `bson:"start"`
	EndDate             Timestamp        `bson:"end"`
	RelatedIDs          []string         `bson:"rel"`
	Locations           Locations        `bson:"locations"`
	SamplingInterval    Duration         `bson:"interval"`
	AccelerationSamples AccelerationData `bson:"samples"`
	AttitudeSamples     AttitudeData     `bson:"att"`
}

type File struct {
	Version  uint      `bson:"ver"`
	Segments []Segment `bson:"series"`
}
