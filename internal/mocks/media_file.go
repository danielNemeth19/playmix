package mocks

import (
	"fmt"
	"os"
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

func CreateData() []uint8 {
	// return []uint8{
	// 0,0,0,32,102,116,121,112,0,8,255,225,109,111,111,118,0,0,0,8,102,114,101,101,13,185,115,113,109,100,97,116,
	// }
	var data []uint8
	ftypHeader := []uint8{0, 0, 0, 32, 102, 116, 121, 112}
	data = append(data, ftypHeader...)
    ftypMeta := []uint8{105,115,111,109,0,0,2,0,105,115,111,109,105,115,111,50,97,118,99,49,109,112,52,49}
	data = append(data, ftypMeta...)
	fmt.Printf("data is: %v\n", data)
	return data
}

func WriteFile() []uint8 {
	f := CreateData()
	fmt.Printf("Len of f: %d\n", len(f))
	err := os.WriteFile("/home/daniel/Videos/test.mp4", f, 0644)
	if err != nil {
		fmt.Println("error")
	}
	return f
}
