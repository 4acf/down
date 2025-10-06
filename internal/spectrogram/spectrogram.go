package spectrogram

import (
	"down/internal/utils"
	"errors"
	"image"
	"math"

	"github.com/r9y9/gossp/stft"
	"github.com/r9y9/gossp/window"
)

type Spectrogram struct {
	stft stft.STFT
}

func NewSpectrogram() Spectrogram {
	frameLength := 2048
	stft := stft.STFT{
		FrameShift: frameLength / 8,
		FrameLen:   frameLength,
		Window:     window.CreateHanning(frameLength),
	}
	return Spectrogram{stft: stft}
}

func (spectrogram *Spectrogram) Image(samples []float64, bounds image.Rectangle) ([][]float64, error) {

	complexValues := spectrogram.stft.STFT(samples)

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

	for i := range xRes {
		for j := range yRes {
			frame := utils.ScaleInt(i, xRes, frames)
			frequencyBin := utils.ScaleInt(j, yRes, frequencyBins)
			abs := math.Abs(real(complexValues[frame][frequencyBin]))
			db := 20 * math.Log10(abs+1e-10)
			db = utils.Clamp(db, minDb, 0)
			result[i][j] = (db + (minDb * -1)) / (minDb * -1)
		}
	}

	return result, nil
}
