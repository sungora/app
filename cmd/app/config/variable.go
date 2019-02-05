package config

type configMain struct {
	SampleConfig *sampleConfig
}

type sampleConfig struct {
	Name     string
	IsAccess bool
	Balance  float64
}

var ServiceAccess *sampleConfig
