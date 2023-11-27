// Copyright 2019 TiKV Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/pingcap/errors"
	"github.com/tikv/pd/pkg/errs"
)

// ValidateURLWithScheme checks the format of the URL.
func ValidateURLWithScheme(rawURL string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return errors.Errorf("%s has no scheme", rawURL)
	}
	return nil
}

// parseUrls parse a string into multiple urls.
func parseUrls(s string) ([]url.URL, error) {
	items := strings.Split(s, ",")
	urls := make([]url.URL, 0, len(items))
	for _, item := range items {
		u, err := url.Parse(item)
		if err != nil {
			return nil, errs.ErrURLParse.Wrap(err).GenWithStackByCause()
		}

		urls = append(urls, *u)
	}

	return urls, nil
}

// FlattenConfigItems flatten config to map.
func FlattenConfigItems(nestedConfig interface{}) (map[string]interface{}, error) {
	jsonValue, err := json.Marshal(nestedConfig)
	if err != nil {
		return nil, err
	}
	nestedConfigMap := make(map[string]interface{})
	err = json.Unmarshal(jsonValue, &nestedConfigMap)
	if err != nil {
		return nil, err
	}
	flatMap := make(map[string]interface{})
	flatten(flatMap, nestedConfigMap, "")
	return flatMap, nil
}

func flatten(flatMap map[string]interface{}, nested interface{}, prefix string) {
	switch nested := nested.(type) {
	case map[string]interface{}:
		for k, v := range nested {
			path := k
			if prefix != "" {
				path = prefix + "." + k
			}
			flatten(flatMap, v, path)
		}
	default: // don't flatten arrays
		flatMap[prefix] = nested
	}
}
