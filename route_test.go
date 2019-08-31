package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFindBestRoute(t *testing.T) {
	testCases := []struct {
		Name     string
		File     string
		From     string
		To       string
		Limit    int
		Expected []RouteResp
	}{
		{
			Name: "Standard file", // The received input for this challenge
			File: `GRU,BRC,10
			BRC,SCL,5
			GRU,CDG,75
			GRU,SCL,20
			GRU,ORL,56
			ORL,CDG,5
			SCL,ORL,20`,
			From:  "GRU",
			To:    "CDG",
			Limit: 1,
			Expected: []RouteResp{
				RouteResp{
					Route: "GRU - BRC - SCL - ORL - CDG",
					Price: 40,
				},
			},
		},
		{
			Name: "No route found",
			File: `GRU,BRC,10
			BRC,SCL,5
			GRU,CDG,75`,
			From:     "GRU",
			To:       "ORL",
			Limit:    1,
			Expected: []RouteResp{},
		},
		{
			Name: "One route found and a infinte loop route ignored",
			File: `GRU,BRC,10
			BRC,GRU,5
			GRU,CDG,75
			BRC,ORL,10`,
			From:  "GRU",
			To:    "ORL",
			Limit: 1,
			Expected: []RouteResp{
				RouteResp{
					Route: "GRU - BRC - ORL",
					Price: 20,
				},
			},
		},
		{
			Name: "Two routes found with same value, should return the shortest route",
			File: `GRU,BRC,10
			BRC,SCL,5
			GRU,CDG,40
			GRU,SCL,20
			GRU,ORL,56
			ORL,CDG,5
			SCL,ORL,20`,
			From:  "GRU",
			To:    "CDG",
			Limit: 1,
			Expected: []RouteResp{
				RouteResp{
					Route: "GRU - CDG",
					Price: 40,
				},
			},
		},
		{
			Name: "File with many routes found and limit greater then 1, should return the best routes",
			File: `GRU,BRC,10
			BRC,SCL,5
			GRU,CDG,75
			GRU,SCL,20
			GRU,ORL,56
			ORL,CDG,5
			SCL,ORL,20`,
			From:  "GRU",
			To:    "CDG",
			Limit: 5,
			Expected: []RouteResp{
				RouteResp{
					Route: "GRU - BRC - SCL - ORL - CDG",
					Price: 40,
				},
				RouteResp{
					Route: "GRU - SCL - ORL - CDG",
					Price: 45,
				},
				RouteResp{
					Route: "GRU - ORL - CDG",
					Price: 61,
				},
				RouteResp{
					Route: "GRU - CDG",
					Price: 75,
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Logf("Running testCase #%d: %s...", i, tc.Name)
		if _, err := parseFile(bytes.NewBufferString(tc.File), true); err != nil {
			t.Fatal(err)
		}
		result := findBestRoute(tc.From, tc.To, tc.Limit)
		if !reflect.DeepEqual(tc.Expected, result) {
			t.Errorf("Unexpected result:\n'%#v'\nexpected:\n'%#v'", result, tc.Expected)
			continue
		}
		t.Logf(">>> OK!")
	}
}
