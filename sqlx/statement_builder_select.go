package sqlx

import (
	"fmt"
	"github.com/Leondroids/gox"
)

func NewSelectStatementBuilder(selectors string, tableName string) (*StatementBuilder) {
	return &StatementBuilder{
		Statement:         fmt.Sprintf("SELECT %v FROM %v", selectors, tableName),
		ConditionParams:   make([]interface{}, 0),
		ConditionPosition: 1,
	}
}

func (it *StatementBuilder) AddLikeCondition(conditionLabel string, conditionValue string) *StatementBuilder {
	newStatement := ""
	if !it.HasOneCondition {
		newStatement = fmt.Sprintf("%v WHERE %v ~ '^[%v]'", it.Statement, conditionLabel, conditionValue)
	} else {
		newStatement = fmt.Sprintf("%v AND %v ~ '^[%v]'", it.Statement, conditionLabel, conditionValue)
	}

	return &StatementBuilder{
		Statement:         newStatement,
		ConditionPosition: it.ConditionPosition,
		ConditionParams:   it.ConditionParams,
		HasOneCondition:   true,
	}
}

func (it *StatementBuilder) AddInCondition(conditionLabel string, values []string) *StatementBuilder {
	if len(values) == 0 {
		return it
	}

	conditionPosCounter := it.ConditionPosition - 1
	methodPlaceholder := gox.CommaSeparatedString(gox.MapStringListWithPos(values, func(key int, value string) string {
		conditionPosCounter++
		return fmt.Sprintf("$%v", conditionPosCounter)
	}))

	newStatement := ""
	if !it.HasOneCondition {
		newStatement = fmt.Sprintf("%v WHERE %v IN (%v)", it.Statement, conditionLabel, methodPlaceholder)
	} else {
		newStatement = fmt.Sprintf("%v AND %v IN (%v)", it.Statement, conditionLabel, methodPlaceholder)
	}

	params := it.ConditionParams
	for _, v := range values {
		params = append(params, v)
	}

	return &StatementBuilder{
		Statement:         newStatement,
		ConditionPosition: conditionPosCounter + 1,
		ConditionParams:   params,
		HasOneCondition:   true,
	}
}

func (it *StatementBuilder) MaybeAddEqualStringCondition(conditionLabel string, conditionValue string) *StatementBuilder {
	if conditionValue == "" {
		return it
	}

	return it.AddEqualCondition(conditionLabel, conditionValue)
}

func (it *StatementBuilder) AddEqualCondition(conditionLabel string, conditionValue interface{}) *StatementBuilder {
	newStatement := ""
	if !it.HasOneCondition {
		newStatement = fmt.Sprintf("%v WHERE %v = $%v", it.Statement, conditionLabel, it.ConditionPosition)
	} else {
		newStatement = fmt.Sprintf("%v AND %v = $%v", it.Statement, conditionLabel, it.ConditionPosition)
	}

	params := it.ConditionParams
	params = append(params, conditionValue)

	return &StatementBuilder{
		Statement:         newStatement,
		ConditionPosition: it.ConditionPosition + 1,
		ConditionParams:   params,
		HasOneCondition:   true,
	}
}

func (it *StatementBuilder) AddDateRange(conditionLabel string, dateFrom int64, dateTo int64) *StatementBuilder {
	if dateFrom == 0 && dateTo == 0 {
		return it
	}

	var df int64 = 0
	var dt int64 = 9223372036854775807

	if dateFrom > 0 {
		df = dateFrom
	}

	if dateTo > 0 {
		dt = dateTo
	}

	return it.AddRange(conditionLabel, df, dt)
}

func (it *StatementBuilder) AddInt64Range(conditionLabel string, from int64, to int64) *StatementBuilder {
	if from == 0 && to == 0 {
		return it
	}
	var bf int64 = 0
	var bt int64 = 9223372036854775807

	if from > 0 {
		bf = from
	}

	if to > 0 {
		bt = to
	}

	return it.AddRange(conditionLabel, bf, bt)
}

func (it *StatementBuilder) AddRange(conditionLabel string, valuesFrom interface{}, valuesTo interface{}) *StatementBuilder {

	newStatement := ""
	if !it.HasOneCondition {
		newStatement = fmt.Sprintf("%v WHERE %v BETWEEN $%v AND $%v", it.Statement, conditionLabel, it.ConditionPosition, it.ConditionPosition+1)
	} else {
		newStatement = fmt.Sprintf("%v AND %v BETWEEN $%v AND $%v", it.Statement, conditionLabel, it.ConditionPosition, it.ConditionPosition+1)
	}

	params := it.ConditionParams
	params = append(params, valuesFrom)
	params = append(params, valuesTo)

	return &StatementBuilder{
		Statement:         newStatement,
		ConditionPosition: it.ConditionPosition + 2,
		ConditionParams:   params,
		HasOneCondition:   true,
	}
}

func (it *StatementBuilder) OrderBy(orderBy *StatementOrderBy) *StatementBuilder {
	if orderBy == nil {
		return it
	}

	direction := "DESC"

	if !orderBy.DirectionDesc {
		direction = "ASC"
	}

	newStatement := fmt.Sprintf("%v ORDER BY %v %v", it.Statement, orderBy.By, direction)

	return &StatementBuilder{
		Statement:         newStatement,
		ConditionPosition: it.ConditionPosition,
		ConditionParams:   it.ConditionParams,
		HasOneCondition:   true,
	}
}

func (it *StatementBuilder) AddOffset(value int) *StatementBuilder {
	if value == 0 {
		return it
	}

	newStatement := fmt.Sprintf("%v OFFSET $%v", it.Statement, it.ConditionPosition)

	params := it.ConditionParams
	params = append(params, value)

	return &StatementBuilder{
		Statement:         newStatement,
		ConditionPosition: it.ConditionPosition + 1,
		ConditionParams:   params,
		HasOneCondition:   true,
	}
}

func (it *StatementBuilder) AddLimit(value int) *StatementBuilder {
	if value == 0 {
		return it
	}

	newStatement := fmt.Sprintf("%v LIMIT $%v", it.Statement, it.ConditionPosition)

	params := it.ConditionParams
	params = append(params, value)

	return &StatementBuilder{
		Statement:         newStatement,
		ConditionPosition: it.ConditionPosition + 1,
		ConditionParams:   params,
		HasOneCondition:   true,
	}
}

func (it *StatementBuilder) GetStatementAndParams() (string, []interface{}) {
	return it.Statement, it.ConditionParams
}

type StatementBuilder struct {
	Statement         string
	ConditionPosition int
	ConditionParams   []interface{}
	HasOneCondition   bool
}

type StatementOrderBy struct {
	By            string
	DirectionDesc bool
}
