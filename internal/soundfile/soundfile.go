package soundfile

import (
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func Wav(img image.Image, outputDirectory string) error {
	outputFilepath := filepath.Join(outputDirectory, "test.wav")

	out, err := os.Create(outputFilepath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", outputFilepath, err)
	}
	defer out.Close()

	const (
		sampleRate  = 44100
		numChannels = 1
		bitDepth    = 16
	)

	encoder := wav.NewEncoder(out, sampleRate, bitDepth, numChannels, 1)
	defer encoder.Close()

	format := &audio.Format{
		NumChannels: numChannels,
		SampleRate:  sampleRate,
	}

	bounds := img.Bounds()
	height := bounds.Dy()

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		freqs := freqsFromColumn(img, x, bounds, height)
		addSine(encoder, freqs, format)
	}

	return nil
}

func freqsFromColumn(img image.Image, x int, bounds image.Rectangle, height int) []float64 {
	freqs := make([]float64, 0, height*2)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		r, g, b, _ := img.At(x, y).RGBA()
		r >>= 8
		g >>= 8
		b >>= 8

		if r > 10 || g > 10 || b > 10 {
			c := 4.25 - 4.25*float64(r+g+b)/768
			percentage := float64(y+1) / float64(height+1)
			d := 22000 - (percentage * 22000)
			freqs = append(freqs, d, c)
		}
	}
	return freqs
}

func addSine(encoder *wav.Encoder, freqs []float64, format *audio.Format) {
	if len(freqs) == 0 {
		return
	}

	const (
		sampleRate = 44100
		duration   = 0.2
	)

	numSamples := int(duration * sampleRate)
	maxAmplitude := math.Pow(2, 15)

	buffer := make([]float64, numSamples)

	for i := 0; i < len(freqs); i += 2 {
		freq := freqs[i]
		attenuation := math.Pow(10, freqs[i+1])
		scaling := 10 / attenuation
		for pos := range numSamples {
			time := float64(pos) / sampleRate
			buffer[pos] += math.Sin(2*math.Pi*freq*time) * scaling
		}
	}

	intBuf := &audio.IntBuffer{
		Format:         format,
		SourceBitDepth: 16,
		Data:           make([]int, numSamples),
	}

	for i, sample := range buffer {
		sample /= float64(len(freqs) / 2)
		sample = math.Max(-1, math.Min(1, sample))
		intBuf.Data[i] = int(sample * maxAmplitude)
	}

	_ = encoder.Write(intBuf)
}
