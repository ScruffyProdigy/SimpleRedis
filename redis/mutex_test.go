package redis

import (
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	m := r.Mutex("Test_Mutex")

	usage := 0
	j := 0

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					if rec.(int) != 5 {
						panic(rec)
					}
				}
			}()
			m.Force(func(int) {
				usage++

				if usage > 1 {
					t.Error("Usage level too hi:", usage)
				}

				time.Sleep(100 * time.Millisecond)

				usage--
				done <- true

				j++
				if j%2 != 0 {
					panic(5) // make sure panic's don't eat up resources
				}
			})
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if usage != 0 {
		t.Error("Usage should be finished")
	}

	for i := 0; i < 10; i++ {
		go func() {
			worked := m.Try(func(int) {
				usage++

				if usage > 1 {
					t.Error("Usage level too hi:", usage)
				}

				time.Sleep(100 * time.Millisecond)

				usage--
			})

			done <- worked
		}()
		time.Sleep(50 * time.Millisecond)
	}

	completed := 0

	for i := 0; i < 10; i++ {
		if <-done {
			completed++
		}
	}

	t.Log(completed, "mutex goroutine(s) successful")
}

func TestSemaphore(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	m := r.Semaphore("Test_Semaphore", 3)

	usage := 0
	j := 0

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			defer func() {
				if rec := recover(); rec != nil {
					if rec.(int) != 5 {
						panic(rec)
					}
				}
			}()
			m.Force(func(int) {
				usage++

				if usage > 3 {
					t.Error("Usage level too hi:", usage)
				}

				time.Sleep(100 * time.Millisecond)

				usage--
				done <- true

				j++
				if j%2 != 0 {
					panic(5) // make sure panic's don't eat up resources
				}

			})
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if usage != 0 {
		t.Error("Usage should be finished")
	}

	for i := 0; i < 10; i++ {
		go func() {
			worked := m.Try(func(int) {
				usage++

				if usage > 3 {
					t.Error("Usage level too hi:", usage)
				}

				time.Sleep(100 * time.Millisecond)

				usage--

			})

			done <- worked
		}()
		time.Sleep(25 * time.Millisecond)
	}

	completed := 0

	for i := 0; i < 10; i++ {
		if <-done {
			completed++
		}
	}

	t.Log(completed, "semaphore goroutines successful")
}
