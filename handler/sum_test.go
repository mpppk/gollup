package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mpppk/cli-template/infra"
	"github.com/mpppk/cli-template/registry"

	"github.com/labstack/echo"

	"github.com/mpppk/cli-template/handler"
)

func TestSum(t *testing.T) {
	h := registry.InitializeHandler(nil)
	e := infra.NewServer(h)

	type params struct {
		path string
	}
	type want struct {
		res  handler.SumResponse
		code int
	}
	tests := []struct {
		name    string
		params  params
		want    want
		wantErr bool
	}{
		{
			params: params{
				path: "/api/sum?a=-1&b=2",
			},
			want: want{
				res:  handler.SumResponse{Result: 1},
				code: http.StatusOK,
			},
		},
		{
			params: params{
				path: "/api/sum?a=-1&b=2&norm=true",
			},
			want: want{
				res:  handler.SumResponse{Result: 3},
				code: http.StatusOK,
			},
		},
		{
			params: params{
				path: "/api/sum?a=1&b=str",
			},
			want: want{
				code: http.StatusBadRequest,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.params.path, nil)
			rec := httptest.NewRecorder()

			err := h.Sum(e.NewContext(req, rec))

			if (err != nil) != tt.wantErr {
				t.Errorf("Handlers.Sum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var code int
			if tt.wantErr {
				httpError, ok := err.(*echo.HTTPError)
				if !ok {
					t.Fatalf("invalid err: %#v", err)
				}
				code = httpError.Code
			} else {
				code = rec.Code
			}

			if tt.want.code != code {
				t.Errorf("HTTP Status Code got = %d, want %d, body = %v", rec.Code, tt.want.code, rec.Body.String())
			}

			if tt.wantErr {
				return
			}

			gotRes := rec.Body.String()
			resJSON := toResponseJSON(t, tt.want.res)
			if resJSON != gotRes {
				t.Errorf("HTTP Response: got = %s, want %s", gotRes, resJSON)
			}
		})
	}
}

func toResponseJSON(t *testing.T, res interface{}) string {
	t.Helper()

	resContents, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("invalid test arg. res: %#v", res)
	}

	return fmt.Sprintln(string(resContents))
}
