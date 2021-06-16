package swagger

import "github.com/labstack/echo/v4"

/*
ws.Route(ws.GET("/").
		To(h.getLog).
		Metadata(restfulspec.KeyOpenAPITags, []string{}).
		Doc("查询资源日志, e.g. /api/v1/log?l=jobId=111&l=logId=111&clusterID=1").
		Param(ws.QueryParameter("clusterID", "集群ID").DataType("int").Required(true)).
		Param(ws.QueryParameter("l", "资源标签, 支持多个标签, e.g. jobId=111").Required(false)).
		Param(ws.QueryParameter("sinceTimestamp", "日志开始时间").Required(false)).
		Param(ws.QueryParameter("pod", "Pod name").Required(false)).
		Param(ws.QueryParameter("container", "Container name").Required(false)).
		Returns(http.StatusBadRequest, "Fail", types.Error{}).
		Returns(http.StatusOK, "Ok", "resource name"))

*/

type param struct {
	kind     string // query|path
	dataType string // int, string, bool
	required bool
}

type reply struct {
}

type route struct {
	document    string
	pathParams  map[string]string
	queryParams map[string]string
}

func (r *route) Doc(doc string) *route {
	r.document = doc
	return r
}

func (r *route) Param(param *param) *route {
	r.pathParams[path] = doc
	return r
}

func (r *route) QueryParam(query string, doc string) *route {
	r.queryParams[query] = doc
	return r
}

func Route(route *echo.Route) *route {

}
