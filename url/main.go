package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

// EncodeDecoder is struct for encoding
type EncodeDecoder struct{}

func (ed *EncodeDecoder) encode(v string) (string, error) {
	return url.QueryEscape(v), nil
}

func (ed *EncodeDecoder) decode(v string) (string, error) {
	decoded, err := url.QueryUnescape(v)
	if err != nil {
		err = errors.Wrapf(err, "failed to decode %s", v)
	}
	return decoded, err
}

func main() {
	var (
		decode = flag.Bool("d", false, "decodes input")
		input  = flag.String("i", "-", "input file (default: \"-\" for stdin)")
		output = flag.String("o", "-", "output file (default: \"-\" for stdout)")
	)
	flag.Parse()

	var r io.Reader
	var w io.Writer

	if *input == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(*input)
		r = f
		defer f.Close()
		if err != nil {
			log.Fatal(errors.Wrapf(err, "%s could not opened", *input))
		}
	}

	if *output == "-" {
		w = os.Stdout
	} else {
		f, err := os.Open(*output)
		w = f
		defer f.Close()
		if err != nil {
			log.Fatal(errors.Wrapf(err, "%s could not opened", *output))
		}
	}

	scanner := bufio.NewScanner(r)
	ec := &EncodeDecoder{}

	for scanner.Scan() {
		var t string
		var err error
		if *decode {
			t, err = ec.decode(scanner.Text())
		} else {
			t, err = ec.encode(scanner.Text())
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(w, t)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
