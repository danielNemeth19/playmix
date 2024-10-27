package mocks

import (
	"fmt"
	"os"
)

const (
	BoxHeaderSize = 8
	timeScale     = 1000
)

// FtypBox
// []uint8 len: 8, cap: 8, [0,0,0,32,102,116,121,112]
// size := buf[0:4] -> 32
// boxtype := buf[4:8] -> "ftyp"
// 105,115,111,109,0,0,2,0,105,115,111,109,105,115,111,50,97,118,99,49,109,112,52,49 -> remaining 24 bytes from 32
// MajorBrand: "isom"
// Minorversion: 512
// CompatibleBrands: ["isom","iso2","avc1","mp41"]

// MoovBox
// []uint8 len: 8, cap: 8, [0,8,255,225,109,111,111,118]
// size := buf[0:4] -> 589793
// boxtype := buf[4:8] -> "moov"
// mvhd starts from 40 [0,0,0,108,109,118,104,100]
// - size 108
// - name "mvhd"
// trak starts from 148 []
// - size 259249
// - name "trak"
// trak starts from 259,397 [0,5,10,90,116,114,97,107]
// - size 330330
// - name "trak"
// udta start from 589727 [0,0,0,98,117,100,116,97]

// Free?
// []uint8 len: 8, cap: 8, [0,0,0,8,102,114,101,101]
// size := buf[0:4] -> 8
// boxtype := buf[4:8] -> "free"

// mdat
// []uint8 len: 8, cap: 8, [13,185,115,113,109,100,97,116]
// size := buf[0:4] -> 230257521
// boxtype := buf[4:8] -> mdat

func addStringAsByte(b []uint8, s string) []uint8 {
	for _, c := range s {
		b = append(b, uint8(c))
	}
	return b
}

func addIntAs4Bytes(b []uint8, number int) []uint8 {
	bytes := []uint8{
		uint8(number >> 24),
		uint8(number >> 16),
		uint8(number >> 8),
		uint8(number),
	}
    for _, byte := range bytes {
        b = append(b, byte)
    }
	return b
}

func addIntAs2Bytes(b []uint8, number int) []uint8 {
	bytes := []uint8{
		uint8(number >> 8),
		uint8(number),
	}
    for _, byte := range bytes {
        b = append(b, byte)
    }
	return b
}

func addPadding(b []uint8, n int) []uint8 {
    for i := 0; i < n; {
        b = append(b, uint8(0))
        i++
    }
    return b
}

// header is 8 bytes, where [0:4] encodes the size of the box
// and [4:8] encodes the name of the box
func createHeader(size int, name string) []uint8 {
    var header []uint8
    header = addIntAs4Bytes(header, size)
    header = addStringAsByte(header, name)
	return header
}

func createFtypBox() []uint8 {
	ftypBox := createHeader(16, "ftyp")
	majorBrand := "isom"
    minorVersion := 512

    ftypBox = addStringAsByte(ftypBox, majorBrand)
    ftypBox = addIntAs4Bytes(ftypBox, minorVersion)
	return ftypBox
}

func createMoovBox(seconds int) []uint8 {
	var data []uint8
	moovHeader := createHeader(48, "moov")
	mhvdHeader := createHeader(40, "mvhd")
	data = append(data, moovHeader...)
	data = append(data, mhvdHeader...)
 
    duration := seconds * timeScale
    rate, volume := 65536, 256

	data = addPadding(data, 12)
    data = addIntAs4Bytes(data, timeScale)
    data = addIntAs4Bytes(data, duration)
    data = addIntAs4Bytes(data, rate)
    data = addIntAs2Bytes(data, volume)
    data = addPadding(data, 6)
	return data
}

func CreateData(seconds int) []uint8 {
	var data []uint8
	ftypHeader := createFtypBox()
	data = append(data, ftypHeader...)
	moovBox := createMoovBox(seconds)
	data = append(data, moovBox...)
	return data
}

func WriteFile(seconds int, path string) []uint8 {
	f := CreateData(seconds)
	fmt.Printf("data is: %v\n", f)
	fmt.Printf("Len of f: %d\n", len(f))
	err := os.WriteFile(path, f, 0644)
	if err != nil {
		fmt.Println("error")
	}
	return f
}
