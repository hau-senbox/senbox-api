package queue

func New() chan func() {
	var queue = make(chan func())

	go func() {
		for true {
			nextFunction := <-queue
			nextFunction()
		}
	}()

	return queue
}
