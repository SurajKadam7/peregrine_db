package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randomeNumber(from, to int) int {
	// Generate a random number between from and to
	randomNumber := rand.Intn(to-from+1) + from
	return randomNumber

}

var iterationRequireTofind = make(map[int]int)

func main() {
	var maxElement uint16 = 7
	t := New(maxElement)

	keys := make(map[int]bool)
	// keyArray := make([]int, 0)
	start := time.Now()
	for i := 2; i < 1200000; i += 1 {

		key := randomeNumber(10, 200000000)
		value := float64(key)
		value = value / 10
		if keys[key] {
			continue
		}
		keys[key] = true
		// keyArray = append(keyArray, key)
		t.set(key, value)
	}

	fmt.Println("total insertion time : ", time.Since(start))

	// fmt.Println("inserted numbers : ", keyArray)
	fmt.Println("inserted numbers : ", len(keys))

	start = time.Now()
	cnt := 0
	wrong := 0

	for key := range keys {
		value := t.get(key)
		if value == 0 {
			continue
		}
		if value != float64(key)/10 {
			wrong += 1
			continue
		}
		cnt++
		// fmt.Println(key, value)
	}

	fmt.Println("total finding time : ", time.Since(start))
	mx := 0
	ke := 0
	for key, value := range iterationRequireTofind {
		mx = max(mx, value)
		if mx == value {
			ke = key
		}
	}

	fmt.Println("find numbers : ", cnt)
	fmt.Println("wrong result : ", wrong)
	fmt.Println("max iterations : ", mx, ke)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
