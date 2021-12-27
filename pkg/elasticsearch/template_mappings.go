package elasticsearch

const taskMapping = `{
    "mappings": {
        "properties": {
            "id": {
                "type": "long"
            },
            "objectID": {
                "type": "keyword"
            },
            "objectType": {
                "type": "keyword"
            },
            "product": {
                "type": "keyword"
            },
            "project": {
                "type": "keyword"
            },
            "execution": {
                "type": "keyword"
            },
            "actor": {
                "type": "keyword"
            },
            "action": {
                "type": "keyword"
            },
            "extra": {
                "type": "keyword"
            },
            "date": {
                "type": "date"
            }
        }
    }
}`

func main() {
	ma := make(map[string]string)
	ma["index"] = "zentao_action"
	ma["mappings"] = taskMapping
	EsClient.CreateIndex(ma)
}