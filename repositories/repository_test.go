package repositories

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBuildQueryClausesSingle(t *testing.T) {
	tables := []struct {
		mValues         map[string]interface{}
		mFilters        map[string]string
		expectedClauses string
		expectedValues  []interface{}
	}{
		{
			map[string]interface{}{
				"category": "foo",
			},
			map[string]string{
				"category": "account_category.name = ",
			},
			"WHERE account_category.name = $1;",
			[]interface{}{"foo"},
		},
		{
			map[string]interface{}{
				"amount": 12.5,
			},
			map[string]string{
				"amount":   "foo.amount = ",
				"category": "test_category.filter >= ",
			},
			"WHERE foo.amount = $1;",
			[]interface{}{12.5},
		},
		{
			map[string]interface{}{},
			map[string]string{},
			" ;",
			[]interface{}{},
		},
	}

	for _, table := range tables {
		actualClauses, actualValues, err := buildQueryClauses(table.mValues, table.mFilters)
		if err != nil {
			t.Errorf("Expected no error building query, got: %s", err)
		}

		if actualClauses != table.expectedClauses {
			t.Errorf("Query clauses built were incorrect, got: %s, want: %s",
				actualClauses, table.expectedClauses)
		}

		for i, v := range actualValues {
			if v != table.expectedValues[i] {
				t.Errorf("Query values built were incorrect, got: %v, want: %v",
					v, table.expectedValues[i])
			}
		}
	}
}

func TestBuildQueryClausesMultipe(t *testing.T) {
	now := time.Now()

	mValues := map[string]interface{}{
		"category": "foo",
		"amount":   45,
		"start":    now,
	}
	mFilters := map[string]string{
		"category": "account_category.name = ",
		"amount":   "test.amount >= ",
		"start":    "time.date_started <= ",
	}
	expectedValues := []interface{}{"foo", 45, now}

	actualClauses, actualValues, err := buildQueryClauses(mValues, mFilters)
	if err != nil {
		t.Errorf("Expected no error building query, got: %s", err)
	}

	if len(actualValues) != len(expectedValues) {
		t.Errorf("Number of values returned incorrect, got: %d, want: %d",
			len(actualValues), len(expectedValues))
	}

	checkParamsContains := func(param string) {
		if !strings.Contains(actualClauses, param) {
			t.Errorf("Query clauses have incorrect parameterization, got: %s, want to contain: %s",
				actualClauses, param)
		}
	}
	checkParamsContains("$1")
	checkParamsContains("$2")
	checkParamsContains("$3")

	for _, v := range mFilters {
		if !strings.Contains(actualClauses, v) {
			t.Errorf("Query clauses were built incorrectly, got: %s, want to contain: %s",
				actualClauses, v)
		}
	}
}

func TestBuildQueryClausesErrors(t *testing.T) {
	tables := []struct {
		mValues  map[string]interface{}
		mFilters map[string]string
	}{
		{
			map[string]interface{}{
				"category": "foo",
			},
			map[string]string{
				"otherCat": "account_category.name = ",
			},
		},
		{
			map[string]interface{}{
				"category": "foo",
				"otherCat": 15,
			},
			map[string]string{
				"otherCat": "account_category.name = ",
			},
		},
	}

	for _, table := range tables {
		actualClauses, actualValues, err := buildQueryClauses(table.mValues, table.mFilters)
		if err == nil {
			t.Error("Expected error building query, got none")
		}

		if err != ErrFiltersToMap {
			t.Errorf("Expected error with map filters, got: %s", err)
		}

		if actualClauses != "" {
			t.Errorf("Expected no clauses to be built, got: %s", actualClauses)
		}
		if len(actualValues) != 0 {
			vStr := ""
			for i, v := range actualValues {
				vStr += fmt.Sprintf("%v", v)
				if i != len(actualValues)-1 {
					vStr += ","
				}
			}
			t.Errorf("Expected no values to be built, got: %s", vStr)
		}
	}
}
