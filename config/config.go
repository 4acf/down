package config

type Config struct {
	audioOutputDirectory       string
	spectrogramOutputDirectory string
	progressEnabled            bool
}

func NewConfig(audioOutputDirectory string, spectrogramOutputDirectory string, progressEnabled bool) Config {
	return Config{
		audioOutputDirectory:       audioOutputDirectory,
		spectrogramOutputDirectory: spectrogramOutputDirectory,
		progressEnabled:            progressEnabled,
	}
}

func (config *Config) AudioOutputDirectory() string {
	return config.audioOutputDirectory
}

func (config *Config) SpectrogramOutputDirectory() string {
	return config.spectrogramOutputDirectory
}

func (config *Config) ProgressEnabled() bool {
	return config.progressEnabled
}
