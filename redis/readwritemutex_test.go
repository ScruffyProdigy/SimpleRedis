package redis

import (
	"testing"
	"time"
)

const reader_max = 4

const read_attempts = 20
const reader_interval = 2
const reader_time = 10

const write_attempts = 2
const writer_interval = 5
const writer_time = 15

func TestReadWriteMutices(t *testing.T) {
	r := GetRedis(t)
	defer r.Close()

	rw := r.ReadWriteMutex("Test_ReadWriteMutex", reader_max)

	usage := 0
	done := make(chan bool)

	go func() {
		for i := 0; i < read_attempts; i++ {
			go func() {
				rw.Read.Force(func(resourceID int) {
					usage++

					if usage > reader_max {
						t.Error("Usage too hi")
					}

					time.Sleep(reader_time * time.Millisecond)

					usage--
				})
				done <- true
			}()
			time.Sleep(reader_interval * time.Millisecond)
		}
	}()
	for i := 0; i < write_attempts; i++ {
		time.Sleep(writer_interval * time.Millisecond)
		go func() {
			rw.Write.Force(func(resourceID int) {
				if usage > 0 {
					t.Error("Should be no other usage when write is happening")
				}

				time.Sleep(writer_time * time.Millisecond)
			})
			done <- true
		}()
	}

	for i := 0; i < read_attempts+write_attempts; i++ {
		<-done
	}
}
