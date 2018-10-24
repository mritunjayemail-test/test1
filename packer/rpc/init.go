package rpc

import "encoding/gob"

func init() {
	// gob.Register(new(*map[string]interface{})) TODO(azr): may be just remove this line
	gob.Register(new(map[string]string))
	gob.Register(make([]interface{}, 0))
	gob.Register(new(BasicError))
}
