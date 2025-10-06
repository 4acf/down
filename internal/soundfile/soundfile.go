package soundfile

import (
	"down/internal/utils"
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/r9y9/gossp/stft"
	"github.com/r9y9/gossp/window"
)

const (
	WAV_EXTENSION = ".wav"
)

type Soundfile struct {
	img         *image.Image
	name        string
	format      *audio.Format
	sampleRate  int
	numChannels int
	bitDepth    int
	data        []float64
}

func NewSoundfile(img *image.Image, name string) Soundfile {
	sampleRate := 44100
	numChannels := 1
	bitDepth := 16

	format := &audio.Format{
		NumChannels: numChannels,
		SampleRate:  sampleRate,
	}

	name = appendWav(name)

	return Soundfile{
		img:         img,
		name:        name,
		format:      format,
		sampleRate:  sampleRate,
		numChannels: numChannels,
		bitDepth:    bitDepth,
		data:        make([]float64, 0),
	}
}

func (soundfile *Soundfile) Name() string {
	return soundfile.name
}

func (soundfile *Soundfile) Data() []float64 {
	return soundfile.data
}

func appendWav(name string) string {
	name = strings.TrimSpace(name)

	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "_")
	}

	if strings.HasSuffix(name, ".") {
		return name + WAV_EXTENSION
	} else if !strings.HasSuffix(strings.ToLower(name), WAV_EXTENSION) {
		return name + WAV_EXTENSION
	}
	return name
}

func (soundfile *Soundfile) Wav(outputDirectory string, progressEnabled bool) error {
	outputFilepath := filepath.Join(outputDirectory, soundfile.name)

	out, err := os.Create(outputFilepath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", soundfile.name, err)
	}
	defer out.Close()

	encoder := wav.NewEncoder(out, soundfile.sampleRate, soundfile.bitDepth, soundfile.numChannels, 1)
	defer encoder.Close()

	bounds := (*soundfile.img).Bounds()

	progressBar := utils.NewProgressBar(
		fmt.Sprintf("Writing to %s...", soundfile.name),
		fmt.Sprintf("Writeout to %s complete.", soundfile.name),
		bounds.Dx(),
	)

	for x := bounds.Min.X; x <= bounds.Max.X; x++ {
		freqs := soundfile.getColumnFrequencies(x)
		soundfile.addSine(encoder, freqs)
		if progressEnabled {
			progressBar.UpdateConsole(x - bounds.Min.X)
		}
	}

	return nil
}

func (soundfile *Soundfile) getColumnFrequencies(x int) []float64 {
	bounds := (*soundfile.img).Bounds()
	height := bounds.Dy()
	freqs := make([]float64, 0, height*2)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		r, g, b, _ := (*soundfile.img).At(x, y).RGBA()
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

func (soundfile *Soundfile) addSine(encoder *wav.Encoder, freqs []float64) {
	if len(freqs) == 0 {
		return
	}

	const (
		sampleRate = 44100
		duration   = 0.2
	)

	numSamples := int(duration * float64(soundfile.sampleRate))
	maxAmplitude := math.Pow(2, float64(soundfile.bitDepth)-1)

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
		Format:         soundfile.format,
		SourceBitDepth: soundfile.bitDepth,
		Data:           make([]int, numSamples),
	}

	for i, sample := range buffer {
		sample /= float64(len(freqs) / 2)
		sample = math.Max(-1, math.Min(1, sample))
		intBuf.Data[i] = int(sample * maxAmplitude)
	}

	// todo: only run one of these blocks depending on if config is to keep audio files or not
	// and determine if its really necessary to go from []float64 to []int and then back to []float64
	_ = encoder.Write(intBuf)

	normalizeValue := math.Pow(2, float64(soundfile.bitDepth)-1)
	float64Array := make([]float64, len(intBuf.Data))
	for i, value := range intBuf.Data {
		float64Array[i] = float64(value) / normalizeValue
	}
	soundfile.data = append(soundfile.data, float64Array...)
}

func (soundfile *Soundfile) Spectrogram() [][]complex128 {
	frameLength := 4096
	stft := stft.STFT{
		FrameShift: frameLength / 2,
		FrameLen:   frameLength,
		Window:     window.CreateHanning(frameLength),
	}

	return stft.STFT(soundfile.data)
}
