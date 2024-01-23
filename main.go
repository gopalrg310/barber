package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numBarbers        = 2
	numWaitingChairs  = 5
	openingTime       = 8 * time.Hour
	closingTime       = 17 * time.Hour
	haircutDuration   = 30 * time.Minute
	clientArrivalTime = 10 * time.Minute
)

type BarberShop struct {
	mu              sync.Mutex
	waitingChairs   chan bool
	barberAvailable chan bool
}

func NewBarberShop() *BarberShop {
	return &BarberShop{
		waitingChairs:   make(chan bool, numWaitingChairs),
		barberAvailable: make(chan bool, numBarbers),
	}
}

func (bs *BarberShop) openShop() {
	closeTime := time.Now().Add(closingTime)

	for {
		select {
		case <-time.After(clientArrivalTime):
			bs.handleClient()
		case <-time.After(closeTime.Sub(time.Now())):
			bs.waitForEmptyWaitingRoom()
			fmt.Println("Barber shop is closing.")
			return
		}
	}
}

func (bs *BarberShop) handleClient() {
	bs.mu.Lock()
	select {
	case bs.barberAvailable <- true:
		// There is an available barber
		bs.mu.Unlock()
		go bs.cutHair()
	case bs.waitingChairs <- true:
		// No available barber, but there's a waiting chair
		bs.mu.Unlock()
		fmt.Println("Client takes a seat.")
	}
}

func (bs *BarberShop) cutHair() {
	fmt.Println("Barber starts cutting hair.")
	time.Sleep(haircutDuration)
	fmt.Println("Barber finishes cutting hair.")
	<-bs.barberAvailable
	bs.mu.Lock()
	if len(bs.waitingChairs) > 0 {
		// There is a waiting customer, wake up the barber
		bs.mu.Unlock()
		<-bs.waitingChairs
		bs.barberAvailable <- true
	} else {
		bs.mu.Unlock()
	}
}

func (bs *BarberShop) waitForEmptyWaitingRoom() {
	fmt.Println("Barber shop is closing. Waiting for the waiting room to be empty.")
	for len(bs.waitingChairs) > 0 {
		<-bs.waitingChairs
	}
}

func main() {
	barberShop := NewBarberShop()

	// Start the barbers
	for i := 0; i < numBarbers; i++ {
		go func() {
			for {
				<-barberShop.barberAvailable
				barberShop.cutHair()
			}
		}()
	}

	// Start the barber shop
	barberShop.openShop()
}
