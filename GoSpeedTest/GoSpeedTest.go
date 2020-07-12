package GoSpeedTest

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	var byteArray [800000000]byte //100 Megabyte
	byteArray2 := byteArray
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	var startTime time.Time
	var endTime time.Time

	fmt.Println("Starting to write Array")
	startTime = time.Now()
	for i := 0; i < len(byteArray); i++ {
		byteArray[i] = byte(r1.Intn(2))
	}
	endTime = time.Now()
	fmt.Println("Done writing Array. Duration: ", endTime.Sub(startTime), ", Speed: ", len(byteArray)/int(endTime.Sub(startTime).Milliseconds())/1000, "mBit")

	fmt.Println("Starting to transfer Array")
	startTime = time.Now()
	byteArray2 = byteArray //data is passed by value not reference
	endTime = time.Now()
	fmt.Println("Done transferring Array. Duration: ", endTime.Sub(startTime), ", Speed: ", len(byteArray)/int(endTime.Sub(startTime).Milliseconds())/1000, "mBit      ", byteArray2[1]) // Go forces me to use the byteArray2

}
