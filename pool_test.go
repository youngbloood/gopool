package gopool_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/youngbloood/gopool"
)

func run(v interface{}) error {
	fmt.Printf("v = %+v\n", v)
	time.Sleep(10 * time.Second)
	return nil
}

func _log(index string, t *testing.T, err error, i int) {
	t.Logf("%s%s=%d\n", t.Name(), index, i)
	if err != nil {
		t.Logf("%s%s.err=%v\n", t.Name(), index, err)
	}
}

func TestGoPool(t *testing.T) {
	pool := gopool.New(0, run)
	pool.StartGo()
	pool.Expand(1)
	_log("0", t, pool.Send("aa"), pool.Cap()) // cap is 1

	pool.Expand(1)
	_log("1", t, pool.Send("bb"), pool.Cap()) // cap is 2

	pool.Expand(-2)
	_log("2", t, pool.Send("cc"), pool.Cap()) // cap is 0

	pool.Expand(-1)
	pool.Expand(-1)
	time.Sleep(3 * time.Second)
	_log("3", t, pool.Send("dd"), pool.Cap()) // cap is 0

	pool.Expand(1)
	time.Sleep(3 * time.Second)               // wait the goroutine started!
	_log("4", t, pool.Send("ee"), pool.Cap()) // cap is 1
}

func TestGoPoolAdd(t *testing.T) {
	pool := gopool.New(5, run)
	_log("1", t, nil, pool.Cap()) // cap is 5

	pool.StartGo()
	_log("2", t, nil, pool.Cap()) // cap is 5
	pool.Expand(2)

	time.Sleep(3 * time.Second)
	_log("3", t, nil, pool.Cap()) // cap is 7
}

func TestGoPoolDone(t *testing.T) {
	pool := gopool.New(5, run)
	_log("0", t, nil, pool.Cap()) // cap is 5
	pool.StartGo()
	_log("1", t, nil, pool.Cap()) // cap is 5
	pool.Expand(-2)
	time.Sleep(3 * time.Second)
	_log("2", t, nil, pool.Cap()) // cap is 3
}

func TestGoPoolZero(t *testing.T) {
	pool := gopool.New(0, run)
	_log("0", t, nil, pool.Cap()) // cap is 0
	pool.StartGo()

	_log("1", t, pool.Send(1), pool.Cap()) // cap is 0
	pool.Expand(1)
	time.Sleep(3 * time.Second)
	_log("3", t, nil, pool.Cap()) // cap is 1
}

func TestGoPoolSend(t *testing.T) {
	pool := gopool.New(5, run) // cap is 5
	// go cycle(pool)
	check(pool.Send("111111"))
	check(pool.Send("222222"))
	check(pool.Send("333333"))
	check(pool.Send("444444"))
	check(pool.Send("555555"))
	check(pool.Send("666666"))
	check(pool.Send("777777"))
	check(pool.Send("888888"))
	fmt.Println("增加容量")
	pool.Expand(2) // cap is 7
	// stop(1 * time.Second)
	check(pool.Send("999999"))
	check(pool.Send("aaaaaa"))
	check(pool.Send("bbbbbb"))
	check(pool.Send("cccccc"))
	check(pool.Send("dddddd"))
	check(pool.Send("eeeeee"))
	check(pool.Send("ffffff"))
	fmt.Println("减少容量")
	pool.Expand(-6) // cap is 1
	// stop(1 * time.Second)
	check(pool.Send("gggggg"))
	check(pool.Send("hhhhhh"))
	check(pool.Send("iiiiii"))
	check(pool.Send("jjjjjj"))
	check(pool.Send("kkkkkk"))
	check(pool.Send("llllll"))
	check(pool.Send("mmmmmm"))

	pool.Wait()

}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func cycle(pool gopool.Pooler) {
	ticker := time.NewTicker(50 * time.Millisecond)
	for range ticker.C {
		fmt.Println("num-goroutine = ", pool.Cap())
	}
}

func stop(d time.Duration) {
	time.Sleep(d)
}

func run2(v interface{}) error {
	return nil
}

type testCase struct {
	name  string
	goNum int
}

func BenchmarkSend(b *testing.B) {
	var tcs []testCase
	for i := 1; i < 50; i++ {
		tcs = append(tcs, testCase{fmt.Sprintf("go%d", i*100), i * 100})
	}

	for _, tc := range tcs {
		b.Run(tc.name, func(b *testing.B) {
			pool := gopool.New(tc.goNum, run2)
			pool.StartGo()
			for i := 0; i < b.N; i++ {
				pool.Send(i)
			}
		})
	}
}

func BenchmarkSend1000(b *testing.B) {
	pool := gopool.New(1000, run2)
	pool.StartGo()
	for i := 0; i < b.N; i++ {
		pool.Send(i)
	}
}
func BenchmarkSend10000(b *testing.B) {
	pool := gopool.New(10000, run2)
	pool.StartGo()
	for i := 0; i < b.N; i++ {
		pool.Send(i)
	}
}

func BenchmarkSend100000(b *testing.B) {
	pool := gopool.New(100000, run2)
	pool.StartGo()
	for i := 0; i < b.N; i++ {
		pool.Send(i)
	}
}

func BenchmarkSend500000(b *testing.B) {
	pool := gopool.New(500000, run2)
	pool.StartGo()
	for i := 0; i < b.N; i++ {
		pool.Send(i)
	}
}

func BenchmarkSendWithout(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go run2(i)
	}
}

func BenchmarkSendParallel(b *testing.B) {
	var tcs []testCase
	for i := 1; i < 50; i++ {
		tcs = append(tcs, testCase{fmt.Sprintf("go%d", i*100), i * 100})
	}

	for _, tc := range tcs {
		b.Run(tc.name, func(b *testing.B) {
			pool := gopool.New(tc.goNum, run2)
			pool.StartGo()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					pool.Send(1)
				}
			})
		})
	}
}
