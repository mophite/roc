package jsonc

import jsoniter "github.com/json-iterator/go"

var JSCodec = &JsonCodec{jsonCodec: jsoniter.ConfigFastest}

type JsonCodec struct {
	jsonCodec jsoniter.API
}

func (j *JsonCodec) Encode(req interface{}) ([]byte, error) {
	b, err := j.jsonCodec.Marshal(req)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (j *JsonCodec) Decode(b []byte, rsp interface{}) error {
	return j.jsonCodec.Unmarshal(b, rsp)
}
