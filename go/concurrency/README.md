# Concurrency

As concurrency is natively supported by Go with the help of `goroutinue`, it is essential to learn more about how to write idiomatic concurrent code in Go.

However, with few opportunities to practice the concurrent programming in simple programming exercises and projects, I may need some guides and inspiration to learn further about this topic.

## Classical Synchronization Problems

There are several famous synchronization problems proposed by [Dijkstra](https://en.wikipedia.org/wiki/Edsger_W._Dijkstra) and other researchers, which are widely-used to test different synchronization schemes.
Getting familiar with these problems is a great start to dive into the concurrent programming.

- [Producer Consumer (Bonded Buffer) Problem](producer-consumer): This problem mainly focuses on the usage of the semaphore and mutex to keep track of the current number of full and empty buffers respectively, but it can also be solved by the channel. [*Wikipedia*](https://en.wikipedia.org/wiki/Producer%E2%80%93consumer_problem#See_also)
- [Dining Philosophers Problem](dining-philosophers): This problem states that K philosophers seated around a circular table with one fork between each pair of philosophers. 
  A philosopher may eat if he can pick up two forks adjacent to him. 
  One fork may be picked up by any one of its adjacent followers but not both.
  This problem involves the allocation of limited resources to a group of processes in a deadlock-free and starvation-free manner. [*Wikipedia*](https://en.wikipedia.org/wiki/Dining_philosophers_problem)
- [Readers Writers Problem](readers-writers): Some threads may read and some may write, with the constraint that no thread may access the shared resource for either reading or writing while another thread is in the act of writing to it. 
  (In particular, we want to prevent more than one thread modifying the shared resource simultaneously and allow for two or more readers to access the shared resource at the same time). [*Wikipedia*](https://en.wikipedia.org/wiki/Readers%E2%80%93writers_problem)
- [Sleeping Barber Problem](sleeping-barber): Barber shop with one barber, one barber chair and N chairs to wait in. 
  When no customers the barber goes to sleep in barber chair and must be woken when a customer comes in. 
  When barber is cutting hair new customers take empty seats to wait, or leave if no vacancy. [Wikipedia](https://en.wikipedia.org/wiki/Sleeping_barber_problem)
