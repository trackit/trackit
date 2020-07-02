package tagging

import (
	"context"
	"time"

	"github.com/trackit/jsonlog"
	"github.com/trackit/trackit/es"
)

const typeTaggingReport = "tagging-reports"
const indexPrefixTaggingReport = "tagging-reports"
const templateNameTaggingReport = "tagging-reports"

// put the ElasticSearch index for *-tagging-reports indices at startup.
func initIndexTemplate() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := es.Client.IndexPutTemplate(templateNameTaggingReport).BodyString(templateTaggingReport).Do(ctx)
	if err != nil {
		jsonlog.DefaultLogger.Error("Failed to put ES indext tagging-reports.", err)
		ctxCancel()
	} else {
		jsonlog.DefaultLogger.Info("Put ES index tagging-reports.", res)
		ctxCancel()
	}
}

const templateTaggingReport = `
{
    "template":"*-tagging-reports",
    "version":2,
    "mappings":{
        "tagging-reports":{
            "properties":{
                "account":{
                    "type":"keyword"
                },
                "region":{
                    "type":"keyword"
                },
                "reportDate":{
                    "type":"date"
                },
                "resourceId":{
                    "type":"keyword"
                },
                "resourceType":{
                    "type":"keyword"
                },
                "tags":{
                    "type":"nested",
                    "properties":{
                        "key":{
                            "type":"keyword"
                        },
                        "value":{
                            "type":"keyword"
                        }
                    }
                },
                "url":{
                    "type":"keyword"
                }
			},
			"_all": {
				"enabled": false
			},
			"date_detection": false,
			"numeric_detection": false
        }
    }
}
`
