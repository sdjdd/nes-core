package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"

	"fyne.io/fyne/app"
)

type logLine struct {
	pc     uint16
	opcode byte
	num    [2]byte
}

func main() {
	app := app.New()

	w := app.NewWindow("Hello")
	w.SetContent(fyne.NewContainer(
		canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
			return color.RGBA{0, 0x80, 0, 0xff}
		}),
	))
	w.ShowAndRun()
}

func main2() {
	logFile, _ := os.Open("nestest.log")
	defer logFile.Close()
	reader := bufio.NewReader(logFile)

	console := new(Console)
	cart, err := loadRomFile("nestest.nes")
	if err != nil {
		log.Fatalf("open rom file: %s", err)
	}

	cpu := new(CPU)
	console.Connect(cpu)
	console.Connect(cart)

	cpu.PC = 0xC000
	cpu.cycles = 7
	for i := 1; ; i++ {
		mycyc := cpu.cycles
		fmt.Printf("%-5d ", i)
		info := cpu.step()
		var opnumsText string
		switch instructionSizes[info.ins.addrMode] {
		case 3:
			opnumsText = fmt.Sprintf("%02X %02X", info.opnums[0], info.opnums[1])
		case 2:
			opnumsText = fmt.Sprintf("%02X", info.opnums[0])
		}
		var addrText string
		if info.ins.addrMode == addrRelative {
			addrText = fmt.Sprintf("*%+X", info.opnums[0])
		} else {
			switch instructionSizes[info.ins.addrMode] {
			case 3:
				addrText = fmt.Sprintf("$%02X%02X", info.opnums[1], info.opnums[0])
			case 2:
				addrText = fmt.Sprintf("$%02X", info.opnums[0])
			}
			if info.ins.addrMode == addrImmediate {
				addrText = "#" + addrText
			}
		}
		l := fmt.Sprintf("%04X  %02X %-5s  %s %-5s", info.PC, info.opcode, opnumsText, instructionNames[info.ins.id], addrText)
		line, _, _ := reader.ReadLine()
		fmt.Print(l)
		fmt.Printf(" A:%02X X:%02X Y:%02X S:%02X P:%02X Z:%d C:%d CYC:%d", cpu.A, cpu.X, cpu.Y, cpu.S, cpu.flag(), cpu.Z, cpu.C, mycyc)
		fmt.Printf("|  %s", string(line))

		fmt.Println()
	}
}
