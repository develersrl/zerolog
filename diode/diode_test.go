package diode_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/develersrl/zerolog"
	"github.com/develersrl/zerolog/diode"
	"github.com/develersrl/zerolog/internal/cbor"
)

func TestNewWriter(t *testing.T) {
	buf := bytes.Buffer{}
	w := diode.NewWriter(&buf, 1000, 10*time.Millisecond, func(missed int) {
		fmt.Printf("Dropped %d messages\n", missed)
	})
	log := zerolog.New(w)
	log.Print("test")

	w.Close()
	want := "{\"level\":\"debug\",\"message\":\"test\"}\n"
	got := cbor.DecodeIfBinaryToString(buf.Bytes())
	if got != want {
		t.Errorf("Diode New Writer Test failed. got:%s, want:%s!", got, want)
	}
}

func Benchmark(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stderr)
	w := diode.NewWriter(ioutil.Discard, 100000, 10*time.Millisecond, nil)
	log := zerolog.New(w)
	defer w.Close()

	b.SetParallelism(1000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("test")
		}
	})

}
