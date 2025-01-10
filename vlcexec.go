package main

import (
	"fmt"
	"log"
	"os/exec"
)


func playMixExec(fileName string) {
    fmt.Printf("playing now: %s\n", fileName)
    cmd := exec.Command("vlc", fileName, "--sub-filter=marq{marquee=test,color=0x00FF00,position=8}")

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
