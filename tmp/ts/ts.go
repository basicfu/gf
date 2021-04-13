package ts

//
//import (
//	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
//	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore/search"
//	"github.com/basicfu/gutil/util"
//)
//
//var client = tablestore.NewClient("", "", "", "")
//
//type IndexQuery struct {
//	TableName string
//	IndexName string
//	Limit     int32
//	NextToken []byte
//	AllColumn bool
//	Columns   []string
//	Query     search.Query
//	//GetTotalCount bool
//	//Term      interface{}
//}
//type PkQuery struct {
//	TableName  string
//	Ids        []interface{}
//	Limit      int32
//	AllColumn  bool
//	Columns    []string
//	MaxVersion int
//	//Term      interface{}
//}
//type Column struct {
//	PrimaryKey bool
//	Name       string
//	Value      interface{}
//	Timestamp  int64
//}
//type Page struct {
//	NextToken []byte
//	Total     int64
//	Rows      []map[string]interface{}
//}
//
//func ListIndex(query IndexQuery) Page {
//	//querys := []search.Query{
//	//	&search.MatchAllQuery{},
//	//	&search.TermQuery{
//	//		FieldName: "Col_Keyword",
//	//		Term:      "tablestore",
//	//	},
//	//}
//	searchRequest := &tablestore.SearchRequest{}
//	searchRequest.SetTableName(query.TableName)
//	if query.IndexName == "" {
//		searchRequest.SetIndexName("default")
//	} else {
//		searchRequest.SetIndexName(query.IndexName)
//	}
//	searchQuery := search.NewSearchQuery()
//	//searchQuery.SetQuery(query)
//	if query.Limit == 0 {
//		searchQuery.SetLimit(20)
//	} else {
//		searchQuery.SetLimit(query.Limit)
//	}
//	if query.NextToken != nil {
//		searchQuery.SetToken(query.NextToken)
//	}
//	//if query.GetTotalCount ==false{
//	//
//	//}
//	searchQuery.SetQuery(query.Query)
//	searchQuery.SetGetTotalCount(true)
//	searchRequest.SetSearchQuery(searchQuery)
//	columnsToGet := &tablestore.ColumnsToGet{}
//	if query.AllColumn {
//		columnsToGet.ReturnAll = true
//	} else {
//		columnsToGet.Columns = query.Columns
//	}
//	searchRequest.SetColumnsToGet(columnsToGet)
//	resp, err := client.Search(searchRequest)
//	if err != nil {
//		panic(err.Error())
//	}
//	var rows []map[string]interface{}
//	for _, v := range resp.Rows {
//		row := make(map[string]interface{})
//		for _, pk := range v.PrimaryKey.PrimaryKeys {
//			row[pk.ColumnName] = pk.Value
//		}
//		for _, col := range v.Columns {
//			row[col.ColumnName] = col.Value
//		}
//		rows = append(rows, row)
//	}
//	return Page{
//		NextToken: resp.NextToken,
//		Total:     resp.TotalCount,
//		Rows:      rows,
//	}
//}
//func UpdateRow(tableName string, rows map[string]interface{}) {
//	updateRowRequest := new(tablestore.UpdateRowRequest)
//	updateRowChange := new(tablestore.UpdateRowChange)
//	updateRowChange.TableName = tableName
//	updatePk := new(tablestore.PrimaryKey)
//	updateRowChange.PrimaryKey = updatePk
//	for key, value := range rows {
//		if key == "id" {
//			updatePk.AddPrimaryKeyColumn(key, value)
//		} else {
//			updateRowChange.PutColumn(key, value)
//		}
//	}
//	updateRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
//	updateRowRequest.UpdateRowChange = updateRowChange
//	_, err := client.UpdateRow(updateRowRequest)
//	if err != nil {
//		panic(err.Error())
//	}
//}
//func ListPk(query PkQuery) []map[string]interface{} {
//	batchGetReq := &tablestore.BatchGetRowRequest{}
//	mqCriteria := &tablestore.MultiRowQueryCriteria{}
//	for _, id := range query.Ids {
//		pkToGet := new(tablestore.PrimaryKey)
//		pkToGet.AddPrimaryKeyColumn("id", id)
//		mqCriteria.AddRow(pkToGet)
//		if query.MaxVersion == 0 {
//			query.MaxVersion = 1
//		}
//		mqCriteria.MaxVersion = query.MaxVersion
//	}
//	mqCriteria.TableName = query.TableName
//	if !query.AllColumn {
//		mqCriteria.ColumnsToGet = query.Columns
//	}
//	batchGetReq.MultiRowQueryCriteria = append(batchGetReq.MultiRowQueryCriteria, mqCriteria)
//	resp, err := client.BatchGetRow(batchGetReq)
//	if err != nil {
//		panic(err.Error())
//	}
//	var rows []map[string]interface{}
//	for _, v := range resp.TableToRowsResult[query.TableName] {
//		row := make(map[string]interface{})
//		for _, pk := range v.PrimaryKey.PrimaryKeys {
//			row[pk.ColumnName] = pk.Value
//		}
//		for _, col := range v.Columns {
//			row[col.ColumnName] = col.Value
//		}
//		rows = append(rows, row)
//	}
//	return rows
//}
//
////cu消耗，范围版本
//func ListPkVersion(query PkQuery) map[interface{}]map[int64]map[interface{}]interface{} {
//	batchGetReq := &tablestore.BatchGetRowRequest{}
//	mqCriteria := &tablestore.MultiRowQueryCriteria{}
//	for _, id := range query.Ids {
//		pkToGet := new(tablestore.PrimaryKey)
//		pkToGet.AddPrimaryKeyColumn("id", id)
//		mqCriteria.AddRow(pkToGet)
//		mqCriteria.MaxVersion = query.MaxVersion
//	}
//	mqCriteria.TableName = query.TableName
//	if !query.AllColumn {
//		mqCriteria.ColumnsToGet = query.Columns
//	}
//	batchGetReq.MultiRowQueryCriteria = append(batchGetReq.MultiRowQueryCriteria, mqCriteria)
//	resp, err := client.BatchGetRow(batchGetReq)
//	if err != nil {
//		panic(err.Error())
//	}
//	rows := make(map[interface{}]map[int64]map[interface{}]interface{})
//	for _, v := range resp.TableToRowsResult[query.TableName] {
//		if !v.IsSucceed {
//			panic("查询失败")
//		}
//		if len(v.PrimaryKey.PrimaryKeys) == 0 && len(v.Columns) == 0 {
//			continue
//		}
//		key := v.PrimaryKey.PrimaryKeys[0].Value
//		timeMap := rows[key]
//		if timeMap == nil {
//			timeMap = map[int64]map[interface{}]interface{}{}
//		}
//		for _, col := range v.Columns {
//			columnMap := timeMap[col.Timestamp]
//			if columnMap == nil {
//				columnMap = map[interface{}]interface{}{}
//			}
//			columnMap[col.ColumnName] = col.Value
//			timeMap[col.Timestamp] = columnMap
//		}
//		rows[key] = timeMap
//	}
//	//list版
//	//var rows []Column
//	//for _, v := range resp.TableToRowsResult[query.TableName] {
//	//	for _, pk := range v.PrimaryKey.PrimaryKeys {
//	//		column := Column{}
//	//		column.PrimaryKey=true
//	//		column.Name=pk.ColumnName
//	//		column.Value=pk.Value
//	//		rows=append(rows, column)
//	//	}
//	//	for _, col := range v.Columns {
//	//		column := Column{}
//	//		column.Name=col.ColumnName
//	//		column.Value=col.Value
//	//		column.Timestamp=col.Timestamp
//	//		rows=append(rows, column)
//	//	}
//	//}
//	return rows
//}
//func BatchWriteRowChange(rows []tablestore.PutRowChange) {
//	if len(rows) == 0 {
//		return
//	}
//	if len(rows) > 200 {
//		panic("单次操作不能超过200条")
//	}
//	batchWriteReq := &tablestore.BatchWriteRowRequest{}
//	for _, v := range rows {
//		batchWriteReq.AddRowChange(&v)
//	}
//	_, err := client.BatchWriteRow(batchWriteReq)
//	if err != nil {
//		panic(err.Error())
//	}
//}
//func BatchWriteRowAndPk(tableName string, pk []interface{}, rows []map[string]interface{}) []interface{} {
//	return batchWriteRowBase(tableName, pk, rows, tablestore.RowExistenceExpectation_IGNORE)
//}
//func BatchWriteRow(tableName string, rows []map[string]interface{}) []interface{} {
//	return batchWriteRowBase(tableName, []interface{}{"id"}, rows, tablestore.RowExistenceExpectation_IGNORE)
//}
//func BatchWriteRowAndCondition(tableName string, rows []map[string]interface{}, condition tablestore.RowExistenceExpectation) []interface{} {
//	return batchWriteRowBase(tableName, []interface{}{"id"}, rows, condition)
//}
//func batchWriteRowBase(tableName string, pk []interface{}, rows []map[string]interface{}, condition tablestore.RowExistenceExpectation) []interface{} {
//	if len(rows) == 0 {
//		return nil
//	}
//	if len(rows) > 200 {
//		panic("单次操作不能超过200条")
//	}
//	var ids = make([]interface{}, 0)
//	batchWriteReq := &tablestore.BatchWriteRowRequest{}
//	for _, v := range rows {
//		putRowChange := new(tablestore.PutRowChange)
//		putRowChange.TableName = tableName
//		putPk := new(tablestore.PrimaryKey)
//		putRowChange.PrimaryKey = putPk
//		for _, pk := range pk {
//			putPk.AddPrimaryKeyColumn(pk.(string), v[pk.(string)])
//			if pk.(string) == "id" {
//				ids = append(ids, v[pk.(string)])
//			}
//		}
//		for key, value := range v {
//			if !util.Contains(pk, key) {
//				putRowChange.AddColumn(key, value)
//			}
//		}
//		putRowChange.SetCondition(condition)
//		batchWriteReq.AddRowChange(putRowChange)
//	}
//	response, err := client.BatchWriteRow(batchWriteReq)
//	if err != nil {
//		panic(err.Error())
//	} else {
//		//只做同一个表的操作
//		var successIds []interface{}
//		for _, v := range response.TableToRowsResult {
//			for _, v := range v {
//				if v.IsSucceed {
//					successIds = append(successIds, ids[v.Index])
//				} else {
//					if v.Error.Code != "OTSConditionCheckFail" {
//						panic("插入数据失败" + v.Error.Code + v.Error.Message)
//					}
//				}
//			}
//		}
//		return successIds
//	}
//}
