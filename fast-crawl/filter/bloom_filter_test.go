package filter

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestNewBloomFilter(t *testing.T) {
	seeds := []uint{7, 11, 13, 31, 37, 61}
	//filter := NewBloomFilter(2<<24, seeds, NewRedisSet(2<<24), defaultHash)
	filter := NewBloomFilter(2<<24, seeds, NewRedisSet(2<<24), DefaultHash)

	target := []string{"google.com", "youtube.com", "facebook.com", "baidu.com", "wikipedia.org", "qq.com",
		"taobao.com", "tmall.com", "yahoo.com", "amazon.com"}
	t0 := time.Now()
	for _, v := range target {
		filter.Add(v)
	}
	t1 := time.Now()
	fmt.Printf("add all: %f s\n", t1.Sub(t0).Seconds())

	var result bytes.Buffer
	for _, v := range target {
		if filter.Contains(v) {
			fmt.Fprintf(&result, "%s: %s\n", "true: ", v)
		} else {
			fmt.Fprintf(&result, "%s: %s\n", "false: ", v)
		}
	}
	t2 := time.Now()
	fmt.Printf("find all: %f s\n", t2.Sub(t1).Seconds())
	fmt.Println(result.String())
}
