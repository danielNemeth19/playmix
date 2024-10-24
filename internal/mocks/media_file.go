package mocks


// FtypBox
// []uint8 len: 8, cap: 8, [0,0,0,32,102,116,121,112]
// size := buf[0:4] -> 32
// boxtype := buf[4:8] -> "ftyp"

// MoovBox
// []uint8 len: 8, cap: 8, [0,8,255,225,109,111,111,118]
// size := buf[0:4] -> 589793
// boxtype := buf[4:8] -> "moov"

// Free?
// []uint8 len: 8, cap: 8, [0,0,0,8,102,114,101,101]
// size := buf[0:4] -> 8 
// boxtype := buf[4:8] -> "free" 

// mdat
// []uint8 len: 8, cap: 8, [13,185,115,113,109,100,97,116]
// size := buf[0:4] -> 230257521
// boxtype := buf[4:8] -> mdat
