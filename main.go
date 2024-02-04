package main
import (
 "context"
 "fmt"
 "sync"
 "time"
)
var wakeBell chan bool
var sleepBell chan bool
func sleepBarber(name string) {
 sleepBell <- true
 fmt.Println("=========== Barber went for sleep")
}
func wakeBarber() {
 fmt.Println("=========== Wakeup barber")
}
func barberCutting(name string) {
 for i := 0; i < 5; i++ {
  fmt.Println("=========== Cutting in progress for ", name)
  time.Sleep(1 * time.Second)
 }
 fmt.Println("=========== Cutting completed for ", name)
}
func barberShop(chairs chan string, ctx context.Context, wg *sync.WaitGroup) {
 wg.Add(1)
 sleepBell <- true
 defer wg.Done()
 fmt.Println("=========== Barber shop opened")
 defer fmt.Println("=========== Barber shop closed")
 closeShop := false
 for {
  select {
  case name := <-chairs:
   barberCutting(name)
   if len(chairs) == 0 {
    if !closeShop {
     sleepBarber(name)
    } else {
     return
    }
   }
  case <-wakeBell:
   wakeBarber()
  case <-ctx.Done():
   if len(chairs) == 0 {
    return
   }
   closeShop = true
  }
}
}
func customerEntry(name string, chairs chan string) {
 fmt.Printf("=========== Customer %s entered Barber shop\n", name)
 select {
case <-sleepBell:
  wakeBell <- true
  chairs <- name
  fmt.Printf("=========== Got a chair for %s\n", name)
 case chairs <- name:
  fmt.Printf("=========== Got a chair for %s\n", name)
 default:
  fmt.Printf("=========== No chair for %s, hence leaving Barber shop\n", name)
 }
}
func main() {
 wakeBell = make(chan bool)
 sleepBell = make(chan bool, 1)
 chairs := make(chan string, 2)
 wg := new(sync.WaitGroup)
 ctx, cancel := context.WithCancel(context.Background())
 go barberShop(chairs, ctx, wg)
 <-sleepBell
 customerEntry("Anil", chairs)
 customerEntry("Dravid", chairs)
 customerEntry("Sachin", chairs)
 customerEntry("Sourav", chairs)
 time.Sleep(20 * time.Second)
 customerEntry("Dhoni", chairs)
 cancel()
 wg.Wait()
 close(wakeBell)
 close(chairs)
 fmt.Println("=========== Exiting...")
}
