package dynamodb

import (
	sdk "github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	defaultReadCapacity  = 1
	defaultWriteCapacity = 1
)

const (
	TableStatusActive = "ACTIVE"
)

// Table Represents the table object of
// Dynamo DB.
type Table struct {
	svc            *DynamoDB
	name           string
	nameWithPrefix string
	blueprint      *TableBlueprint

	putSpool   []*sdk.PutItemInput
	errorItems []*sdk.PutItemInput
}

func (t *Table) GetBlueprint() *TableBlueprint {
	return t.blueprint
}

func (t *Table) SetBlueprint(bp *TableBlueprint) {
}

func (t *Table) RefreshBlueprint() (*TableBlueprint, error) {
	return nil, nil
}

func (t *Table) UpdateThroughput(r, w int64) error {
	return nil
}

func (t *Table) UpdateWriteThroughput(w int64) error {
	return nil
}

func (t *Table) UpdateReadThroughput(r int64) error {
	return nil
}

//// AddItem adds the item that to be put to the table.
//func (t *Table) AddItem(item *PutItem) {
//
//}
//
//func (t *Table) Put() error {
//	return nil
//}
//
//func (t *Table) BatchPut() error {
//	return nil
//}
//
//func (t *Table) Scan() (*QueryResult, error) {
//	return nil, nil
//}
//
//func (t *Table) ScanWithCondition(cond *ConditionList) (*QueryResult, error) {
//	return nil, nil
//}
//
//func (t *Table) Query(cond *ConditionList) (*QueryResult, error) {
//	return nil, nil
//}
//
//func (t *Table) Count(cond *ConditionList) (*QueryResult, error) {
//	return nil, nil
//}
//
//func (t *Table) NewConditionList() *ConditionList {
//	return nil
//}
//
//func (t *Table) GetOne(hashValue interface{}, rangeValue ...interface{}) (map[string]interface{}, error) {
//	return nil, nil
//}
//
//func (t *Table) Delete(hashValue interface{}, rangeValue ...interface{}) error {
//	return nil
//}
//
//func (t *Table) ForceDeleteAll() error {
//	return nil
//}
