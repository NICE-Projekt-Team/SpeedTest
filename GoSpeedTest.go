package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/Azer0s/quacktors"
	"github.com/Azer0s/quacktors/pid"
)

type sizeInputStruct struct {
	sender pid.Pid
	wg     sync.WaitGroup
	size   int
}

type arrayInputStruct struct {
	sender    pid.Pid
	size      int
	wg        sync.WaitGroup
	byteArray []byte
}

type arrayOutputStruct struct {
	sender    pid.Pid
	byteArray []byte
}

func writeArray() {
	sizeInput := quacktors.Receive()
	sender := sizeInput.(sizeInputStruct).sender
	size := sizeInput.(sizeInputStruct).size
	wg := sizeInput.(sizeInputStruct).wg
	var byteArray = make([]byte, size) //creating the array
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)          //defining the random number generator
	for i := 0; i < size; i++ { //looping through the array
		byteArray[i] = byte(r1.Intn(2)) //setting a random value
	}
	quacktors.Send(sender, arrayOutputStruct{
		sender:    quacktors.Self(),
		byteArray: byteArray,
	})
	wg.Done() //telling the waitgroup that the routine is finished
}

func transfer() {
	arrayInput := quacktors.Receive()
	sender := arrayInput.(arrayInputStruct).sender
	size := arrayInput.(arrayInputStruct).size
	wg := arrayInput.(arrayInputStruct).wg
	inputByteArray := arrayInput.(arrayInputStruct).byteArray
	var byteArray = make([]byte, size)
	for i := 0; i < size; i++ { //looping through the array
		byteArray[i] = inputByteArray[i] //transferring the value
	}
	quacktors.Send(sender, arrayOutputStruct{
		sender:    quacktors.Self(),
		byteArray: byteArray,
	})
	wg.Done() //telling the waitgroup that the routine is finished
}

func sendPackages(size int) float64 {

	fmt.Println("Size of the Array is set to: ", size/8000000, " megabyte")

	var wg sync.WaitGroup //creating the waitgroup
	self := quacktors.Self()

	var byteArray1 = make([]byte, size)
	var byteArray2 = make([]byte, size) //Creating the 3 byteArrays
	var byteArray3 = make([]byte, size)

	fmt.Println("Starting to write Arrays")

	sizeInput := sizeInputStruct{self, wg, size}

	wg.Add(1)
	pid1 := quacktors.Spawn(writeArray)
	wg.Add(1)
	pid2 := quacktors.Spawn(writeArray)
	wg.Add(1)
	pid3 := quacktors.Spawn(writeArray)

	quacktors.Send(pid1, sizeInput)
	quacktors.Send(pid2, sizeInput)
	quacktors.Send(pid3, sizeInput)

	byteArray1 = quacktors.Receive().(arrayOutputStruct).byteArray
	byteArray2 = quacktors.Receive().(arrayOutputStruct).byteArray //receiving data from channels
	byteArray3 = quacktors.Receive().(arrayOutputStruct).byteArray

	wg.Wait() //waiting until all routines are finished

	//now transferring

	//TODO: implement the actor model framework for the transfer as well
	fmt.Println("Starting to transfer Arrays")
	startTine := time.Now()

	arrayInput1 := arrayInputStruct{self, size, wg, byteArray1}
	arrayInput2 := arrayInputStruct{self, size, wg, byteArray2}
	arrayInput3 := arrayInputStruct{self, size, wg, byteArray3}

	wg.Add(1) //add 1 task to the GoRoutine
	quacktors.Send(quacktors.Spawn(transfer), arrayInput1)
	wg.Add(1)
	quacktors.Send(quacktors.Spawn(transfer), arrayInput2)
	wg.Add(1)
	quacktors.Send(quacktors.Spawn(transfer), arrayInput3)

	byteArray3 = quacktors.Receive().(arrayOutputStruct).byteArray
	byteArray2 = quacktors.Receive().(arrayOutputStruct).byteArray //receiving data from channels
	byteArray1 = quacktors.Receive().(arrayOutputStruct).byteArray

	//wg.Wait() //waiting until all routines are done

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
