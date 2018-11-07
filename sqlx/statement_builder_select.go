package sqlx

import (
	"fmt"
	"github.com/Leondroids/gox"
)

func NewSelectStatement(selectors string, tableName string) (*StatementBuilder) {
	return &StatementBuilder{
		statement:         fmt.Sprintf("SELECT %v FROM %v", selectors, tableName),
		conditionParams:   make([]interface{}, 0),
		conditionPosition: 1,
	}
}

func (it *StatementBuilder) AddLikeCondition(conditionLabel string, conditionValue string) *StatementBuilder {
	newStatement := ""
	if !it.hasOneCondition {
		newStatement = fmt.Sprintf("%v WHERE %v ~ '^[%v]'", it.statement, conditionLabel, conditionValue)
	} else {
		newStatement = fmt.Sprintf("%v AND %v ~ '^[%v]'", it.statement, conditionLabel, conditionValue)
	}

	return &StatementBuilder{
		statement:         newStatement,
		conditionPosition: it.conditionPosition,
		conditionParams:   it.conditionParams,
		hasOneCondition:   true,
	}
}

func (it *StatementBuilder) AddInCondition(conditionLabel string, values []string) *StatementBuilder {
	if len(values) == 0 {
		return it
	}

	conditionPosCounter := it.conditionPosition - 1
	methodPlaceholder := gox.CommaSeparatedString(gox.MapStringListWithPos(values, func(key int, value string) string {
		conditionPosCounter++
		return fmt.Sprintf("$%v", conditionPosCounter)
	}))

	newStatement := ""
	if !it.hasOneCondition {
		newStatement = fmt.Sprintf("%v WHERE %v IN (%v)", it.statement, conditionLabel, methodPlaceholder)
	} else {
		newStatement = fmt.Sprintf("%v AND %v IN (%v)", it.statement, conditionLabel, methodPlaceholder)
	}

	params := it.conditionParams
	for _, v := range values {
		params = append(params, v)
	}

	return &StatementBuilder{
		statement:         newStatement,
		conditionPosition: conditionPosCounter + 1,
		conditionParams:   params,
		hasOneCondition:   true,
	}
}

func (it *StatementBuilder) MaybeAddEqualStringCondition(conditionLabel string, conditionValue string) *StatementBuilder {
	if conditionValue == "" {
		return it
	}

	return it.AddEqualCondition(conditionLabel, conditionValue)
}



func (it *StatementBuilder) FilterBoolean(conditionLabel string, conditionValue bool) *StatementBuilder {
	return it.AddEqualCondition(conditionLabel, conditionValue)
}

func (it *StatementBuilder) AddEqualCondition(conditionLabel string, conditionValue interface{}) *StatementBuilder {
	newStatement := ""
	if !it.hasOneCondition {
		newStatement = fmt.Sprintf("%v WHERE %v = $%v", it.statement, conditionLabel, it.conditionPosition)
	} else {
		newStatement = fmt.Sprintf("%v AND %v = $%v", it.statement, conditionLabel, it.conditionPosition)
	}

	params := it.conditionParams
	params = append(params, conditionValue)

	return &StatementBuilder{
		statement:         newStatement,
		conditionPosition: it.conditionPosition + 1,
		conditionParams:   params,
		hasOneCondition:   true,
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
	if !it.hasOneCondition {
		newStatement = fmt.Sprintf("%v WHERE %v BETWEEN $%v AND $%v", it.statement, conditionLabel, it.conditionPosition, it.conditionPosition+1)
	} else {
		newStatement = fmt.Sprintf("%v AND %v BETWEEN $%v AND $%v", it.statement, conditionLabel, it.conditionPosition, it.conditionPosition+1)
	}

	params := it.conditionParams
	params = append(params, valuesFrom)
	params = append(params, valuesTo)

	return &StatementBuilder{
		statement:         newStatement,
		conditionPosition: it.conditionPosition + 2,
		conditionParams:   params,
		hasOneCondition:   true,
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

	newStatement := fmt.Sprintf("%v ORDER BY %v %v", it.statement, orderBy.By, direction)

	return &StatementBuilder{
		statement:         newStatement,
		conditionPosition: it.conditionPosition,
		conditionParams:   it.conditionParams,
		hasOneCondition:   true,
	}
}

func (it *StatementBuilder) AddOffset(value int) *StatementBuilder {
	if value == 0 {
		return it
	}

	newStatement := fmt.Sprintf("%v OFFSET $%v", it.statement, it.conditionPosition)

	params := it.conditionParams
	params = append(params, value)

	return &StatementBuilder{
		statement:         newStatement,
		conditionPosition: it.conditionPosition + 1,
		conditionParams:   params,
		hasOneCondition:   true,
	}
}

func (it *StatementBuilder) AddLimit(value int) *StatementBuilder {
	if value == 0 {
		return it
	}

	newStatement := fmt.Sprintf("%v LIMIT $%v", it.statement, it.conditionPosition)

	params := it.conditionParams
	params = append(params, value)

	return &StatementBuilder{
		statement:         newStatement,
		conditionPosition: it.conditionPosition + 1,
		conditionParams:   params,
		hasOneCondition:   true,
	}
}

func (it *StatementBuilder) GetStatementAndParams() (string, []interface{}) {
	return it.statement, it.conditionParams
}

type StatementBuilder struct {
	statement         string
	conditionPosition int
	conditionParams   []interface{}
	hasOneCondition   bool
}

type StatementOrderBy struct {
	By            string
	DirectionDesc bool
}
