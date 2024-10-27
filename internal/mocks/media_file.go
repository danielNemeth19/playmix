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

func convertStringByte(s string) []uint8 {
	var data []uint8
	for _, c := range s {
		data = append(data, uint8(c))
	}
	return data
}

func convertIntByte(number int) []uint8 {
	bytes := []uint8{
		uint8(number >> 24),
		uint8(number >> 16),
		uint8(number >> 8),
		uint8(number),
	}
	return bytes
}

func convertInt2Bytes(number int) []uint8 {
	bytes := []uint8{
		uint8(number >> 8),
		uint8(number),
	}
	return bytes
}

func createEmpySequence(n int) []uint8 {
	var data []uint8
	for i := 0; i < n; {
		data = append(data, uint8(0))
		i++
	}
	fmt.Println(data)
	return data
}

// header is 8 bytes, where [0:4] encodes the size of the box
// and [4:8] encodes the name of the box
func createHeader(size int, name string) []uint8 {
	header := convertIntByte(size)
	nameInByte := convertStringByte(name)
	header = append(header, nameInByte...)
	return header
}

func createFtypBox() []uint8 {
	ftypBox := createHeader(16, "ftyp")
	majorBrand := convertStringByte("isom")
	minorVersion := convertIntByte(512)

	ftypBox = append(ftypBox, majorBrand...)
	ftypBox = append(ftypBox, minorVersion...)
	return ftypBox
}

func createMoovBox(seconds int) []uint8 {
	var data []uint8
	moovHeader := createHeader(48, "moov")
	mhvdHeader := createHeader(40, "mvhd")
	data = append(data, moovHeader...)
	data = append(data, mhvdHeader...)

	mvhdData := createEmpySequence(12)
	fmt.Printf("Len of start sequence: %d\n", len(mvhdData))

	tS := convertIntByte(timeScale)
	duration := convertIntByte(seconds * timeScale)
	rate := convertIntByte(65536)
	volume := convertInt2Bytes(256)
	filler := createEmpySequence(6)

	mvhdData = append(mvhdData, tS...)
	mvhdData = append(mvhdData, duration...)
	mvhdData = append(mvhdData, rate...)
	mvhdData = append(mvhdData, volume...)
	mvhdData = append(mvhdData, filler...)

	data = append(data, mvhdData...)
	return data
}

func CreateData(seconds int) []uint8 {
	var data []uint8
	ftypHeader := createFtypBox()
	data = append(data, ftypHeader...)
	fmt.Printf("Len of ftyp: %d\n", len(ftypHeader))

	moovBox := createMoovBox(seconds)
	data = append(data, moovBox...)
	fmt.Printf("Len of moovBox: %d\n", len(moovBox))
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
