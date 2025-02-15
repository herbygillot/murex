package bson

import (
	"github.com/lmorg/murex/lang"
	"github.com/lmorg/murex/lang/stdio"
	"github.com/lmorg/murex/lang/types"
	"go.mongodb.org/mongo-driver/bson"
)

/*func readArray(read stdio.Io, callback func([]byte)) error {
	b, err := read.ReadAll()
	if err != nil {
		return err
	}

	j := make([]interface{}, 0)
	err = bson.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	for i := range j {
		switch j[i].(type) {
		case string:
			callback(bytes.TrimSpace([]byte(j[i].(string))))

		default:
			jBytes, err := bson.Marshal(j[i])
			if err != nil {
				return err
			}
			callback(jBytes)
		}
	}

	return nil
}*/

func readArray(read stdio.Io, callback func([]byte)) error {
	return lang.ArrayTemplate(bson.Marshal, bson.Unmarshal, read, callback)
}

func readArrayWithType(read stdio.Io, callback func([]byte, string)) error {
	return lang.ArrayWithTypeTemplate(types.Json, bson.Marshal, bson.Unmarshal, read, callback)
}

/*func readArrayWithType(read stdio.Io, callback func([]byte, string)) error {
	b, err := read.ReadAll()
	if err != nil {
		return err
	}

	j := make([]interface{}, 0)
	err = bson.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	for i := range j {
		switch j[i].(type) {
		case string:
			callback([]byte(j[i].(string)), types.String)

		default:
			jBytes, err := bson.Marshal(j[i])
			if err != nil {
				return err
			}
			callback(jBytes, typeName)
		}
	}

	return nil
}*/
