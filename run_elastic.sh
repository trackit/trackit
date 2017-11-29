#!/bin/bash

set -ex

docker run -d --rm --name=es-go-tests -p 9200:9200  -e "http.host=0.0.0.0" -e "transport.host=127.0.0.1" -e "bootstrap.memory_lock=true" -e "ES_JAVA_OPTS=-Xms1g -Xmx1g" docker.elastic.co/elasticsearch/elasticsearch:5.5.3 elasticsearch -Expack.security.enabled=false -Enetwork.host=_local_,_site_ -Enetwork.publish_host=_local_
sleep 20
wget "https://s3-us-west-2.amazonaws.com/trackit-public-artifacts/elasticsearch/awsdetailedlineitem/es_data.ndjson"
wget "https://s3-us-west-2.amazonaws.com/trackit-public-artifacts/elasticsearch/awsdetailedlineitem/mapping.json"
curl -XPUT 'localhost:9200/awsdetailedlineitem?pretty' -H 'Content-Type: application/json' -d'
{}
'
curl -XPUT 'localhost:9200/awsdetailedlineitem/_mapping/a_ws_detailed_lineitem?pretty' -H 'Content-Type: application/json' -d"@mapping.json"
curl -s -H 'Content-Type: application/x-ndjson' -XPOST 'localhost:9200/_bulk' --data-binary "@es_data.ndjson"
rm -f es_data.ndjson mapping.json
