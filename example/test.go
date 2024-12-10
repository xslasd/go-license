package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"os"
)

func main() {
	fmt.Println("---------------------------------------------")
	fmt.Println(host.Info())
	fmt.Println("---------------------------------------------")
	fmt.Println(cpu.Info())
	fmt.Println("---------------------------------------------")
	fmt.Println(os.Getwd())
}
