package soundfile

import (
	"down/config"
	"down/internal/filesystem"
	"down/internal/imagefile"
	"down/internal/utils"
	"fmt"
	"image"
	"math"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type Soundfile struct {
	config        *config.Config
	imageContents *image.Image
	imagefile     *imagefile.Imagefile
	format        *audio.Format
	sampleRate    int
	numChannels   int
	bitDepth      int
	data          []float64
}

func NewSoundfile(config *config.Config, imageContents *image.Image, imagefile *imagefile.Imagefile) Soundfile {
	sampleRate := 44100
	numChannels := 1
	bitDepth := 16

	format := &audio.Format{
		NumChannels: numChannels,
		SampleRate:  sampleRate,
	}

	return Soundfile{
		config:        config,
		imageContents: imageContents,
		imagefile:     imagefile,
		format:        format,
		sampleRate:    sampleRate,
		numChannels:   numChannels,
		bitDepth:      bitDepth,
		data:          make([]float64, 0),
	}
}

func (soundfile *Soundfile) Data() []float64 {
	return soundfile.data
}

func (soundfile *Soundfile) Wav() error {

	filename := utils.AppendWav(soundfile.imagefile.Name())
	outputFilepath, err := filesystem.CreateFinalPath(soundfile.config.AudioOutputDirectory(), soundfile.imagefile.Path(), filename)
	if err != nil {
		return err
	}

	out, err := os.Create(outputFilepath)
	if err != nil {
		return err
	}
	defer out.Close()

	encoder := wav.NewEncoder(out, soundfile.sampleRate, soundfile.bitDepth, soundfile.numChannels, 1)
	defer encoder.Close()

	bounds := (*soundfile.imageContents).Bounds()

	progressBar := utils.NewProgressBar(
		fmt.Sprintf("Writing to %s...", filename),
		fmt.Sprintf("Writeout to %s complete.", filename),
		bounds.Dx(),
	)

	for x := bounds.Min.X; x <= bounds.Max.X; x++ {
		freqs := soundfile.getColumnFrequencies(x)
		soundfile.addSine(encoder, freqs)
		if soundfile.config.ProgressEnabled() {
			progressBar.UpdateConsole(x - bounds.Min.X)
		}
	}

	return nil
}

func (soundfile *Soundfile) getColumnFrequencies(x int) []float64 {
	bounds := (*soundfile.imageContents).Bounds()
	height := bounds.Dy()
	freqs := make([]float64, 0, height*2)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		r, g, b, _ := (*soundfile.imageContents).At(x, y).RGBA()
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

	_ = encoder.Write(intBuf)

	normalizeValue := math.Pow(2, float64(soundfile.bitDepth)-1)
	float64Array := make([]float64, len(intBuf.Data))
	for i, value := range intBuf.Data {
		float64Array[i] = float64(value) / normalizeValue
	}
	soundfile.data = append(soundfile.data, float64Array...)
}
