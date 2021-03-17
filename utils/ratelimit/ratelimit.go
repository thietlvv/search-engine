package ratelimit

import (
	"auth/models"
	"fmt"
	"strconv"
	"time"

	"auth/utils/array"
	dt "auth/utils/datastore"
	"auth/utils/jwt"

	"cloud.google.com/go/datastore"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
)

func Skipper(c echo.Context) bool {
	bypass_endpoint := []string{
		"/",
	}
	if array.Contains(bypass_endpoint, c.Path()) {
		return true
	}
	return false
}

func Process(client *datastore.Client, redis *redis.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if Skipper(c) {
				return next(c)
			}

			filter := []dt.DatastoreFilter{}
			docs := new([]models.LimitConfig)
			err := dt.Get(client, c.Request().Context(), "limit_configs", docs, filter, 100)
			if err != nil {
				fmt.Println(err)
				return echo.ErrBadRequest
			}

			for _, doc := range *docs {
				fmt.Println(doc.Prefix, doc.Max, doc.Duration)
				key_lock := ""
				if doc.Prefix == "IP" {
					fmt.Println("Header: ", c.Request().Header)
					ip := c.Request().Header.Get("X-Custom-User-Ip")
					if ip != "" {
						key_lock = fmt.Sprintf(`LIMIT_%s_%s`, doc.Prefix, ip)
					}
				} else if doc.Prefix == "TOKEN" {
					token, err := jwt.GetJWTFromHeader(c, "Bearer")
					if err != nil {
						continue
					}
					key_lock = fmt.Sprintf(`LIMIT_%s_%s`, doc.Prefix, token)
				}
				if key_lock != "" {
					val, err := redis.Get(key_lock).Result()
					i, _ := strconv.Atoi(val)
					fmt.Println(key_lock, i)
					if err == nil {
						if i < doc.Max {
							_, err = redis.Incr(key_lock).Result()
						} else {
							return echo.ErrTooManyRequests
						}
					} else {
						_, err = redis.Set(key_lock, 1, time.Minute*time.Duration(doc.Duration)).Result()
					}
				}
			}
			return next(c)
		}
	}
}
