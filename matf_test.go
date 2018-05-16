package matf

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

var (
	noMatf     = []byte{0x6e, 0x6f, 0x4d, 0x61, 0x74, 0x66}
	matfHeader = []byte{0x41, 0x4d, 0x4c, 0x54, 0x42, 0x41, 0x35, 0x20, 0x30, 0x2e, 0x4d, 0x20, 0x54, 0x41, 0x66, 0x2d,
		0x6c, 0x69, 0x2c, 0x65, 0x77, 0x20, 0x69, 0x72, 0x74, 0x74, 0x6e, 0x65, 0x62, 0x20, 0x20, 0x79,
		0x63, 0x4f, 0x61, 0x74, 0x65, 0x76, 0x34, 0x20, 0x32, 0x2e, 0x32, 0x2e, 0x20, 0x2c, 0x30, 0x32,
		0x38, 0x31, 0x30, 0x2d, 0x2d, 0x34, 0x31, 0x30, 0x31, 0x20, 0x3a, 0x32, 0x34, 0x35, 0x32, 0x3a,
		0x20, 0x36, 0x54, 0x55, 0x20, 0x43, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x01, 0x00, 0x4d, 0x49}
	matfHeaderMI = []byte{0x41, 0x4d, 0x4c, 0x54, 0x42, 0x41, 0x35, 0x20, 0x30, 0x2e, 0x4d, 0x20, 0x54, 0x41, 0x66, 0x2d,
		0x6c, 0x69, 0x2c, 0x65, 0x77, 0x20, 0x69, 0x72, 0x74, 0x74, 0x6e, 0x65, 0x62, 0x20, 0x20, 0x79,
		0x63, 0x4f, 0x61, 0x74, 0x65, 0x76, 0x34, 0x20, 0x32, 0x2e, 0x32, 0x2e, 0x20, 0x2c, 0x30, 0x32,
		0x38, 0x31, 0x30, 0x2d, 0x2d, 0x34, 0x31, 0x30, 0x31, 0x20, 0x3a, 0x32, 0x34, 0x35, 0x32, 0x3a,
		0x20, 0x36, 0x54, 0x55, 0x20, 0x43, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x01, 0x00, 0x49, 0x4d}
	oneXoneMatrix = []byte{0x00, 0x0e, 0x00, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00,
		0x00, 0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x08, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x63, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x01, 0x00, 0x02, 0x00, 0x00}
)

func TestCheckIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int
		out  int
	}{
		{name: "1", in: 1, out: 8},
		{name: "7", in: 7, out: 8},
		{name: "8", in: 8, out: 8},
		{name: "9", in: 9, out: 16},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ret := checkIndex(tc.in)
			if ret != tc.out {
				t.Fatalf("Expected: %v \t Got: %v", tc.out, ret)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	t.Parallel()

	tdir, ferr := ioutil.TempDir("", "TestOpen")
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer os.RemoveAll(tdir)

	notValid, ferr := ioutil.TempFile(tdir, "noMatf.mat")
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer os.Remove(notValid.Name())

	ferr = ioutil.WriteFile(notValid.Name(), noMatf, 0644)
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer notValid.Close()

	headerOnly, ferr := ioutil.TempFile(tdir, "headerOnly.mat")
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer os.Remove(headerOnly.Name())

	ferr = ioutil.WriteFile(headerOnly.Name(), matfHeader, 0644)
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer headerOnly.Close()

	headerOnlyMI, ferr := ioutil.TempFile(tdir, "headerOnlyMI.mat")
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer os.Remove(headerOnlyMI.Name())

	ferr = ioutil.WriteFile(headerOnlyMI.Name(), matfHeaderMI, 0644)
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer headerOnlyMI.Close()

	tests := []struct {
		name string
		in   string
		err  string
		mat  *Matf
	}{
		{name: "Empty Input", in: "", err: "no such file or directory"},
		{name: "No Matf", in: notValid.Name(), err: "Could not read enough bytes"},
		{name: "Header Only", in: headerOnly.Name()},
		{name: "Header Only MI", in: headerOnlyMI.Name()},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Open(tc.in)
			if err != nil {
				if matched, _ := regexp.MatchString(tc.err, err.Error()); matched == false {
					t.Fatalf("Error matching regex: %v \t Got: %v", tc.err, err)
				}
			} else if len(tc.err) != 0 {
				t.Fatalf("Expected error, got none")
			}
		})
	}
}

func TestDecompressData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  []byte
		output []byte
		err    string
	}{
		{name: "DeadCell", input: []byte{0x78, 0x9c, 0xba, 0xb7, 0xf6, 0x9c, 0x20, 0x20, 0x00, 0x00, 0xff, 0xff, 0x07, 0x30, 0x02, 0x6b}, output: []byte{0xDE, 0xAD, 0xCE, 0x11}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := decompressData(tc.input)
			if err != nil {
				if matched, _ := regexp.MatchString(tc.err, err.Error()); matched == false {
					t.Fatalf("Error matching regex: %v \t Got: %v", tc.err, err)
				}
				return
			} else if len(tc.err) != 0 {
				t.Fatalf("Expected error, got none")
			}
			fmt.Printf("Expected: %#v\tGot: %#v\n", tc.output, output)
		})
	}
}
