package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestFilter(t *testing.T) {
	MongoQuery, Results := Cases(true)

	// fmt.Println(MongoQuery[3], Results[3])

	tests := []struct {
		input string
		want  string
	}{
		{MongoQuery[0], Results[0]},
		{MongoQuery[1], Results[1]},
		{MongoQuery[2], Results[2]},
		{MongoQuery[3], Results[3]},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("testno=%d", i), func(t *testing.T) {
			var m map[string]interface{}
			err := json.Unmarshal([]byte(tc.input), &m)
			if err != nil {
				log.Println(err)
			}
			got, err := MongoToSql(m)
			if err != nil {
				t.Fatalf("got err := %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %v; want %v", got, tc.want)
			} else {
				t.Logf("Success !")
			}

		})
	}
}
