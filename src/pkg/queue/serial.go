package queue

func New() chan func() {
	var queue = make(chan func())

	go func() {
		for {
			nextFunction := <-queue
			nextFunction()
		}
	}()

	return queue
}
