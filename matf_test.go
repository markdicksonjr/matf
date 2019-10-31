package matf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

var (
	noMatf     = []byte{0x6e, 0x6f, 0x4d, 0x61, 0x74, 0x66}
	matfHeader = []byte{0x41, 0x4d, 0x4c, 0x54, 0x42, 0x41, 0x35, 0x20, 0x30,
		0x2e, 0x4d, 0x20, 0x54, 0x41, 0x66, 0x2d, 0x6c, 0x69, 0x2c, 0x65,
		0x77, 0x20, 0x69, 0x72, 0x74, 0x74, 0x6e, 0x65, 0x62, 0x20, 0x20,
		0x79, 0x63, 0x4f, 0x61, 0x74, 0x65, 0x76, 0x34, 0x20, 0x32, 0x2e,
		0x32, 0x2e, 0x20, 0x2c, 0x30, 0x32, 0x38, 0x31, 0x30, 0x2d, 0x2d,
		0x34, 0x31, 0x30, 0x31, 0x20, 0x3a, 0x32, 0x34, 0x35, 0x32, 0x3a,
		0x20, 0x36, 0x54, 0x55, 0x20, 0x43, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x01, 0x00, 0x4d, 0x49}
	matfHeaderMI = []byte{0x41, 0x4d, 0x4c, 0x54, 0x42, 0x41, 0x35, 0x20,
		0x30, 0x2e, 0x4d, 0x20, 0x54, 0x41, 0x66, 0x2d, 0x6c, 0x69, 0x2c,
		0x65, 0x77, 0x20, 0x69, 0x72, 0x74, 0x74, 0x6e, 0x65, 0x62, 0x20,
		0x20, 0x79, 0x63, 0x4f, 0x61, 0x74, 0x65, 0x76, 0x34, 0x20, 0x32,
		0x2e, 0x32, 0x2e, 0x20, 0x2c, 0x30, 0x32, 0x38, 0x31, 0x30, 0x2d,
		0x2d, 0x34, 0x31, 0x30, 0x31, 0x20, 0x3a, 0x32, 0x34, 0x35, 0x32,
		0x3a, 0x20, 0x36, 0x54, 0x55, 0x20, 0x43, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x01, 0x00, 0x49, 0x4d}
	compressedMatf = []byte{0x4d, 0x41, 0x54, 0x4c, 0x41, 0x42, 0x20, 0x35,
		0x2e, 0x30, 0x20, 0x4d, 0x41, 0x54, 0x2d, 0x66, 0x69, 0x6c, 0x65,
		0x2c, 0x20, 0x77, 0x72, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x20, 0x62,
		0x79, 0x20, 0x4f, 0x63, 0x74, 0x61, 0x76, 0x65, 0x20, 0x34, 0x2e,
		0x32, 0x2e, 0x32, 0x2c, 0x20, 0x32, 0x30, 0x31, 0x38, 0x2d, 0x30,
		0x35, 0x2d, 0x32, 0x35, 0x20, 0x30, 0x39, 0x3a, 0x31, 0x36, 0x3a,
		0x33, 0x38, 0x20, 0x55, 0x54, 0x43, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x00, 0x01, 0x49, 0x4d, 0x0f,
		0x00, 0x00, 0x00, 0x31, 0x00, 0x00, 0x00, 0x78, 0x9c, 0xe3, 0x63,
		0x60, 0x60, 0xf0, 0x00, 0x62, 0x36, 0x20, 0xe6, 0x80, 0xd2, 0x8c,
		0x40, 0xcc, 0x0a, 0xe5, 0x33, 0x22, 0x61, 0x4e, 0x20, 0x0e, 0xc8,
		0xc8, 0xcc, 0x89, 0xf7, 0x4e, 0x2c, 0xa9, 0x62, 0x80, 0x00, 0x4e,
		0xa8, 0x3a, 0x30, 0x10, 0x99, 0xef, 0x00, 0x00, 0x88, 0xd3, 0x05,
		0x0f}
)

func TestAlignIndex(t *testing.T) {

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

	fakeReader := bytes.NewReader([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ret := alignIndex(fakeReader, binary.LittleEndian, tc.in)
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

	testdir, ferr := ioutil.TempDir(tdir, "TestDir")
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer os.RemoveAll(testdir)

	tests := []struct {
		name string
		in   string
		err  string
	}{
		{name: "Empty Input", in: "", err: "no such file or directory"},
		{name: "No Matf", in: notValid.Name(), err: "Could not read enough bytes"},
		{name: "Header Only", in: headerOnly.Name()},
		{name: "Header Only MI", in: headerOnlyMI.Name()},
		{name: "Folder As Input", in: testdir, err: "is not a file"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Open(tc.in)
			if err != nil {
				if matched, _ := regexp.MatchString(tc.err, err.Error()); !matched {
					t.Fatalf("Error matching regex: %v \t Got: %v", tc.err, err)
				} else {
					return
				}
				t.Fatalf("Expected no error, got: %v", err)
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
				if matched, _ := regexp.MatchString(tc.err, err.Error()); !matched {
					t.Fatalf("Error matching regex: %v \t Got: %v", tc.err, err)
				} else {
					return
				}
				t.Fatalf("Expected no error, got: %v", err)
			} else if len(tc.err) != 0 {
				t.Fatalf("Expected error, got none")
			}
			fmt.Printf("Expected: %#v\tGot: %#v\n", tc.output, output)
		})
	}
}

func TestDimensions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mat  MatMatrix
		x    int
		y    int
		z    int
		err  string
	}{
		{name: "1 Dim", mat: MatMatrix{Dim: Dim{X: 2, Y: 0, Z: 0}}, x: 2},
		{name: "2 Dim", mat: MatMatrix{Dim: Dim{X: 3, Y: 5, Z: 0}}, x: 3, y: 5},
		{name: "3 Dim", mat: MatMatrix{Dim: Dim{X: 7, Y: 11, Z: 13}}, x: 7, y: 11, z: 13},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			x, y, z, extra, err := tc.mat.Dimensions()
			if err != nil {
				if matched, _ := regexp.MatchString(tc.err, err.Error()); !matched {
					t.Fatalf("Error matching regex: %v \t Got: %v", tc.err, err)
				} else {
					return
				}
				t.Fatalf("Expected no error, got: %v", err)
			} else if len(tc.err) != 0 {
				t.Fatalf("Expected error, got none")
			}
			if tc.x != x {
				t.Fatalf("Expected x: %d\tgot: %d", tc.x, x)
			} else if tc.y != y {
				t.Fatalf("Expected y: %d\tgot: %d", tc.y, y)
			} else if tc.z != z {
				t.Fatalf("Expected z: %d\tgot: %d", tc.x, z)
			} else if extra != nil {
				t.Fatalf("Expected no extra dimensions, got: %d", len(extra))
			}
		})
	}
}

func TestExtractMatrix(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		err  string
	}{
		{name: "notExpectedArrayFlagSize", data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, err: "will not read 0 bytes"},
		{name: "invalidSmallTag", data: []byte{0x00, 0x01, 0x02, 0x03}, err: "EOF"},
		{name: "tooFewBytes", data: []byte{0x00, 0x01}, err: "EOF"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.data)
			mat, index, err := extractMatrix(r, binary.LittleEndian)
			fmt.Println(err)
			if err != nil {
				if matched, _ := regexp.MatchString(tc.err, err.Error()); !matched {
					t.Fatalf("Error matching regex: %v \t Got: %v", tc.err, err)
				} else {
					return
				}
				t.Fatalf("Expected no error, got: %v", err)
			} else if len(tc.err) != 0 {
				t.Fatalf("Expected error, got none")
			}
			_ = mat
			_ = index
		})
	}
}

func TestReadDimensions(t *testing.T) {
	tests := []struct {
		name string
		data []interface{}
		err  string
	}{
		{name: " 1", data: []interface{}{1, 1}},
		{name: " 2", data: []interface{}{2, 1, 2}},
		{name: " 3", data: []interface{}{3, 1, 2, 3}, err: "More dimensions than exptected"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dim, err := readDimensions(tc.data)
			if err != nil {
				if matched, _ := regexp.MatchString(tc.err, err.Error()); !matched {
					t.Fatalf("Error matching regex: %v \t Got: %v", tc.err, err)
				} else {
					return
				}
				t.Fatalf("Expected no error, got: %v", err)
			} else if len(tc.err) != 0 {
				t.Fatalf("Expected error, got none")
			}
			_ = dim
		})
	}
}

func TestMatf(t *testing.T) {
	tdir, ferr := ioutil.TempDir("", "TestMatf")
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer os.RemoveAll(tdir)

	simple, ferr := ioutil.TempFile(tdir, "simple.mat")
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer os.Remove(simple.Name())

	ferr = ioutil.WriteFile(simple.Name(), compressedMatf, 0644)
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer simple.Close()

	corrupt, ferr := ioutil.TempFile(tdir, "corrupt.mat")
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer os.Remove(corrupt.Name())

	ferr = ioutil.WriteFile(corrupt.Name(), noMatf, 0644)
	if ferr != nil {
		t.Fatal(ferr)
	}
	defer simple.Close()

	tests := []struct {
		name string
		file string
		err  string
	}{
		{name: "simple", file: simple.Name()},
		{name: "corrupt", file: corrupt.Name(), err: "Could not read enough bytes"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			x, err := Open(tc.file)
			if err != nil {
				if matched, _ := regexp.MatchString(tc.err, err.Error()); !matched {
					t.Fatalf("Error matching regex: %v \t Got: %v", tc.err, err)
				} else {
					return
				}
				t.Fatalf("Expected no error, got: %v", err)
			} else if len(tc.err) != 0 {
				t.Fatalf("Expected error, got none")
			}
			defer Close(x)

			for {
				mat, err := ReadDataElement(x)
				if err == io.EOF {
					break
				} else if err != nil {
					t.Fatalf("Could not open test file: %v", err)
				}
				_ = mat
			}
		})
	}
}
