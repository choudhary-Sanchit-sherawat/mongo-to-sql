package main

// for the test CASES
func Cases(i bool) ([]string, []string) {

	s0 := `{
		"$and": [
		{ "x": { "$gt": 0 } },
			   { "x": { "$lt": 0 } },
		 	  {"x": { "$eq":3}},
			   {"x":{"$in":["1","2","3"]}},
		 	   {"$or": [
			   { "x": { "$ne": 0 } },
			   { "x": { "$ne": 0 } }
		 	   ]}
			]
	
		  }  `
	s1 := `{ "$and": [ { "age" : { "$gte" : 20 } },  { "$or" : [  { "postalCode" : 1000 },  { "$and" : [  { "lastname" : "johnson" }, { "firstname" : "ohn"} ] }  ] }  ] } `

	s2 := `{ "age" : { "$gte" : 20 } }`
	s3 := `{ "postalCode" : 1000 }`
	MongoQuery := []string{s0, s1, s2, s3}
	j0 := "(x > 0 AND x < 0 AND x = 3 AND x IN (1,2,3) AND (x != 0 OR x != 0))"
	j1 := `(age >= 20 AND (postalCode = "1000" OR (lastname = "johnson" AND firstname = "ohn")))`
	j2 := "age >= 20"
	j3 := `postalCode = "1000"`
	Results := []string{j0, j1, j2, j3}
	// fmt.Println(s)

	return MongoQuery, Results

}
