package config

type Config struct {
	outputDirectory string
	progressEnabled bool
	keepAudio       bool
	keepSpectrogram bool
}

func NewConfig(outputDirectory string, progressEnabled bool, keepAudio bool, keepSpectrogram bool) Config {
	return Config{
		outputDirectory: outputDirectory,
		progressEnabled: progressEnabled,
		keepAudio:       keepAudio,
		keepSpectrogram: keepSpectrogram,
	}
}

func (config *Config) OutputDirectory() string {
	return config.outputDirectory
}

func (config *Config) ProgressEnabled() bool {
	return config.progressEnabled
}

func (config *Config) KeepAudio() bool {
	return config.keepAudio
}

func (config *Config) KeepSpectrogram() bool {
	return config.keepSpectrogram
}
