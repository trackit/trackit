# Elastic Search

In order to add an index in Elastic search, you must send its template to Elastic Search.

You can learn about ES templates here: [Index templates](https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-templates.html).

### Example
```go
const TemplateNameEBSReport = "ebs-reports"
const TemplateEBSReport = "..."

// put the ElasticSearch index for *-ebs-reports indices at startup.
func init() {
    ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
    res, err := es.Client.IndexPutTemplate(TemplateNameEBSReport).BodyString(TemplateEBSReport).Do(ctx)
    ctxCancel()
    if err != nil {
        jsonlog.DefaultLogger.Error("Failed to put ES index EBSReport.", err)
    } else {
        jsonlog.DefaultLogger.Info("Put ES index EBSReport.", res)
    }
}
```
