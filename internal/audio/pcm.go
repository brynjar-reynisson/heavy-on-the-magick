package audio

import "encoding/binary"

// ToStereo16 converts mono PCM samples (in [-1, 1]) into interleaved
// 16-bit signed little-endian stereo bytes (left/right duplicated) — the
// format ebiten's audio player (and most live audio backends) expect.
// This is the bridge between our square-wave synthesis (synth.go, which
// only knows about the ZX Spectrum's single mono beeper channel) and any
// stereo playback API.
func ToStereo16(samples []float32) []byte {
	out := make([]byte, 0, len(samples)*4) // 2 bytes/channel * 2 channels
	for _, s := range samples {
		if s > 1 {
			s = 1
		}
		if s < -1 {
			s = -1
		}
		v := int16(s * 32767)
		var buf [2]byte
		binary.LittleEndian.PutUint16(buf[:], uint16(v))
		out = append(out, buf[0], buf[1], buf[0], buf[1]) // L, R
	}
	return out
}
