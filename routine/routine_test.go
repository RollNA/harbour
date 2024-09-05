package routine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSafe(t *testing.T) {
	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan struct{})
	go RunSafe(func() {
		defer func() {
			ch <- struct{}{}
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestRunSafeCtx(t *testing.T) {
	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan struct{})
	go RunSafeCtx(context.Background(), func() {
		defer func() {
			ch <- struct{}{}
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestGoSafeCtx(t *testing.T) {
	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan struct{})
	GoSafeCtx(context.Background(), func() {
		defer func() {
			ch <- struct{}{}
		}()

		panic("panic")
	})

	<-ch
	i++
}
