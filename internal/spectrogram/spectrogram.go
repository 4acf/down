package spectrogram

import (
	"down/config"
	"down/internal/utils"
	"errors"
	"image"
	"math"

	"github.com/r9y9/gossp/stft"
	"github.com/r9y9/gossp/window"
)

type Spectrogram struct {
	config *config.Config
	stft   stft.STFT
}

func NewSpectrogram(config *config.Config) Spectrogram {
	frameLength := 2048
	stft := stft.STFT{
		FrameShift: frameLength / 8,
		FrameLen:   frameLength,
		Window:     window.CreateHanning(frameLength),
	}
	return Spectrogram{
		config: config,
		stft:   stft,
	}
}

func (spectrogram *Spectrogram) Image(samples []float64, bounds image.Rectangle) ([][]float64, error) {

	spectrogramProgressBar := utils.NewProgressBar(
		"Running short time fourier transform...",
		"Short time fourier transform complete.",
		1,
	)

	if spectrogram.config.ProgressEnabled() {
		spectrogramProgressBar.UpdateConsole(0)
	}

	complexValues := spectrogram.stft.STFT(samples)

	if spectrogram.config.ProgressEnabled() {
		spectrogramProgressBar.UpdateConsole(1)
	}

	frames := len(complexValues)
	if frames == 0 {
		return nil, errors.New("spectrogram has insufficient frame information")
	}

	frequencyBins := len(complexValues[0])
	if frequencyBins == 0 {
		return nil, errors.New("spectrogram has insufficient frequency information")
	}
	frequencyBins /= 2

	xRes := bounds.Dx()
	yRes := bounds.Dy()

	minDb := float64(-120)

	result := make([][]float64, xRes)
	for i := range result {
		result[i] = make([]float64, yRes)
	}

	normalizationProgressBar := utils.NewProgressBar(
		"Generating spectrogram...",
		"Spectrogram generation complete.",
		xRes,
	)

	for i := range xRes {
		for j := range yRes {
			frame := utils.ScaleInt(i, xRes, frames)
			frequencyBin := utils.ScaleInt(j, yRes, frequencyBins)
			abs := math.Abs(real(complexValues[frame][frequencyBin]))
			db := 20 * math.Log10(abs+1e-10)
			db = utils.ClampFloat64(db, minDb, 0)
			result[i][j] = (db + (minDb * -1)) / (minDb * -1)
		}
		if spectrogram.config.ProgressEnabled() {
			normalizationProgressBar.UpdateConsole(i - bounds.Min.X + 1)
		}
	}

	return result, nil
}
