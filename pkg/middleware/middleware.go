package middleware

import (
	"context"
	"fmt"
	"net/http"
	"watcharis/go-poc-protocal/pkg/dto"

	"github.com/labstack/echo/v4"
)

func AddProjectNameFromContext(ctx context.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			echoContext := c.Request().Context()
			projectName, ok := echoContext.Value(dto.APP_NAME).(string)
			fmt.Println("ok :", ok)
			if !ok {
				projectName := ctx.Value(dto.APP_NAME)
				if projectName != nil {
					nctx := context.WithValue(echoContext, dto.APP_NAME, projectName)
					c.SetRequest(c.Request().WithContext(nctx))
					if err := next(c); err != nil {
						fmt.Println("err :", err.Error())
						c.Error(err)
					}
				}
			} else {
				nctx := context.WithValue(echoContext, dto.APP_NAME, projectName)
				c.SetRequest(c.Request().WithContext(nctx))
				if err := next(c); err != nil {
					fmt.Println("err :", err.Error())
					c.Error(err)
				}
			}
			return nil
		}
	}
}

func AddProjectNameFromContextEchoWarpHttp(ctx context.Context) echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			echoContext := r.Context()
			projectName, ok := echoContext.Value(dto.APP_NAME).(string)
			if !ok {
				projectName := ctx.Value(dto.APP_NAME)
				if projectName != nil {
					nctx := context.WithValue(echoContext, dto.APP_NAME, projectName)
					r = r.WithContext(nctx)
					next.ServeHTTP(w, r)
				}
			} else {
				nctx := context.WithValue(echoContext, dto.APP_NAME, projectName)
				r = r.WithContext(nctx)
				next.ServeHTTP(w, r)
			}
		})
	})
}
