package Counter

import (
	"log"
	"strconv"
	"time"
)

type (
	Counter struct {
		duration int
		state    int // 0 running & 1 paused
		timer    chan int
		quit     chan bool
	}
)

func NewCounter(duration int) *Counter {
	command := make(chan int)
	quit := make(chan bool)
	return &Counter{duration, 0, command, quit}
}

func (c *Counter) StartTicker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	count := (*c).duration

	for {
		select {
		case <-ticker.C:
			if c.state == 0 {
				count--
			}
			if count == 0 {
				c.StopCounter()
				return
			} else {
				c.timer <- count
			}
		case <- c.quit:
			return
		}
	}
}

func (c *Counter) StartCounter(out chan string) {
	c.ResumeCounter()
	go c.StartTicker()
	for {
		select {
		case t := <-c.timer:
			out <- strconv.Itoa(c.state) + "-" + strconv.Itoa(t)
		case <-c.quit:
			out <- "done"
			return
		}
	}
}

func (c *Counter) PauseCounter() {
	c.state = 1
	log.Println("Pause boy")
}

func (c *Counter) ResumeCounter() {
	c.state = 0
	log.Println("Resume boy")
}

func (c *Counter) StopCounter() {
	c.quit <- true
	log.Println("Stop boy")
}

func (c *Counter) HandleCommand(command string, outChannel chan string) {
	switch command {
	case "start":
		go c.StartCounter(outChannel)
	case "pause":
		c.PauseCounter()
	case "resume":
		c.ResumeCounter()
	case "stop":
		c.StopCounter()
	}
}