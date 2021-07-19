package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func ApiPostWrapper(api func([]byte) ([]byte, error), obj interface{}, result interface{}) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	fmt.Println("xx ", string(body))

	resp, err := api(body)
	if err != nil {
		return err
	}

	if result != nil {
		return json.Unmarshal(resp, result)
	}
	return nil
}

func ApiPostWrapperEx(
	api func([]byte, url.Values) ([]byte, error),
	obj interface{}, params url.Values, result interface{},
) error {
	body, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	resp, err := api(body, params)
	if err != nil {
		return err
	}

	if result != nil {
		return json.Unmarshal(resp, result)
	}
	return nil
}

func ApiGetWrapper(api func(url.Values) ([]byte, error), paramFunc func(url.Values), result interface{}) error {
	params := url.Values{}
	paramFunc(params)
	resp, err := api(params)
	if err != nil {
		return err
	}

	if result != nil {
		return json.Unmarshal(resp, result)
	}
	return nil
}

func ApiGetNullWrapper(api func() ([]byte, error), result interface{}) error {
	resp, err := api()
	if err != nil {
		return err
	}

	if result != nil {
		return json.Unmarshal(resp, result)
	}
	return nil
}
