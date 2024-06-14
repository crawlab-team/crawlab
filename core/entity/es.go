package entity

/* ElasticsearchResponseData JSON format
{
  "took" : 6,
  "timed_out" : false,
  "_shards" : {
    "total" : 1,
    "successful" : 1,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : {
      "value" : 60,
      "relation" : "eq"
    },
    "max_score" : 1.0,
    "hits" : [
      {
        "_index" : "test_table",
        "_id" : "c39ad9a2-9a37-49fb-b7ea-f1b55913e0af",
        "_score" : 1.0,
        "_source" : {
          "_tid" : "62524ac7f5f99e7ef594de64",
          "author" : "James Baldwin",
          "tags" : [
            "love"
          ],
          "text" : "“Love does not begin and end the way we seem to think it does. Love is a battle, love is a war; love is a growing up.”"
        }
      }
    ]
  }
}
*/

type ElasticsearchResponseData struct {
	Took    int64 `json:"took"`
	Timeout bool  `json:"timeout"`
	Hits    struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string      `json:"_index"`
			Id     string      `json:"_id"`
			Score  float64     `json:"_score"`
			Source interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
