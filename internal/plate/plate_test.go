package plate

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	RespData = `2.41 2.13 1.23 0.05 0.94 0.43 2.03 0.44 0.88 2.55 2.81 2.30
2.77 2.53 1.87 1.19 2.23 1.11 1.21 0.04 0.02 2.27 1.69 2.50
0.33 0.06 0.16 0.44 0.09 2.53 1.13 1.99 0.89 3.00 1.84 2.10
1.85 0.11 2.91 2.19 1.62 1.54 2.52 2.58 2.44 1.44 0.75 0.64
2.40 1.42 2.69 1.18 2.67 2.80 2.88 0.31 0.87 1.10 0.29 2.00
2.77 2.28 0.07 1.78 2.36 0.37 0.44 2.81 1.58 2.77 2.95 2.83
2.24 2.87 0.11 1.83 1.33 0.41 2.18 1.42 1.98 2.43 2.83 1.07
0.33 0.44 2.57 2.52 1.30 2.51 2.48 2.87 2.30 2.63 2.90 1.29`
)

func TestPlate_GetAvg(t *testing.T) {
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := RespData
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}))
	httpFailServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
	}))

	type fields struct {
		URL string
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				URL: successServer.URL,
			},
			want:    1.67,
			wantErr: false,
		},
		{
			name: "http_fail",
			fields: fields{
				URL: httpFailServer.URL,
			},
			want:    math.NaN(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Plate{
				URL: tt.fields.URL,
			}
			got, err := p.GetAvg()
			if (err != nil) != tt.wantErr {
				t.Errorf("Plate.GetAvg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("Plate.GetAvg() = %v, want %v", got, tt.want)
			}
		})
	}
}
