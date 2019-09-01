package route

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	input := `GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
GRU,SCL,20
GRU,ORL,56
ORL,CDG,5
SCL,ORL,20`
	expectedParsed := map[string][][]string{
		"GRU": [][]string{
			[]string{"GRU", "BRC", "10"},
			[]string{"GRU", "CDG", "75"},
			[]string{"GRU", "SCL", "20"},
			[]string{"GRU", "ORL", "56"},
		},
		"BRC": [][]string{
			[]string{"BRC", "SCL", "5"},
		},
		"ORL": [][]string{
			[]string{"ORL", "CDG", "5"},
		},
		"SCL": [][]string{
			[]string{"SCL", "ORL", "20"},
		},
	}

	resultParsed, err := ParseFile(bytes.NewBufferString(input), false)
	if err != nil {
		t.Fatal(err)
	}
	for from, expectedRoutes := range expectedParsed {
		resultRoutes, ok := resultParsed[from]
		if !ok {
			t.Errorf("FromRoute '%s' not found at result", from)
			continue
		}
		if len(resultRoutes) != len(expectedRoutes) {
			t.Errorf("Unexpected routes size %d of FromRoute '%s', expected %d", len(resultRoutes), from, len(expectedRoutes))
			continue
		}
		for i := range expectedRoutes {
			ex := fmt.Sprintf("%v", expectedRoutes[i])
			rs := fmt.Sprintf("%v", resultRoutes[i])
			if ex != rs {
				t.Errorf("Unexpected route '%s' at index #%d, expected '%s'", rs, i, ex)
			}
		}
	}
}

func TestParseFileValidation(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    string
		ErrorMsg string
	}{
		{
			Name:     "No error detected",
			Input:    "GRU,BRC,10",
			ErrorMsg: "",
		},
		{
			Name:     "Error on line size at line 1",
			Input:    "GRU,BRC",
			ErrorMsg: "line 1: Invalid line format",
		},
		{
			Name:     "Error on src route at line 1 (first column of line)",
			Input:    " ,BRC,10",
			ErrorMsg: "line 1: Invalid src route",
		},
		{
			Name:     "Error on dst route at line 2 (second column of line)",
			Input:    "GRU,BRC,10\nGRU, ,5",
			ErrorMsg: "line 2: Invalid dst route",
		},
		{
			Name:     "Error route price at line 1 (third column of line)",
			Input:    "GRU,BRC,1ABC",
			ErrorMsg: "line 1: Invalid price",
		},
	}

	for i, tc := range testCases {
		t.Logf("Running testCase '%s'...", tc.Name)
		_, err := ParseFile(bytes.NewBufferString(tc.Input), true)
		if err != nil && tc.ErrorMsg == "" {
			t.Errorf("Unexpected non <nil> error at testCase #%d: %v", i, err)
			continue
		}
		if err == nil && tc.ErrorMsg != "" {
			t.Errorf("Unexpected <nil> error at testCase #%d, expected something like '%s'", i, tc.ErrorMsg)
			continue
		}
		if err != nil && !strings.Contains(err.Error(), tc.ErrorMsg) {
			t.Errorf("Unexpected error '%s' at testCase #%d, expected something like '%s'", err.Error(), i, tc.ErrorMsg)
			continue
		}
		t.Logf(">>> OK!")
	}
}

func TestAddRoute(t *testing.T) {
	File = "sample.csv"

	testCases := []struct {
		Name            string
		Route           Route
		ExpectedCreated bool
		ExpectedContent string
	}{
		{
			Name: "Add new route",
			Route: Route{
				From:  "ABC",
				To:    "DEF",
				Price: 50,
			},
			ExpectedCreated: true,
			ExpectedContent: `GRU,BRC,10
BRC,SCL,5
ABC,DEF,50
`,
		},
		{
			Name: "Update existent route",
			Route: Route{
				From:  "GRU",
				To:    "BRC",
				Price: 35,
			},
			ExpectedCreated: false,
			ExpectedContent: `GRU,BRC,35
BRC,SCL,5
`,
		},
	}

	for i, tc := range testCases {
		t.Logf("Running testcase #%d: %s...", i, tc.Name)
		route := testCases[i].Route
		newContent, created, err := addRoute(&route, false)
		if err != nil {
			t.Fatal(err)
		}
		if tc.ExpectedCreated != created {
			t.Errorf("Unexpected created '%v' at testcase #%d, expected '%v'", created, i, tc.ExpectedCreated)
			continue
		}
		if string(newContent) != tc.ExpectedContent {
			t.Errorf("Unexpected file at testcase #%d:\n%s\nexpected:\n%s", i, newContent, tc.ExpectedContent)
			continue
		}
		t.Logf(">>> OK!")
	}
}
