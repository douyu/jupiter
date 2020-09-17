package mongox

import (
	"net/http"

	"github.com/douyu/jupiter/pkg/application"
	"github.com/douyu/jupiter/pkg/xlog"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	// govern.RegisterStatSnapper("mongo", Stats)
	// govern.RegisterConfSnapper("mongo", Configs)
	http.HandleFunc("/debug/mongo/stats", func(w http.ResponseWriter, r *http.Request) {
		type mongoStatus struct {
			application.RuntimeStats
			Mongos map[string]interface{} `json:"mongos"`
		}
		var rets = mongoStatus{
			RuntimeStats: application.NewRuntimeStats(),
			Mongos:       make(map[string]interface{}, 0),
		}
		Range(func(name string, cc *mongo.Client) bool {
			rets.Mongos[name] = map[string]interface{}{
				"numberSessionsInProgress": cc.NumberSessionsInProgress(),
			}
			return true
		})

		_ = jsoniter.NewEncoder(w).Encode(rets)
	})
}

var _logger = xlog.JupiterLogger.With(xlog.FieldMod("mongodb"))
