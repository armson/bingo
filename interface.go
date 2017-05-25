package bingo

func InterfaceToInt64(a interface{}) (int64) {
    switch a.(type) {
        case int64:
            if i, ok := a.(int64); ok {
                return i
            }
        case int32:
            if i, ok := a.(int32); ok {
                return int64(i)
            }
        case int16:
            if i, ok := a.(int16); ok {
                return int64(i)
            }
        case int8:
            if i, ok := a.(int8); ok {
                return int64(i)
            }
        case int:
            if i, ok := a.(int); ok {
                return int64(i)
            } 
        case float32:
            if i, ok := a.(float32); ok {
                return int64(i)
            }
        case float64:
            if i, ok := a.(float64); ok {
                return int64(i)
            }                    
        default:
            return 0     
    }
    return 0
}