//go:build !server && porttest

package ail

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"hash"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/opennox/libs/datapath"
)

// Freeze only after repeated captures from the unchanged C decoder.
const mp3AssetBaselineSHA256 = "7dc3362576eb00660fd68a836d17af6012c4d1c0595d137958109d62228add0b"

const mp3BaselineMaxPCM = 1152 * 2 // MINIMP3_MAX_SAMPLES_PER_FRAME: maximum stereo frame.

type mp3BaselineRecord struct {
	Name              string `json:"name"`
	InputSHA256       string `json:"input_sha256"`
	Format            string `json:"format"`
	Channels          int    `json:"channels"`
	SampleRate        int    `json:"sample_rate"`
	PCMSHA256         string `json:"pcm_sha256"`
	ReportedSamples   int64  `json:"reported_samples"`
	DecodeSequenceSHA string `json:"decode_sequence_sha256"`
	DecodeCalls       int    `json:"decode_calls"`
	OutputWriteSHA256 string `json:"full_capacity_output_write_sha256"`
}

func TestMP3AssetBaseline(t *testing.T) {
	dataRoot := datapath.Data()
	if dataRoot == "" || !datapath.Found() {
		t.Fatal("Nox data directory not found; set NOX_DATA to the extracted game data")
	}
	dir := filepath.Join(dataRoot, "Dialog")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read Dialog asset directory %q: %v", dir, err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })

	var records []mp3BaselineRecord
	seenMP3 := 0
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		input, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read input %q: %v", path, err)
		}
		if len(input) < 12 || string(input[:4]) != "RIFF" || string(input[8:12]) != "WAVE" {
			continue
		}
		r, err := openWav(path)
		if err != nil {
			t.Fatalf("open WAV %q: %v", path, err)
		}
		if r.Format() != "WAV+MP3" {
			r.Close()
			continue // Select from the parsed RIFF format tag, never the extension.
		}
		seenMP3++
		record := decodeMP3Baseline(t, filepath.ToSlash(filepath.Join("Dialog", entry.Name())), input, r)
		records = append(records, record)
	}
	if seenMP3 == 0 {
		t.Fatalf("no WAV assets with RIFF MP3 format tag found under %q", dir)
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	digest := sha256.Sum256(data)
	got := hex.EncodeToString(digest[:])
	out := os.Getenv("OPENNOX_MP3_BASELINE_CAPTURE")
	if mp3AssetBaselineSHA256 == "" && out == "" {
		t.Fatal("MP3 baseline is not frozen; an explicit capture path is required for baseline collection")
	}
	if out != "" {
		if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(out, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("captured %d MP3-in-WAVE assets; sha256=%s", seenMP3, got)
	if mp3AssetBaselineSHA256 != "" && got != mp3AssetBaselineSHA256 {
		t.Fatalf("MP3 baseline mismatch: got %s, want %s", got, mp3AssetBaselineSHA256)
	}
}

func decodeMP3Baseline(t *testing.T, name string, input []byte, r *wavReader) mp3BaselineRecord {
	t.Helper()
	defer r.Close()
	const guardWords = 8
	const guardLeft int16 = 0x6a5a
	const guardRight int16 = 0x3bc7
	backing := make([]int16, guardWords+mp3BaselineMaxPCM+guardWords)
	for i := 0; i < guardWords; i++ {
		backing[i] = guardLeft
		backing[guardWords+mp3BaselineMaxPCM+i] = guardRight
	}
	out := backing[guardWords : guardWords+mp3BaselineMaxPCM : guardWords+mp3BaselineMaxPCM]

	pcmHash := sha256.New()
	sequenceHash := sha256.New()
	writeHash := sha256.New()
	channels, rate := r.Channels(), r.SampleRate()
	if channels < 1 || channels > 2 || rate <= 0 {
		t.Fatalf("%s: unexpected decoder metadata channels=%d rate=%d", name, channels, rate)
	}
	var total int64
	noProgress := 0
	const noProgressLimit = 64
	const maxSamplesPerChannel = 5 * 60
	const maxCalls = maxSamplesPerChannel*192000*2/576 + 10000
	calls := 0
	for {
		if total > int64(maxSamplesPerChannel*rate*channels) {
			t.Fatalf("%s: decoded more than five minutes of reported PCM", name)
		}
		for i := range out {
			out[i] = 0 // Makes unreported/full-capacity writes observable per call.
		}
		for i := 0; i < guardWords; i++ {
			backing[i] = guardLeft
			backing[guardWords+mp3BaselineMaxPCM+i] = guardRight
		}
		n, ok := r.Decode(out)
		calls++
		if n < 0 || n > len(out) {
			t.Fatalf("%s: Decode returned invalid count %d for capacity %d", name, n, len(out))
		}
		for i := 0; i < guardWords; i++ {
			if backing[i] != guardLeft || backing[guardWords+mp3BaselineMaxPCM+i] != guardRight {
				t.Fatalf("%s: Decode wrote outside the %d-word output at call %d", name, len(out), calls)
			}
		}
		appendSequence(sequenceHash, int32(n), ok)
		hashPCMBytes(pcmHash, out[:n])
		hashPCMBytes(writeHash, out)
		total += int64(n)
		if n == 0 && ok {
			noProgress++
			if noProgress > noProgressLimit {
				t.Fatalf("%s: Decode made no reported progress for %d consecutive calls", name, noProgress)
			}
		} else {
			noProgress = 0
		}
		if !ok {
			break
		}
		if calls >= maxCalls {
			t.Fatalf("%s: exceeded bounded Decode call count %d", name, maxCalls)
		}
	}
	inputHash := sha256.Sum256(input)
	return mp3BaselineRecord{
		Name:              name,
		InputSHA256:       hex.EncodeToString(inputHash[:]),
		Format:            r.Format(),
		Channels:          channels,
		SampleRate:        rate,
		PCMSHA256:         hex.EncodeToString(pcmHash.Sum(nil)),
		ReportedSamples:   total,
		DecodeSequenceSHA: hex.EncodeToString(sequenceHash.Sum(nil)),
		DecodeCalls:       calls,
		OutputWriteSHA256: hex.EncodeToString(writeHash.Sum(nil)),
	}
}

func appendSequence(h hash.Hash, n int32, ok bool) {
	var b [5]byte
	binary.LittleEndian.PutUint32(b[:4], uint32(n))
	if ok {
		b[4] = 1
	}
	_, _ = h.Write(b[:])
}

func hashPCMBytes(h hash.Hash, samples []int16) {
	var encoded [2 * mp3BaselineMaxPCM]byte
	for i, sample := range samples {
		binary.LittleEndian.PutUint16(encoded[2*i:], uint16(sample))
	}
	_, _ = h.Write(encoded[:2*len(samples)])
}
