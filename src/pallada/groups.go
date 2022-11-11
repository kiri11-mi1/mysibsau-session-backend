package pallada

type Group struct {
	Name string `xmlrpc:"name"`
	Id   int64  `xmlrpc:"id"`
}

type GroupInfo struct {
	Id           int64         `xmlrpc:"id"`
	SessionsIds  []interface{} `xmlrpc:"session_ids"`
	CurrentYears string        `xmlrpc:"cur_year_header"`
}
type GroupInfoArray []GroupInfo

func (gi GroupInfo) ConvertIdsToInt() []int64 {
	var result = make([]int64, 0)
	for _, sessionId := range gi.SessionsIds {
		result = append(result, sessionId.(int64))
	}
	return result
}
