package imports

import "net/http"

func Getlist(request string) (*http.Response, error){
	var err error
	var response *http.Response
	if response, err = http.Head(request); err != nil {
		return response, err
	}
	return response, err

}
