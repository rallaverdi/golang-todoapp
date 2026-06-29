package core_http_server

import (
	"fmt"
	"net/http"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
}

func NewAPIVersionRouter(apiVersion ApiVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.Handler)
	}
}

/* версионирование схемы АПИ - добавляем префикс /api/v1/path_name
смысл в том, что мы не удаляем старую схему а добавляем новую, это дает возможность потребителям без проблем
перейти на новую схему без потери работоспособности приложения.

После того как все перейдут на новую версию - мы тогда уже спокойно можем удалить старую.
*/
