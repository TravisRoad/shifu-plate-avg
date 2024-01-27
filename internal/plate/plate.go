package plate

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Plate struct {
	URL string
}

func (p *Plate) GetAvg() (float64, error) {
	resp, err := http.Get(p.URL)
	if err != nil {
		return math.NaN(), err
	}
	if resp.StatusCode != http.StatusOK {
		return math.NaN(), fmt.Errorf("http request err")
	}
	defer resp.Body.Close()

	var sum float64
	var cnt int
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		for _, valStr := range values {
			// TODO: 这里可能需要精确值
			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return math.NaN(), err
			}
			sum += val
			cnt++
		}
	}
	if err := scanner.Err(); err != nil {
		return math.NaN(), err
	}

	return sum / float64(cnt), nil
}

func (p *Plate) Poll(ctx context.Context, d time.Duration) {
	for {
		select {
		case <-time.Tick(d):
			avg, err := p.GetAvg()
			if err != nil {
				slog.Error("", err)
				// jump out select
				break
			}
			slog.Info("", slog.Float64("avg", avg))
		case <-ctx.Done():
			return
		}
	}
}
