package serverError

import "fmt"

var BadInternalApiCall = fmt.Errorf("Failed internal API call")
