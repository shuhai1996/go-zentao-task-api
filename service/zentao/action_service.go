package zentao

import (
	"github.com/olivere/elastic/v7"
	"go-zentao-task-api/model/zentao"
)

var es = zentao.NewEsAction()

func GetActionList (page int, size int, query string) (int64,[]*elastic.SearchHit) {
	data :=es.List(page, size, query)
	total:= data.TotalHits.Value
	hits := data.Hits

	return total, hits
}

func FindOne (id string) (*elastic.GetResult, error) {
	return es.Find(id)
}

func DeleteOne (id string) (string, error) {
	return es.Delete(id)
}


