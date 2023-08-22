package shared

import (
	"github.com/gofiber/fiber/v2/log"
	"time"
)

func Top(label string) func() {
	x := time.Now()

	return func() {
		if time.Since(x) > 100*time.Millisecond {
			log.Warn(label, " time: ", time.Since(x))
		}
	}
}
