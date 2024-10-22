package fanin_test

import (
	"context"
	"testing"
	"time"

	"github.com/cyrusaf/fanin"
)

func TestFanIn(t *testing.T) {
	f, _ := fanin.WithContext[int](context.Background())
	for i := 0; i < 5; i++ {
		f.Go(func() (*int, error) {
			time.Sleep(time.Millisecond * 50)
			return &i, nil
		})
	}
	results, err := f.Wait()
	if err != nil {
		t.Fatalf("expected no error but got %v instead", err)
	}
	set := map[int]struct{}{}
	for _, result := range results {
		set[result] = struct{}{}
	}
	if len(set) != 5 {
		t.Fatalf("expected results to have 5 different values but got %v instead", results)
	}
	for i := 0; i < 5; i++ {
		if _, ok := set[i]; !ok {
			t.Fatalf("expected %d to be in results: %v", i, results)
		}
	}
}
