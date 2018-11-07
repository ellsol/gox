package sqlx

import (
	"testing"
	"fmt"
	"github.com/ellsol/gox/testx"
)

func TestStatementBuilder(t *testing.T) {

	stb, params := NewSelectStatement("*", "tablename").
		AddEqualCondition("label", 45).
		AddInCondition("inlabel", []string{"p1", "p2"}).
		AddEqualCondition("label2", 88).
		AddOffset(1200).
		AddLimit(120).GetStatementAndParams()

	expectedStatement := "SELECT * FROM tablename WHERE label = $1 AND inlabel IN ($2,$3) AND label2 = $4 OFFSET $5 LIMIT $6"
	if testx.CompareString("statement", expectedStatement, stb, t) {
		return
	}

	fmt.Println(params)
	if len(params) != 6 {
		t.Errorf("Param length is wrong [Expected 6, Actual: %v", len(params))
		return
	}

	if params[0] != 45 {
		t.Errorf("Param[0[ is wrong [Expected 45, Actual: %v", params[0])
		return
	}
	if params[1] != "p1" {
		t.Errorf("Param[1] is wrong [Expected 45, Actual: %v", params[1])
		return
	}

	if params[2] != "p2" {
		t.Errorf("Param[2] is wrong [Expected p2, Actual: %v", params[2])
		return
	}

	if params[3] != 88 {
		t.Errorf("Param[3] is wrong [Expected 45, Actual: %v", params[3])
		return
	}

	if params[4] != 1200 {
		t.Errorf("Param[4] is wrong [Expected 1200, Actual: %v", params[4])
		return
	}

	if params[5] != 120 {
		t.Errorf("Param[5] is wrong [Expected 120, Actual: %v", params[5])
		return
	}
}
