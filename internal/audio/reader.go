package audio

import (
	"os"

	"github.com/go-audio/wav"
)

func ReadWav(path string) ([]float64, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)

	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, 0, err
	}

	samples := make([]float64, len(buf.Data))
	for i, v := range buf.Data {
		samples[i] = float64(v)
	}

	return samples, int(decoder.SampleRate), nil
}
