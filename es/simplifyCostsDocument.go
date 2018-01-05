//   Copyright 2017 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/trackit/jsonlog"

	"gopkg.in/olivere/elastic.v5"
)

// SimplifiedCostsDocument contains the data necessary to show a clean version
// of the costs breakdowns returned from the ElasticSearch.
type SimplifiedCostsDocument struct {
	Key          string
	Children     []SimplifiedCostsDocument
	ChildrenKind string
	HasValue     bool
	Value        float64
}

// ToJsonable returns the simplified costs document as a map that can easily be
// marshaled to JSON.
func (scd SimplifiedCostsDocument) ToJsonable() map[string]interface{} {
	children := make(map[string]interface{})
	for _, c := range scd.Children {
		if c.HasValue {
			children[c.Key] = c.Value
		} else {
			children[c.Key] = c.ToJsonable()
		}
	}
	return map[string]interface{}{
		scd.ChildrenKind: children,
	}
}

type value = map[string]interface{}
type aggregation = map[string]interface{} // aggregation has buckets
type bucket = map[string]interface{}      // bucket has aggregations and values

const (
	BucketPrefix         = "by-"
	BucketKeyKey         = "key"
	AggBucketKey         = "buckets"
	BucketKeyAsStringKey = "key_as_string"
	BucketValueKey       = "value"
	BucketValueValueKey  = "value"
)

var (
	ErrNoSingleRootAggregation   = errors.New("document does not have a single aggregation at its root")
	ErrNoSingleAggregationBranch = errors.New("document's aggregations branch")
	ErrFailedJsonParsing         = errors.New("failed to parse JSON document")
	ErrKeyNotFound               = errors.New("could not find 'key' field")
	ErrNoAggregation             = errors.New("found no next aggregation and no value")
)

func SimplifyCostsDocument(ctx context.Context, sr *elastic.SearchResult) (SimplifiedCostsDocument, error) {
	var scdz SimplifiedCostsDocument
	if len(sr.Aggregations) == 1 {
		for k, v := range sr.Aggregations {
			if v != nil {
				return simplifyCostsDocumentWithSingleAggregation(ctx, k, sr.Aggregations[k])
			}
		}
	}
	return scdz, ErrNoSingleRootAggregation
}

func simplifyCostsDocumentWithSingleAggregation(ctx context.Context, rootAgg string, rm *json.RawMessage) (SimplifiedCostsDocument, error) {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	var parsedDocument bucket
	var scdz SimplifiedCostsDocument
	err := json.Unmarshal(*rm, &parsedDocument)
	if err != nil {
		logger.Error("Failed to parse JSON costs document.", err.Error())
		return scdz, ErrFailedJsonParsing
	} else {
		return simplifyCostsDocumentRec(ctx, map[string]interface{}{rootAgg: parsedDocument}, true)
	}
}

func simplifyCostsDocumentRec(ctx context.Context, doc bucket, root bool) (SimplifiedCostsDocument, error) {
	var scd SimplifiedCostsDocument
	if !root {
		var err error
		scd.Key, err = getKey(doc)
		if err != nil {
			return scd, err
		}
	}
	if value, ok := getValue(ctx, doc); ok {
		scd.HasValue = true
		scd.Value = value
		return scd, nil
	} else if childrenKind, children, err := getChildren(ctx, doc); err == nil {
		scd.Children = children
		scd.ChildrenKind = childrenKind
		return scd, nil
	} else {
		return scd, err
	}
}

func getKey(doc bucket) (string, error) {
	if key, ok := doc[BucketKeyKey]; ok {
		if tkey, ok := key.(string); ok {
			return tkey, nil
		}
	}
	if key, ok := doc[BucketKeyAsStringKey]; ok {
		if tkey, ok := key.(string); ok {
			return tkey, nil
		}
	}
	return "", ErrKeyNotFound
}

func getValue(ctx context.Context, doc bucket) (float64, bool) {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	if value, ok := doc[BucketValueKey]; ok {
		if tvalue, ok := value.(map[string]interface{}); ok {
			if valuevalue, ok := tvalue[BucketValueValueKey]; ok {
				if tvaluevalue, ok := valuevalue.(float64); ok {
					return tvaluevalue, true
				} else {
					logger.Warning(fmt.Sprintf("Found non float value: %[1]T %[1]v.", valuevalue), nil)
				}
			}
		}
	}
	return 0, false
}

func getChildren(ctx context.Context, doc bucket) (string, []SimplifiedCostsDocument, error) {
	var logger = jsonlog.LoggerFromContextOrDefault(ctx)
	if childKey, err := getChildKey(ctx, doc); err != nil {
		logger.Error("Failed to get child key.", err.Error())
		logger.Debug("Document is.", doc)
		return "", nil, err
	} else if childAgg, ok := doc[BucketPrefix+childKey].(map[string]interface{}); !ok {
		logger.Error(fmt.Sprintf("Failed to get buckets: value under '%s' is not an aggregation.", childKey), nil)
		logger.Debug("Document is.", doc)
		return "", nil, ErrFailedJsonParsing
	} else if childAggsBuckets, ok := childAgg[AggBucketKey]; !ok {
		logger.Error(fmt.Sprintf("Failed to get buckets: value under '%s' does not have '%s' field.", childKey, AggBucketKey), nil)
		logger.Debug("Document is.", doc)
		return "", nil, ErrFailedJsonParsing
	} else if children, ok := childAggsBuckets.([]interface{}); !ok {
		logger.Error(fmt.Sprintf("Failed to get buckets: value under '%s.%s' is not a slice.", childKey, AggBucketKey), nil)
		logger.Debug("Document is.", doc)
		return "", nil, ErrFailedJsonParsing
	} else {
		cs := make([]SimplifiedCostsDocument, len(children))
		for i, child := range children {
			tchild, ok := child.(bucket)
			if !ok {
				logger.Error("Child under '%s' is not a bucket.", nil)
				logger.Debug("Document is.", doc)
				return "", nil, ErrFailedJsonParsing
			}
			if cs[i], err = simplifyCostsDocumentRec(ctx, tchild, false); err != nil {
				return "", nil, err
			}
		}
		return childKey, cs, nil
	}
}

func getChildKey(ctx context.Context, doc map[string]interface{}) (string, error) {
	var childKey string
	for k := range doc {
		if strings.HasPrefix(k, BucketPrefix) {
			if childKey != "" {
				return "", ErrNoSingleAggregationBranch
			} else {
				childKey = strings.TrimPrefix(k, BucketPrefix)
			}
		}
	}
	if childKey == "" {
		return childKey, ErrNoAggregation
	}
	return childKey, nil
}
