package cache

import (
	"fmt"
	"testing"
	"time"
)

func Test_Time(t *testing.T) {

}

func Test_Set(t *testing.T) {
	c := NewCache()

	ti := time.Now()
	k := fmt.Sprintf("%d", ti.Nanosecond())
	v := ti.Nanosecond()

	c.Set(k, v, NoExpiration)
}

func Test_Get(t *testing.T) {
	c := NewCache()

	ti := time.Now()
	k := fmt.Sprintf("k_%d", ti.Nanosecond())
	v := fmt.Sprintf("v_%d", ti.Nanosecond())

	c.Set(k, v, NoExpiration)

	vv, err := c.Get(k)

	if err != nil {
		t.Fatal(err)
	}

	if _v, ok := vv.(string); !ok {
		t.Fatal("Fail to get Value: %s", _v)
	}
}

func Test_Incr(t *testing.T) {
	c := NewCache()

	k := "q2"

	v, err := c.Incr(k)
	if err != nil {
		t.Fatal(err)
	}

	if v != 1 {
		t.Fatal("incr error: %d != 1", v)
	}
}

func Test_IncrBy(t *testing.T) {
	c := NewCache()
	k := "q3"

	v, err := c.IncrBy(k, 10)
	if err != nil {
		t.Fatal(err)
	}
	if v != 10 {
		t.Fatal("incrby err: %d = 10", v)
	}

	v, err = c.IncrBy(k, 20)
	if err != nil {
		t.Fatal(err)
	}

	if v != 30 {
		t.Fatal("incrby err: %d != 30", v)
	}
}
