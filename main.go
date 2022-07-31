package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

//Producer is type for structs that holds two channel;one for pizzas, with all information for a given pizza order including wheather it was made successfully,
//and another to handle end of processing (when we quit the channel)
type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

// PizzaOrder is a type for structs that describe a given pizza order. It has the order number,a message indicating what happened to the order, and a boolean
//indicating if the order was successfully completed.
type PizzaOrder struct {
	pizzaNumber int
	message     string
	seccess     bool
}

//Close is simply a method of closing the channel when we are done with it. (i.e. something is pushed to the quit channel)
func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch

	return <-ch
}

// makePizza attempts to make a pizza, We generate a random number from 1-12, and put in two cases where we can't make the pizza in time.
//Otherwise, we make the pizza without issue. To make things intresting, each pizza will take a different length of time to produce (some pizzas are harder than others).
func makePizza(pizzaNumner int) *PizzaOrder {
	pizzaNumner++
	if pizzaNumner <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Receicved order #%d!\n", pizzaNumner)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizza #%d. It will take %d seconds.....\n", pizzaNumner, delay)

		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!", pizzaNumner)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d!", pizzaNumner)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumner)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumner,
			message:     msg,
			seccess:     success,
		}

		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumner,
	}

}

//pizzeria is a goroutine that runs in the background and calls makePizza to  try to make one order each time it iterates through the for loop.
//It executes until it receives something on the quit channel. The quit channel does not receive anything until the consumer sends it(When the number of
//orders is greather than or equal to the constant NumberOfPizza)

func pizzeria(pizzaMaker *Producer) {
	//keep track of which  pizza we are making
	var i = 0

	//run forever or until we receive a quit notification

	//try to make pizzas
	for {
		currentPizza := makePizza(i)
		//try to make pizza
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			//we try to male pizza(We sent something to the data channel)
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				//close channels
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
		//decision structure
	}

}
func main() {

	//seed random number generator
	rand.Seed(time.Now().UnixNano())

	//print out a message
	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	//create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}
	//run a producer in background
	go pizzeria(pizzaJob)

	//create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.seccess {
				color.Green(i.message)
				color.Green("Order #%d is our for delivery!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizzas....")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channel!", err)
			}
		}
	}
	//print out a ending message
	color.Cyan("-----------------")
	color.Cyan("Done for the day.")

	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		color.Red("It was an aweful day...")
	case pizzasFailed >= 6:
		color.Red("It was not a very good day...")
	case pizzasFailed >= 4:
		color.Yellow("It was an okay day...")
	case pizzasFailed >= 2:
		color.Yellow("It was an pretty good day!")
	default:
		color.Green("It was a great day.")
	}
}
