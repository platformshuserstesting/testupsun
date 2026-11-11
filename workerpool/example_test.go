package workerpool_test

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/platformshuserstesting/testupsun/workerpool"
)

func ExampleNew() {
	pool, err := workerpool.New(4, 8)
	if err != nil {
		log.Fatal(err)
	}

	var sum int64
	for i := 1; i <= 10; i++ {
		if err := pool.Submit(context.Background(), func() error {
			atomic.AddInt64(&sum, int64(i*i))
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}

	// Close waits for every submitted task to finish.
	if err := pool.Close(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("sum of squares:", atomic.LoadInt64(&sum))
	// Output: sum of squares: 385
}
