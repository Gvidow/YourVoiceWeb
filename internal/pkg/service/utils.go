package service

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gvidow/YourVoiceWeb/internal/pkg/repository/chat"
)

type M map[string]string

func getIntParam(param string) (int, bool) {
	param = strings.Trim(param, "/")
	param = strings.Split(param, "/")[0]
	res, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, false
	} else {
		return int(res), true
	}
}

func getStringParam(param string) string {
	param = strings.Trim(param, "/")
	return strings.Split(param, "/")[0]
}

func readChatSettings(r io.Reader) (*chat.Setting, error) {
	d := json.NewDecoder(r)
	res := &chat.Setting{}
	err := d.Decode(res)
	if err != nil {
		return nil, fmt.Errorf("decode chat settings: %w", err)
	}
	return res, nil
}

func unmarshalBody(r io.Reader) (map[string]string, error) {
	d := json.NewDecoder(r)
	var m map[string]string
	err := d.Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("decode request body: %w", err)
	}
	return m, nil
}

func responseOK(comment string) M {
	return M{
		"status":  "ok",
		"comment": comment,
	}
}

func responseErr(code, message string) M {
	return M{
		"status":  "error",
		"code":    code,
		"message": message,
	}
}
