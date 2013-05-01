package redis

//Mutexes are useful when you need to make sure that two separate processes aren't using the same underlying resources
//But what do you do when the processes are an separate machines
//Redis can be used to facilitate the network-wide Mutex, and this is the interface I will be using
type Mutex interface {
	//Try will attempt to gain access to the mutex, but stop if it can't immediately gain access and return false.
	//If it succeeds, it executes the function and returns true
	Try(action func(resourceID int)) bool

	//Force will block until the mutex is available, and then execute the function
	Force(action func(resourceID int))
}

type redisMutex struct {
	init      String
	processes IntList
}

func newMutex(client Prefix, key string, count int) Mutex {
	m := new(redisMutex)
	m.init = client.String(key + ":Initialized")
	m.processes = client.IntList(key)
	m.initialize(count)

	return m
}

func (this *redisMutex) initialize(count int) {
	if <-this.init.Replace("initialized") != "initialized" {
		for i := 0; i < count; i++ {
			this.processes.RightPush(i) //give each one a unique resource ID
		}
	}
}

func (this *redisMutex) Try(action func(resourceID int)) bool {
	val, available := <-this.processes.LeftPop()
	if !available {
		return false
	}

	defer func() {
		<-this.processes.RightPush(val)
	}()

	action(val)

	return true
}

func (this *redisMutex) Force(action func(resourceID int)) {
	val := <-this.processes.BlockUntilLeftPop()

	defer func() {
		<-this.processes.RightPush(val)
	}()

	action(val)
}
