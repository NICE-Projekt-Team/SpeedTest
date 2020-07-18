package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Azer0s/quacktors"
	"github.com/Azer0s/quacktors/pid"
)

type sizeInput struct {
	sender pid.Pid
	wg     sync.WaitGroup
	size   int
}

func writeArray() {
	var sizeInput = quacktors.Receive()
	sender := sizeInput(sender)
	size := sizeInput(size)
	wg := sizeInput(wg)
	var byteArray = make([]byte, size) //creating the array
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)          //defining the random number generator
	for i := 0; i < size; i++ { //looping through the array
		byteArray[i] = byte(r1.Intn(2)) //setting a random value
	}
	result <- byteArray //sending data back through channel
	defer wg.Done()     //telling the waitgroup that the routine is finished
}

func transfer(size int, wg *sync.WaitGroup, input []byte, result chan<- []byte) {
	var byteArray = make([]byte, size)
	for i := 0; i < size; i++ { //looping through the array
		byteArray[i] = input[i] //transferring the value
	}
	result <- byteArray //sending data back through channel
	defer wg.Done()     //telling the waitgroup that the routine is finished
}

func sendPackages(size int) float64 {

	fmt.Println("Size of the Array is set to: ", size/8000000, " megabyte")

	var wg sync.WaitGroup //creating the waitgroup
	self := quacktors.Self()

	var byteArray1 = make([]byte, size)
	var byteArray2 = make([]byte, size) //Creating the 3 byteArrays
	var byteArray3 = make([]byte, size)

	fmt.Println("Starting to write Arrays")

	sizeInput := sizeInput{self, wg, size}

	wg.Add(1) //add 1 task to the GoRoutine
	quacktors.Spawn(writeArray)
	quacktors.Send(self, sizeInput)
	wg.Add(1)
	quacktors.Spawn(writeArray)
	quacktors.Send(self, sizeInput)
	wg.Add(1)
	quacktors.Spawn(writeArray)
	quacktors.Send(self, sizeInput)

	byteArray1 = quacktors.Receive()
	byteArray2 = quacktors.Receive() //receiving data from channels
	byteArray3 = quacktors.Receive()

	wg.Wait() //waiting until all routines are finished

	//now transferring

	fmt.Println("Starting to transfer Arrays")
	startTine := time.Now()
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

	speed := float64(size/1000000) / (float64(time.Now().Sub(startTine).Milliseconds()) / 1000)
	fmt.Println("Done transferring Array. Duration: ", time.Now().Sub(startTine), ", Speed:", speed, "mBit/s")
	return speed
}

func main() {
	var size = 80000000     //Setting the size of the Arrays
	var duration = 1        //duration in minutes
	var startTime time.Time //Creating the start and end Time variable
	startTime = time.Now()
	var speeds = make([]float64, 2000)
	sum := float64(0)
	var counter = 1
	fmt.Println("time", time.Now().Sub(startTime))
	for time.Now().Sub(startTime).Minutes() < float64(duration) {

		speeds[counter-1] = sendPackages(size)
		counter++
	}

	for i := 0; i < counter; i++ {
		sum += speeds[i]
	}

	fmt.Println("3-Way transfer was executed ", counter, " times in ", duration, "minutes. The average speed was ", sum/float64(counter), "mBit")

}
