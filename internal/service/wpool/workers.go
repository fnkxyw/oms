package wpool

func ChangeNumOfWorkers(wp *WorkerPool, numOfWorkers int) error {

	if numOfWorkers > 0 {
		return wp.AddWorker(numOfWorkers)
	}
	if numOfWorkers < 0 {
		numOfWorkers *= -1
		return wp.RemoveWorkers(numOfWorkers)

	}
	return nil
}
