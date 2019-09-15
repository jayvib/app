package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	sdk "github.com/aws/aws-sdk-go/service/dynamodb"
)

type TableBlueprint struct {
	// name is the table name.
	name          string
	readCapacity  int64
	writeCapacity int64

	// HashKey a key schema element declares
	// the hash key of the table without
	// declaring yet the type.
	HashKey *sdk.KeySchemaElement

	// RangKey is a key schema element that
	// declares the hash key of the table
	// without declaring yet the type.
	RangeKey               *sdk.KeySchemaElement
	LocalSecondaryIndexes  []*sdk.LocalSecondaryIndex
	GlobalSecondaryIndexes []*sdk.GlobalSecondaryIndex

	// AttributeDefinitions is a list of attribute definition
	// wherein declaring type type of the attribute.
	// possible types are:
	// S - string
	// N - number
	// B - binary
	AttributeDefinitions map[string]*sdk.AttributeDefinition

	itemCount              int64
	status                 string
	numberOfDecreasesToday int64
}

type TableBlueprintOptFunc func(blueprint *TableBlueprint)

func NewTableBlueprint(tableName, hashKeyName string, opts ...TableBlueprintOptFunc) *TableBlueprint {
	tb := &TableBlueprint{
		name:                   tableName,
		HashKey:                NewHashKey(hashKeyName), // This is required.
		AttributeDefinitions:   make(map[string]*sdk.AttributeDefinition),
		LocalSecondaryIndexes:  make([]*sdk.LocalSecondaryIndex, 0),
		GlobalSecondaryIndexes: make([]*sdk.GlobalSecondaryIndex, 0),
		readCapacity:           defaultReadCapacity,
		writeCapacity:          defaultWriteCapacity,
	}

	if opts != nil {
		for _, opt := range opts {
			opt(tb)
		}
	}
	return tb
}

func newTableBlueprintFromDescription(desc *sdk.TableDescription) *TableBlueprint {
	if desc == nil {
		return nil
	}

	bp := &TableBlueprint{
		name:                   *desc.TableName,
		status:                 *desc.TableStatus,
		itemCount:              *desc.ItemCount,
		readCapacity:           *desc.ProvisionedThroughput.ReadCapacityUnits,
		writeCapacity:          *desc.ProvisionedThroughput.WriteCapacityUnits,
		numberOfDecreasesToday: *desc.ProvisionedThroughput.NumberOfDecreasesToday,
		AttributeDefinitions:   make(map[string]*sdk.AttributeDefinition),
		LocalSecondaryIndexes:  make([]*sdk.LocalSecondaryIndex, 0),
		GlobalSecondaryIndexes: make([]*sdk.GlobalSecondaryIndex, 0),
	}

	for _, attr := range desc.AttributeDefinitions {
		bp.AttributeDefinitions[aws.StringValue(attr.AttributeName)] = attr
	}

	for _, schema := range desc.KeySchema {
		switch aws.StringValue(schema.KeyType) {
		case "HASH":
			bp.HashKey = schema
		case "RANGE":
			bp.RangeKey = schema
		}
	}

	for _, lsi := range desc.LocalSecondaryIndexes {
		bp.LocalSecondaryIndexes = append(bp.LocalSecondaryIndexes, &sdk.LocalSecondaryIndex{
			IndexName:  lsi.IndexName,
			KeySchema:  lsi.KeySchema,
			Projection: lsi.Projection,
		})
	}

	for _, gsi := range desc.GlobalSecondaryIndexes {
		bp.GlobalSecondaryIndexes = append(bp.GlobalSecondaryIndexes, &sdk.GlobalSecondaryIndex{
			IndexName:  gsi.IndexName,
			KeySchema:  gsi.KeySchema,
			Projection: gsi.Projection,
			ProvisionedThroughput: &sdk.ProvisionedThroughput{
				WriteCapacityUnits: gsi.ProvisionedThroughput.WriteCapacityUnits,
				ReadCapacityUnits:  gsi.ProvisionedThroughput.ReadCapacityUnits,
			},
		})
	}
	return bp
}

func (tb *TableBlueprint) CreateTableInput() *sdk.CreateTableInput {
	input := &sdk.CreateTableInput{
		TableName:             aws.String(tb.GetTableName()),
		KeySchema:             tb.GetKeySchemaList(),
		AttributeDefinitions:  tb.GetAttributeDefinitionList(),
		ProvisionedThroughput: newProvisionThroughput(tb.readCapacity, tb.writeCapacity),
	}

	if tb.HasLocalSecondaryIndex() {
		input.LocalSecondaryIndexes = tb.LocalSecondaryIndexes
	}

	if tb.HasGlobalSecondaryIndex() {
		input.GlobalSecondaryIndexes = tb.GlobalSecondaryIndexes
	}
	return input
}

func (tb *TableBlueprint) HasRangeKey() bool {
	return tb.RangeKey != nil
}

func (tb *TableBlueprint) HasLocalSecondaryIndex() bool {
	return len(tb.LocalSecondaryIndexes) != 0
}

func (tb *TableBlueprint) HasGlobalSecondaryIndex() bool {
	return len(tb.GlobalSecondaryIndexes) != 0
}

func (tb *TableBlueprint) ListLocalSecondaryIndexes() []*sdk.LocalSecondaryIndex {
	return tb.LocalSecondaryIndexes
}

func (tb *TableBlueprint) ListGlobalSecondaryIndexes() []*sdk.GlobalSecondaryIndex {
	return tb.GlobalSecondaryIndexes
}

func (tb *TableBlueprint) GetReadCapacity() int64 {
	return tb.readCapacity
}

func (tb *TableBlueprint) GetWriteCapacity() int64 {
	return tb.writeCapacity
}

func (tb *TableBlueprint) GetTableName() string {
	return tb.name
}

func (tb *TableBlueprint) GetKeys() map[string]*sdk.KeySchemaElement {
	keys := map[string]*sdk.KeySchemaElement{
		tb.GetHashKeyName(): tb.HashKey,
	}

	if tb.HasRangeKey() {
		keys[tb.GetRangeKeyName()] = tb.RangeKey
	}
	return keys
}

func (tb *TableBlueprint) GetKeySchemaList() []*sdk.KeySchemaElement {
	keySchema := []*sdk.KeySchemaElement{
		tb.HashKey,
	}
	if tb.HasRangeKey() {
		keySchema = append(keySchema, tb.RangeKey)
	}
	return keySchema
}

func (tb *TableBlueprint) GetHashKeyName() string {
	return aws.StringValue(tb.HashKey.AttributeName)
}

func (tb *TableBlueprint) GetRangeKeyName() string {
	if tb.HasRangeKey() {
		return aws.StringValue(tb.RangeKey.AttributeName)
	}
	return ""
}

func (tb *TableBlueprint) GetAttributeDefinitionList() []*sdk.AttributeDefinition {
	ad := make([]*sdk.AttributeDefinition, 0)
	for _, d := range tb.AttributeDefinitions {
		ad = append(ad, d)
	}
	return ad
}

func (tb *TableBlueprint) IsActive() bool {
	if tb.status == TableStatusActive {
		return true
	}
	return false
}

// TableBlueprint Options
func HashKeyNumberOpt(hashKeyName string) TableBlueprintOptFunc {
	return func(bp *TableBlueprint) {
		if bp.AttributeDefinitions == nil {
			bp.AttributeDefinitions = map[string]*sdk.AttributeDefinition{
				hashKeyName: NewNumberAttributeDefinition(hashKeyName),
			}
		} else {
			bp.AttributeDefinitions[hashKeyName] = NewNumberAttributeDefinition(hashKeyName)
		}
	}
}

func HashKeyStringOpt(hashKeyName string) TableBlueprintOptFunc {
	return func(bp *TableBlueprint) {
		if bp.AttributeDefinitions == nil {
			bp.AttributeDefinitions = map[string]*sdk.AttributeDefinition{
				hashKeyName: NewStringAttributeDefinition(hashKeyName),
			}
		} else {
			bp.AttributeDefinitions[hashKeyName] = NewStringAttributeDefinition(hashKeyName)
		}
	}
}

func RangeKeyStringOpt(keyName string) TableBlueprintOptFunc {
	return func(bp *TableBlueprint) {
		bp.RangeKey = NewRangeKey(keyName)

		// Assuming here that the hash key is already set.
		// And the attribute map has at least one attribute
		// definition.
		bp.AttributeDefinitions[keyName] = NewStringAttributeDefinition(keyName)
	}
}

func RangeKeyNumberOpt(keyName string) TableBlueprintOptFunc {
	return func(bp *TableBlueprint) {
		bp.RangeKey = NewRangeKey(keyName)

		// Assuming here that the hash key is already set.
		// And the attribute map has at least one attribute
		// definition.
		bp.AttributeDefinitions[keyName] = NewNumberAttributeDefinition(keyName)
	}
}

func LocalSecondaryIndexStringOpt(indexName, keyName string) TableBlueprintOptFunc {
	return func(bp *TableBlueprint) {
		// Adding keyName to the list of attribute
		// because probably this is a new attribute.

		// In local secondary index the keyName is just an ordinary
		// attribute in the base table.
		bp.AttributeDefinitions[keyName] = NewStringAttributeDefinition(keyName)
		rangeKey := NewRangeKey(keyName)
		schema := NewKeySchemaElements(bp.HashKey, rangeKey)
		lsi := NewLocalSecondaryIndex(indexName, schema)
		bp.LocalSecondaryIndexes = append(bp.LocalSecondaryIndexes, lsi)
	}
}

func LocalSecondaryIndexNumberOpt(indexName, keyName string) TableBlueprintOptFunc {
	return func(bp *TableBlueprint) {
		bp.AttributeDefinitions[keyName] = NewNumberAttributeDefinition(keyName)
		rangeKey := NewRangeKey(keyName)
		schema := NewKeySchemaElements(bp.HashKey, rangeKey)
		lsi := NewLocalSecondaryIndex(indexName, schema)
		bp.LocalSecondaryIndexes = append(bp.LocalSecondaryIndexes, lsi)
	}
}

func ThroughputOpt(r, w int64) TableBlueprintOptFunc {
	return func(bp *TableBlueprint) {
		bp.readCapacity = r
		bp.writeCapacity = w
	}
}

func GlobalSecondaryIndexStringOpt(indexName string, keyName ...string) TableBlueprintOptFunc {
	return func(bp *TableBlueprint) {
		// TODO: Implement me
	}
}
