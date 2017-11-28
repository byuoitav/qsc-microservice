package qsysremote

//BaseRequest are the common parts of every qsc jsonrpc request
type BaseRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
}

type QSCStatusReport struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Platform    string `json:"Platform"`
		State       string `json:"State"`
		DesignName  string `json:"DesignName"`
		DesignCode  string `json:"DesignCode"`
		IsRedundant bool   `json:"IsRedundant"`
		IsEmulator  bool   `json:"IsEmulator"`
		Status      struct {
			Code   int    `json:"Code"`
			String string `json:"String"`
		} `json:"Status"`
	} `json:"params"`
}

type QSCGetStatusResponse struct {
	BaseRequest
	Result []QSCGetStatusResult `json:"result"`
}
type QSCGetStatusResult struct {
	Name     string
	Value    float64
	String   string
	Position float64
}

type QSCGetStatusRequest struct {
	BaseRequest
	Params []string `json:"params"`
}

type QSCSetStatusRequest struct {
	BaseRequest
	Params QSCSetStatusParams `json:"params"`
}

type QSCSetStatusParams struct {
	Name  string
	Value float64
}

type QSCSetStatusResponse struct {
	BaseRequest
	Result QSCGetStatusResult `json:"result"`
}

func GetGenericSetStatusRequest() QSCSetStatusRequest {
	return QSCSetStatusRequest{BaseRequest: BaseRequest{JSONRPC: "2.0", ID: 1, Method: "Control.Set"}, Params: QSCSetStatusParams{}}
}

func GetGenericGetStatusRequest() QSCGetStatusRequest {
	return QSCGetStatusRequest{BaseRequest: BaseRequest{JSONRPC: "2.0", ID: 1, Method: "Control.Set"}, Params: []string{}}
}
