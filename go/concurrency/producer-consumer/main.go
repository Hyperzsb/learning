package main

//import spsc "concurrency/producerconsumer/single-producer-single-consumer"
import spmc "concurrency/producerconsumer/single-producer-multiple-consumer"

func main() {
	spmc.SingleProducerMultipleConsumer()
}
