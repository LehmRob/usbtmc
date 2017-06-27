// Copyright (c) 2015-2017 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gotmc/usbtmc"
	_ "github.com/gotmc/usbtmc/driver/truveris"
)

func main() {
	flag.Parse()

	start := time.Now()
	log.Println("Create new USBTMC context")
	ctx, err := usbtmc.NewContext()
	log.Println("Created new USBTMC context")
	if err != nil {
		log.Fatalf("Error creating new USB context: %s", err)
	}
	defer ctx.Close()

	start = time.Now()
	log.Println("Here")
	fg, err := ctx.NewDeviceByVIDPID(0x957, 0x407)
	log.Println("Boo")
	if err != nil {
		log.Fatalf("NewDevice error: %s", err)
	}
	defer fg.Close()
	log.Printf("%.2fs to setup instrument\n", time.Since(start).Seconds())
	start = time.Now()
	// fmt.Printf(
	// "Found the Arb Wave Gen S/N %s by Vendor ID %d with Product ID %d\n",
	// fg.USBDevice.Descriptor.SerialNumber,
	// fg.USBDevice.Descriptor.Vendor,
	// fg.USBDevice.Descriptor.Product)
	// Send commands to waveform generator
	fg.Write([]byte("apply:sinusoid 2340, 0.1, 0.0")) // Write using byte slice
	io.WriteString(fg, "burst:internal:period 0.112") // WriteString using io's Writer interface
	fg.WriteString("burst:internal:period 0.112")     // WriteString
	fg.WriteString("burst:ncycles 131")
	fg.WriteString("burst:state on")
	fg.WriteString("*idn?")

	start = time.Now()
	var buf [1024]byte
	bytesRead, err := fg.Read(buf[:])
	log.Printf("%.2fs to read %d bytes\n", time.Since(start).Seconds(), bytesRead)
	if err != nil {
		log.Printf("Error reading: %s", err)
	}
	fmt.Printf("Read %d bytes = %s", bytesRead, buf[12:bytesRead])
	fmt.Printf("Last rune read = %x\n", buf[bytesRead-1:bytesRead])
	fmt.Printf("Last rune read = %q\n", buf[bytesRead-1:bytesRead])
	fmt.Printf("Read %d bytes = %v\n", bytesRead, buf[:12])

	// log.Print(fg.Write("freq 2340"))
	// log.Print(scope.Ask("*idn?"))
	defer fg.Close()
	// fmt.Printf("Goodbye arbitrary waveform generator %s\n", fg.USBDevice.Descriptor.SerialNumber)

}
