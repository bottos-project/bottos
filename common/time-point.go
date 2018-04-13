package dpos

import(
	"time"
)
type Microseconds struct{
	int64 micro
}
func NowToSeconds() int64 {
	return time.Now().Unix()
}
func microseconds() int64{

}
func millseconds(s int64) int64{
	return s * 1000
}
func seconds(s int64) int64{
	return  s * 1000000
}


func getEpochTime() int64{
	now := time.Now() // get current time
	epoch := time.Since(now)
	return epoch.Nanoseconds*1000 // microseconds
}
  