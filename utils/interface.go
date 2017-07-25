package utils

import(
	"strconv"
	"errors"
)

type binInterface string
var Interface *binInterface

func (_ *binInterface) Int(a interface{}) (int, error) {
    switch a.(type) {
        case int64:
            if i, ok := a.(int64); ok {
                return int(i),nil
            }
        case int32:
            if i, ok := a.(int32); ok {
                return int(i),nil
            }
        case int16:
            if i, ok := a.(int16); ok {
                return int(i),nil
            }
        case int8:
            if i, ok := a.(int8); ok {
                return int(i),nil
            }
        case int:
            if i, ok := a.(int); ok {
                return i,nil
            } 
        case float32:
            if i, ok := a.(float32); ok {
                return int(i),nil
            }
        case float64:
            if i, ok := a.(float64); ok {
                return int(i),nil
            }
		case string:
			if str, ok := a.(string); ok {
				i, err := strconv.ParseInt(str, 10, 0)
				if err == nil {
					return int(i),nil
				}
			}
        default:
            return 0,errors.New("Don't change to Int")
    }
    return 0,errors.New("Don't change to Int")
}
func (_ *binInterface) Int64(a interface{}) (int64, error) {
	switch a.(type) {
	case int64:
		if i, ok := a.(int64); ok {
			return i,nil
		}
	case int32:
		if i, ok := a.(int32); ok {
			return int64(i),nil
		}
	case int16:
		if i, ok := a.(int16); ok {
			return int64(i),nil
		}
	case int8:
		if i, ok := a.(int8); ok {
			return int64(i),nil
		}
	case int:
		if i, ok := a.(int); ok {
			return int64(i),nil
		}
	case float32:
		if i, ok := a.(float32); ok {
			return int64(i),nil
		}
	case float64:
		if i, ok := a.(float64); ok {
			return int64(i),nil
		}
	case string:
		if str, ok := a.(string); ok {
			i, err := strconv.ParseInt(str, 10, 0)
			if err == nil {
				return i,nil
			}
		}
	default:
		return 0,errors.New("Don't change to Int64")
	}
	return 0,errors.New("Don't change to Int64")
}
func (_ *binInterface) String(a interface{}) (string, error) {
	switch a.(type) {
	case int64:
		if i, ok := a.(int64); ok {
			return strconv.FormatInt(i, 10),nil
		}
	case int32:
		if i, ok := a.(int32); ok {
			return strconv.FormatInt(int64(i), 10),nil
		}
	case int16:
		if i, ok := a.(int16); ok {
			return strconv.FormatInt(int64(i), 10),nil
		}
	case int8:
		if i, ok := a.(int8); ok {
			return strconv.FormatInt(int64(i), 10),nil
		}
	case int:
		if i, ok := a.(int); ok {
			return strconv.FormatInt(int64(i), 10),nil
		}
	case float32:
		if i, ok := a.(float32); ok {
			return strconv.FormatInt(int64(i), 10),nil
		}
	case float64:
		if i, ok := a.(float64); ok {
			return strconv.FormatInt(int64(i), 10),nil
		}
	case string:
		if i, ok := a.(string); ok {
			return i,nil
		}
	default:
		return "",errors.New("Don't change to String")
	}
	return "",errors.New("Don't change to String")
}