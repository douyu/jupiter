package tstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search"
	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	tableName = "ts_v5_test"
	indexName = "idx_ts_v5_test"
)

func TestCURD(t *testing.T) {
	t.Skip()

	type args struct {
		name string
		opts []interface{}
	}
	tests := []struct {
		name    string
		args    args
		config  string
		wantErr bool
	}{
		{
			name: "std new",
			args: args{
				name: "demo",
				opts: []interface{}{},
			},
			wantErr: false,
			config: `
			[jupiter.tablestore.demo]
			   debug = false
			   enableAccessLog = false
			   endPoint ="` + tablestoreEndponint + `"
			   instance = "` + tablestoreInstance + `"
			   accessKeyId ="` + accessKeyId + `"
			   accessKeySecret = "` + accessKeySecret + `"
			   requestTimeout = "30s"
			   slowThreshold = "10s"
			   maxIdleConnections = 2000
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, cfg.LoadFromReader(bytes.NewReader([]byte(tt.config)), toml.Unmarshal))
			// 初始化连接
			client := StdConfig(tt.args.name).MustSingleton()
			// 清理历史
			dropSearchIndex(client, true, t)
			dropTable(client, true, t)

			// 创建表
			createTableKeyAutoIncrementSample(client, t)
			defer dropTable(client, false, t)

			// 创建多元索引
			createSearchIndex(client, t)
			defer dropSearchIndex(client, false, t)

			// 写入数据
			primaryKeys := insertData(client, t)

			// 更新数据
			updateDate(client, primaryKeys, t)

			// 查询数据
			queryData(client, t)

			// 删除数据
			deleteData(client, primaryKeys, t)
		})
	}
}

func createTableKeyAutoIncrementSample(client *tablestore.TableStoreClient, t *testing.T) {
	createtableRequest := new(tablestore.CreateTableRequest)
	tableMeta := new(tablestore.TableMeta)
	tableMeta.TableName = tableName
	tableMeta.AddPrimaryKeyColumn("uid", tablestore.PrimaryKeyType_INTEGER)
	tableMeta.AddPrimaryKeyColumn("vid", tablestore.PrimaryKeyType_INTEGER)
	tableMeta.AddPrimaryKeyColumnOption("id", tablestore.PrimaryKeyType_INTEGER, tablestore.AUTO_INCREMENT) // 自增长列
	tableOption := new(tablestore.TableOption)
	tableOption.TimeToAlive = -1
	tableOption.MaxVersion = 1
	reservedThroughput := new(tablestore.ReservedThroughput)
	reservedThroughput.Readcap = 0
	reservedThroughput.Writecap = 0
	createtableRequest.TableMeta = tableMeta
	createtableRequest.TableOption = tableOption
	createtableRequest.ReservedThroughput = reservedThroughput

	_, err := client.CreateTable(createtableRequest)
	if err != nil {
		t.Fatal(errors.Wrap(err, "create table failed"))
	}
	t.Log("Create table finished")
}

func dropTable(client *tablestore.TableStoreClient, ignore bool, t *testing.T) {
	deleteReq := new(tablestore.DeleteTableRequest)
	deleteReq.TableName = tableName
	_, err := client.DeleteTable(deleteReq)
	if err != nil && ignore == false {
		t.Fatal(errors.Wrap(err, "Failed to delete table with error"))
	}
	t.Log("Delete table finished")
}

func createSearchIndex(client *tablestore.TableStoreClient, t *testing.T) {
	request := &tablestore.CreateSearchIndexRequest{}
	request.TableName = tableName
	request.IndexName = indexName

	schemas := []*tablestore.FieldSchema{}
	field1 := &tablestore.FieldSchema{
		FieldName:        proto.String("uid"),
		FieldType:        tablestore.FieldType_LONG,
		Index:            proto.Bool(true),
		EnableSortAndAgg: proto.Bool(true),
	}
	field2 := &tablestore.FieldSchema{
		FieldName:        proto.String("vid"),
		FieldType:        tablestore.FieldType_LONG,
		Index:            proto.Bool(true),
		EnableSortAndAgg: proto.Bool(true),
	}
	field3 := &tablestore.FieldSchema{
		FieldName:        proto.String("is_del"),
		FieldType:        tablestore.FieldType_LONG,
		Index:            proto.Bool(true),
		EnableSortAndAgg: proto.Bool(true),
	}
	field4 := &tablestore.FieldSchema{
		FieldName:        proto.String("source"),
		FieldType:        tablestore.FieldType_KEYWORD,
		Index:            proto.Bool(true),
		EnableSortAndAgg: proto.Bool(true),
	}
	schemas = append(schemas, field1, field2, field3, field4)

	indexSort := &search.Sort{
		Sorters: []search.Sorter{
			&search.FieldSort{
				FieldName: "uid",
				Order:     search.SortOrder_DESC.Enum(),
			},
		},
	}

	request.IndexSchema = &tablestore.IndexSchema{
		FieldSchemas: schemas,
		IndexSort:    indexSort,
	}
	_, err := client.CreateSearchIndex(request)
	if err != nil {
		t.Fatal(errors.Wrap(err, "create index failed"))
		return
	}
	t.Log("Create SearchIndex finished")
}

func dropSearchIndex(client *tablestore.TableStoreClient, ignore bool, t *testing.T) {
	request := &tablestore.DeleteSearchIndexRequest{}
	request.TableName = tableName
	request.IndexName = indexName
	_, err := client.DeleteSearchIndex(request)
	if err != nil && ignore == false {
		t.Fatal(errors.Wrap(err, "drop search index failed"))
		return
	}
	t.Log("Delete SearchIndex finished")
}

func insertData(client *tablestore.TableStoreClient, t *testing.T) []*tablestore.PrimaryKeyColumn {
	putRowRequest := new(tablestore.PutRowRequest)
	putRowChange := new(tablestore.PutRowChange)
	putRowChange.TableName = tableName
	putPk := new(tablestore.PrimaryKey)
	putPk.AddPrimaryKeyColumn("uid", int64(20000703))
	putPk.AddPrimaryKeyColumn("vid", int64(5058790))
	putPk.AddPrimaryKeyColumnWithAutoIncrement("id")
	putRowChange.PrimaryKey = putPk
	putRowChange.AddColumn("source", "web")
	putRowChange.AddColumn("is_del", int64(0))
	putRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
	putRowChange.ReturnType = tablestore.ReturnType_RT_PK
	putRowRequest.PutRowChange = putRowChange
	rest, err := client.PutRow(putRowRequest)
	if err != nil {
		t.Fatal(errors.Wrap(err, "insert date failed"))
	}
	t.Log("Insert Data finished")
	return rest.PrimaryKey.PrimaryKeys
}

func updateDate(client *tablestore.TableStoreClient, primaryKeys []*tablestore.PrimaryKeyColumn, t *testing.T) {
	updateRowRequest := new(tablestore.UpdateRowRequest)
	updateRowChange := new(tablestore.UpdateRowChange)
	updateRowChange.TableName = tableName
	updatePk := new(tablestore.PrimaryKey)
	updatePk.AddPrimaryKeyColumn("uid", primaryKeys[0].Value)
	updatePk.AddPrimaryKeyColumn("vid", primaryKeys[1].Value)
	updatePk.AddPrimaryKeyColumn("id", primaryKeys[2].Value)
	updateRowChange.PrimaryKey = updatePk
	updateRowChange.PutColumn("source", "h5")
	updateRowChange.PutColumn("is_del", int64(1))
	updateRowChange.SetCondition(tablestore.RowExistenceExpectation_EXPECT_EXIST)
	updateRowRequest.UpdateRowChange = updateRowChange
	_, err := client.UpdateRow(updateRowRequest)
	if err != nil {
		t.Fatal(errors.Wrap(err, "update data failed"))
	}
	t.Log("Update Data finished")
}

func queryData(client *tablestore.TableStoreClient, t *testing.T) {
	searchRequest := &tablestore.SearchRequest{}
	searchRequest.SetTableName(tableName)
	searchRequest.SetIndexName(indexName)
	query := &search.TermQuery{
		FieldName: "is_del",
		Term:      int64(1),
	}

	searchQuery := search.NewSearchQuery()
	searchQuery.SetQuery(query)
	searchQuery.SetGetTotalCount(true)
	searchRequest.SetSearchQuery(searchQuery)
	searchRequest.SetColumnsToGet(&tablestore.ColumnsToGet{ReturnAll: true})

	searchResponse, err := client.Search(searchRequest)
	if err != nil {
		t.Fatal("Failed to search with error: ", err)
		return
	}
	t.Log("RowCount：", searchResponse.TotalCount)
	for _, row := range searchResponse.Rows {
		st, _ := json.Marshal(row)
		fmt.Println("============》11", 11)
		t.Log("RowObj：", string(st))
	}
}

func deleteData(client *tablestore.TableStoreClient, primaryKeys []*tablestore.PrimaryKeyColumn, t *testing.T) {
	deleteRowReq := new(tablestore.DeleteRowRequest)
	deleteRowReq.DeleteRowChange = new(tablestore.DeleteRowChange)
	deleteRowReq.DeleteRowChange.TableName = tableName
	deletePk := new(tablestore.PrimaryKey)
	deletePk.AddPrimaryKeyColumn("uid", primaryKeys[0].Value)
	deletePk.AddPrimaryKeyColumn("vid", primaryKeys[1].Value)
	deletePk.AddPrimaryKeyColumn("id", primaryKeys[2].Value)
	deleteRowReq.DeleteRowChange.PrimaryKey = deletePk
	deleteRowReq.DeleteRowChange.SetCondition(tablestore.RowExistenceExpectation_EXPECT_EXIST)
	_, err := client.DeleteRow(deleteRowReq)
	if err != nil {
		t.Fatal(errors.Wrap(err, "delete data failed"))
	}
	t.Log("delete row finished")
}
