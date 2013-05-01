package redis

//ReadWriteMutexes are useful for making sure nothing is trying to read data while you're trying to write to it
//When you're trying to read and write across a network, the mutex needs to work across the network too
//And redis works well for this
type ReadWriteMutex struct {
	readers int
	write   Mutex

	//Read is a semaphore that gives access to be able to read information and be sure nothing is trying to write
	Read Mutex

	//Write is a mutex that gives access to be able to write information such that nothing can read
	Write Mutex
}

type writeMutex struct {
	*ReadWriteMutex
}

func lockAllReads(rw *ReadWriteMutex, finalAction func(resourceID int)) func(int) {
	return func(resourceID int) {
		in := make(chan bool)
		out := make(chan bool)
		for i := 0; i < rw.readers; i++ {
			go func(j int) {
				rw.Read.Force(func(int) {
					in <- true
					<-out //this one is here to make sure the final action has been completed before releasing any of the readers
				})
				<-out //this one is here to make sure all of the readers have been released before returning control to the calling function
			}(i)
		}

		for i := 0; i < rw.readers; i++ {
			<-in
		}
		finalAction(resourceID)
		for i := 0; i < 2*rw.readers; i++ {
			out <- true
		}
	}
}

func (this writeMutex) Try(action func(i int)) bool {
	return this.write.Try(lockAllReads(this.ReadWriteMutex, action))
}

func (this writeMutex) Force(action func(i int)) {
	this.write.Force(lockAllReads(this.ReadWriteMutex, action))
}

func newRWMutex(client Prefix, key string, readers int) *ReadWriteMutex {
	rw := new(ReadWriteMutex)
	rw.readers = readers
	rw.write = client.Mutex(key + ":Write")
	rw.Read = client.Semaphore(key+":Read", readers)
	rw.Write = writeMutex{rw}
	return rw
}
