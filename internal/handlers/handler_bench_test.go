package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/configs"
	"github.com/KokoulinM/go-musthave-shortener-tpl/internal/workers"
	"github.com/golang/mock/gomock"
)

func BenchmarkHandlers_RetrieveShortURL(b *testing.B) {
	b.Run("RetrieveShortURL", func(b *testing.B) {
		ctx := context.Background()
		ctrl := gomock.NewController(b)
		defer ctrl.Finish()

		cfg := configs.New()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/Vq7zU8E5b7sLZo3qY82UKYRvQ-A=", nil)
		wp := workers.New(context.Background(), cfg.Workers, cfg.WorkersBuffer)

		go func() {
			wp.Run(ctx)
		}()

		defer wp.Stop()

		repoMock := NewMockRepository(ctrl)

		h := New(repoMock, cfg.BaseURL, wp)

		r := router(h)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			repoMock.EXPECT().GetURL(gomock.Any(), "Vq7zU8E5b7sLZo3qY82UKYRvQ-A=").Return("https://go.dev", nil).AnyTimes()

			r.ServeHTTP(w, req)
		}
	})
}
