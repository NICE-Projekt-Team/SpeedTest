package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func writeArray(size int, wg *sync.WaitGroup, result chan<- []byte) {
	fmt.Println("Channel opened, Starting to write Array")
	var byteArray = make([]byte, size) //creating the array
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)          //defining the random number generator
	for i := 0; i < size; i++ { //looping through the array
		byteArray[i] = byte(r1.Intn(2)) //setting a random value
	}
	fmt.Println("Done writing array")
	result <- byteArray //sending data back through channel
	defer wg.Done()     //telling the waitgroup that the routine is finished
}

func transfer(size int, wg *sync.WaitGroup, input []byte, result chan<- []byte) {
	fmt.Println("Channel opened, starting to transfer array")
	var byteArray = make([]byte, size)
	for i := 0; i < size; i++ { //looping through the array
		byteArray[i] = input[i] //transferring the value
	}
	result <- byteArray //sending data back through channel
	defer wg.Done()     //telling the waitgroup that the routine is finished

}

func main() {
	var size = 80000000 //Setting the size of the Arrays
	fmt.Println("Size of the Array is set to: ", size/8000000, " megabyte")

	var wg sync.WaitGroup //creating the waitgroup

	var byteArray1 = make([]byte, size)
	var byteArray2 = make([]byte, size) //Creating the 3 byteArrays
	var byteArray3 = make([]byte, size)

	var startTime time.Time //Creating the start and end Time variable
	var endTime time.Time

	ch1 := make(chan []byte, 1)
	ch2 := make(chan []byte, 1) //Creating the 3 byteArray1-Channels
	ch3 := make(chan []byte, 1)

	fmt.Println("Starting to write Arrays")
	startTime = time.Now() //Starting the timer

	wg.Add(1)                     //add 1 task to the GoRoutine
	go writeArray(size, &wg, ch1) //define goroutine
	wg.Add(1)
	go writeArray(size, &wg, ch2)
	wg.Add(1)
	go writeArray(size, &wg, ch3)

	byteArray1 = <-ch1
	byteArray2 = <-ch2 //receiving data from channels
	byteArray3 = <-ch3

	wg.Wait() //waiting until all routines are finished

	close(ch1)
	close(ch2) //closing channels
	close(ch3)

	endTime = time.Now()                  //setting end time
	var duration = endTime.Sub(startTime) //calculating duration
	fmt.Println("Done writing Arrays with routines. Duration: ", duration, ", Speed:", (float64(size/1000000)/(duration.Seconds()+(float64(duration.Milliseconds())/1000)))*3, "mBit/s")

	//now transferring

	fmt.Println("Starting to transfer Arrays")
	startTime = time.Now() //setting starttime

	ch4 := make(chan []byte, 1)
	ch5 := make(chan []byte, 1) //Creating the 3 byteArray-Channels
	ch6 := make(chan []byte, 1)

	wg.Add(1)                               //add 1 task to the GoRoutine
	go transfer(size, &wg, byteArray1, ch4) //define goroutine
	wg.Add(1)
	go transfer(size, &wg, byteArray2, ch5)
	wg.Add(1)
	go transfer(size, &wg, byteArray3, ch6)

	byteArray1 = <-ch1
	byteArray2 = <-ch2 //receiving data from channels
	byteArray3 = <-ch3

	wg.Wait() //waiting until all routines are done

	close(ch4)
	close(ch5) //closing channels
	close(ch6)

	endTime = time.Now()              //stopping time
	duration = endTime.Sub(startTime) //calculating duration
	fmt.Println("Done transferring Array. Duration: ", duration, ", Speed:", (float64(size/1000000)/(duration.Seconds()+(float64(duration.Milliseconds())/1000)))*3, "mBit/s")
}
