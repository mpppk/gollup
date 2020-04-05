package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mpppk/cli-template/registry"

	"github.com/mpppk/cli-template/domain/model"

	"github.com/labstack/echo"

	"github.com/mpppk/cli-template/handler"
)

func TestSumHistory(t *testing.T) {

	history1 := &model.SumHistory{
		IsNorm:  false,
		Date:    time.Time{},
		Numbers: model.Numbers{1, 2},
		Result:  3,
	}
	history2 := &model.SumHistory{
		IsNorm:  true,
		Date:    time.Time{},
		Numbers: model.Numbers{-1, 2},
		Result:  3,
	}
	history := []*model.SumHistory{history1, history2}
	h := registry.InitializeHandler(history)
	e := registry.InitializeServer(nil)

	type params struct {
		path string
	}
	type want struct {
		res  handler.SumHistoryResponse
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
				path: "/api/sum-history?limit=1",
			},
			want: want{
				res: handler.SumHistoryResponse{Result: []*model.SumHistory{
					history2,
				}},
				code: http.StatusOK,
			},
		},
		{
			params: params{
				path: "/api/sum-history?limit=2",
			},
			want: want{
				res: handler.SumHistoryResponse{Result: []*model.SumHistory{
					history1,
					history2,
				}},
				code: http.StatusOK,
			},
		},
		{
			params: params{
				path: "/api/sum-history?limit=xxx",
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

			err := h.SumHistory(e.NewContext(req, rec))

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
