package misc

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func JSON(x interface{}) string {
	js, err := json.Marshal(&x)
	if err != nil {
		x = map[string]interface{}{
			"object": fmt.Sprintf("%#v", x),
		}
		js, _ = json.Marshal(&x)
	}
	return string(js)
}

func SHA(js []byte) string {
	hasher := sha1.New()
	hasher.Write(js)
	return hex.EncodeToString(hasher.Sum(nil))
}
