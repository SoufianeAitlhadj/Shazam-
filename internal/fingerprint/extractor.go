package fingerprint

import (
	"math"

	"gonum.org/v1/gonum/dsp/fourier"
)

type Fingerprint struct {
	Hash       int64
	TimeOffset int
}

func Extract(samples []float64, sampleRate int) []Fingerprint {
	windowSize := 2048
	hopSize := 1024

	if len(samples) < windowSize {
		return nil
	}

	fft := fourier.NewFFT(windowSize)

	var fingerprints []Fingerprint

	for i := 0; i+windowSize < len(samples); i += hopSize {

		window := samples[i : i+windowSize]
		coeffs := fft.Coefficients(nil, window)

		maxMag := 0.0
		maxIndex := 0

		for k := 0; k < len(coeffs)/2; k++ {
			mag := magnitude(coeffs[k])
			if mag > maxMag {
				maxMag = mag
				maxIndex = k
			}
		}

		hash := int64(maxIndex)

		fingerprints = append(fingerprints, Fingerprint{
			Hash:       hash,
			TimeOffset: i,
		})
	}

	return fingerprints
}

func magnitude(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}
