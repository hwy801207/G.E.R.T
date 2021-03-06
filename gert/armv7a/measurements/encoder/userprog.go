package main

import (
	"./embedded"
	"fmt"
	"math"
	"time"
	//	"unsafe"
)

var event_chan chan interface{}
var drive *embedded.MDD10A_controller
var adc *embedded.MCP3008_controller

func user_init() {

	//play with the SD card a bit
	//	good, root := embedded.Fat32_som_start(embedded.Init_som_sdcard, embedded.Read_som_sdcard)
	//	if !good {
	//		fmt.Println("fat32 init failure")
	//	}
	//	fmt.Println(root.Getfilenames())
	//	fmt.Println(root.Getsubdirnames())
	//	good, bootdir := root.Direnter("BOOT")
	//	if !good {
	//		panic("dir entry failed")
	//	} else {
	//		fmt.Println(bootdir.Getfilenames())
	//		good, contents := bootdir.Fileread("UENV.TXT")
	//		if !good {
	//			panic("file read failure")
	//		}
	//		fmt.Println(string(contents))
	//	}
	adc = embedded.MakeMCP3008(embedded.WB_SPI1)
	drive = embedded.MakeMDD10A(embedded.WB_PWM1, embedded.WB_PWM2, embedded.WB_JP4_4, embedded.WB_JP4_6)
	event_chan = make(chan interface{}, 10)
	go func() {
		for {
			event_chan <- string(embedded.WB_DEFAULT_UART.Read(1)[:])
		}
	}()

	go func() {
		for {
			event_chan <- adc.Read(0)
			time.Sleep(2 * time.Second)
		}
	}()
	fmt.Printf("pi is %v \n", pi(50))

	go func() {
		for {
			old := count
			time.Sleep(1 * time.Second)
			new := count
			event_chan <- new - old
		}
	}()

	//	go func() {
	//		for {
	//			time.Sleep(100 * time.Millisecond)
	//			event_chan <- "Poll fast!"
	//		}
	//	}()
	//fmt.Printf("press key to continue\n")
	//embedded.WB_DEFAULT_UART.Read(1)
	//fmt.Printf("wait 10 sec... \n")
	//time.Sleep(10 * time.Second)
	//embedded.WB_JP4_10.SetOutput()
	//embedded.WB_JP4_10.SetLO()
	embedded.WB_JP4_10.SetInput()
	embedded.WB_JP4_10.EnableIntr(embedded.INTR_FALLING, inc)
	embedded.Enable_interrupt(99, 0) //send GPIO1 interrupt to CPU0
	//embedded.Enable_interrupt(103, 0) //send GPIO3 interrupt to CPU0

	//	embedded.WB_PWM1.Begin(0x10)
	//	embedded.WB_PWM1.SetDuty(0.5)
	//
	//	embedded.WB_PWM2.Begin(0xF000)
	//	embedded.WB_PWM2.SetDuty(0.5)
	//
	//	embedded.WB_PWM3.Begin(0xFF00)
	//	embedded.WB_PWM3.SetDuty(0.5)

	//	embedded.WB_SPI1.Begin(0, 10, 8, 0)

	//send the GPT interrupt to CPU1
	//embedded.Enable_interrupt(87, 1)

	//start the GPT
	//embedded.StartGPT()

	//start a little watch
	//go embedded.Gopherwatch()
}

func user_loop() {
	//make an event loop
	select {
	case event := <-event_chan:
		//if event != oldevent {
		fmt.Printf("%v\n", event)
		switch event {
		case "p":
			//embedded.WB_SPI1.Send(0xAA)
			val := adc.Read(0)
			fmt.Printf("adc reads %v\n", val)
		case "w":
			drive.Forward(0.2)
		case "s":
			drive.Backward(0.2)
		case "a":
			drive.TurnRight(0.2)
		case "d":
			drive.TurnLeft(0.2)
		case " ":
			drive.Stop()
		}
		//}
		//oldevent = event
		//	default:
		//drive.Stop()
	}

	//	fmt.Printf("waiting for input: ")
	//	data := string(embedded.WB_DEFAULT_UART.Read(10)[:])
	//	fmt.Printf("got %s\n", data)
	//embedded.WB_JP4_6.SetLO()
	//embedded.WB_JP4_6.SetHI()
	//	embedded.WB_JP4_4.SetLO()
	//	embedded.WB_JP4_4.SetHI()
	//	if ping {
	//		fmt.Printf("count is %d\n", count)
	//		ping = false
	//	}
	//embedded.Sleep(2)
	//fmt.Printf("count is %d\n", count)
}

var count uint32

func inc() {
	count += 1
}

// pi launches n goroutines to compute an
// approximation of pi.
func pi(n int) float64 {
	ch := make(chan float64)
	for k := 0; k <= n; k++ {
		go term(ch, float64(k))
	}
	f := 0.0
	for k := 0; k <= n; k++ {
		f += <-ch
	}
	return f
}

func term(ch chan float64, k float64) {
	ch <- 4 * math.Pow(-1, k) / (2*k + 1)
}
