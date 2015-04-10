package cache

import (
	"testing"
)

func Test_Time(t *testing.T) {

}

func Test_NewCache(t *testing.T) {
    c := NewCache(DefaultExpiration)
    c.Set("k123", "v123", DefaultExpiration)  
}
