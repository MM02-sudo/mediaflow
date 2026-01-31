// this code can be used by other programs

package shared

type Request struct {
	Action string `json:"action"`
	Path   string `json:"path"`
}

type Response struct {
	Succes bool        `json:"success"`
	Data   interface{} `json:"data"`
}

// interface{} means that it czan be any type of data
