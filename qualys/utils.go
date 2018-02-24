package qualys

import (
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

func containsString(strList []string, testStr string) bool {
	for _, str := range strList {
		if testStr == str {
			return true
		}
	}
	return false
}

func addURLParameters(urlString string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return urlString, nil
	}

	origURL, err := url.Parse(urlString)
	if err != nil {
		return urlString, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return urlString, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}
