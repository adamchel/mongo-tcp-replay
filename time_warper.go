package main

import (
	"fmt"
	"io/ioutil"
	"github.com/smallfish/simpleyaml"
	"strconv"
)

func ReadConfig() *simpleyaml.Yaml {
	dat, err := ioutil.ReadFile("warp_config.yml")
	if err != nil {
		fmt.Println("Error reading warp config file")
	}
	y, err := simpleyaml.NewYaml(dat)
	if err != nil {
		fmt.Println("Error parsing warp config file as YAML")
	}
	return y
}

func ApplyTransformations(time uint64) uint64 {
	yamlData := ReadConfig()
	multStr, _ := yamlData.Get("multiply").String()
	mult, _ := strconv.ParseFloat(multStr, 64)
	time = TimeMultiply(time, mult)
	shift, _ := yamlData.Get("shift").Int()
	time = TimeShift(time, shift)
	return time
}

func TimeMultiply(time uint64,
				  multiplier float64) uint64{
	return uint64(float64(time) * multiplier)
}

func TimeShift(time uint64,
			   shift int) uint64 {
	return uint64(int(time) + shift)
}