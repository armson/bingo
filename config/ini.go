package config

import (
    "bufio"
    "bytes"
    "io"
    "os"
    "strings"
    "errors"
    "strconv"
    "sync"
    "fmt"
    "time"
)

var (
    bCommentA       byte = '#'
    bCommentB       byte = ';'
    bEqual          []byte = []byte{'='}
    bQuote          string = "'"
    bDoubleQuote    string = "\""
    bSectionStart   byte = '['
    bSectionEnd     byte = ']'
    lineBreak       byte = '\n'
    mutex sync.Mutex
    once  sync.Once
)
const defaultSection  = "default"

func Int(args ...string) (int) {
    value, err := Get(args...)
    if err != nil { return -1 }
    i, err := strconv.ParseInt(value, 10, 64)
    if err != nil {
        return -1
    }
    return int(i)
}

func Bool(args ...string) (bool) {
    value, err := Get(args...)
    if err != nil { return false}
    b, err := strconv.ParseBool(value)
    if err != nil {
        return false
    } 
    return b
}

func String(args ...string) (string) {
    value, err := Get(args...)
    if err != nil { return ""}
    return value
}

func Float(args ...string) (float64) {
    value, err := Get(args...)
    if err != nil { return -1 }
    i, err := strconv.ParseFloat(value, 64)
    if err != nil {
        return -1
    } 
    return i
}

func Slice(args ...string) (slice []string) {
	value, err := Get(args...)
	if value = strings.TrimSpace(value); err != nil || value == "" {
		return nil
	}
	for _ , arr := range strings.Split(value, ",") {
		if arr = strings.TrimSpace(arr); arr != "" {
			slice = append(slice, arr)
		}
	}
	return
}

func Map(sectionName string) (map[string]string) {
    sectionName = strings.ToLower(sectionName)
    if _, ok := defaultConfigs[sectionName]; !ok {
        return map[string]string{}
    }
    return defaultConfigs[sectionName]
}

func Default() (map[string]map[string]string) {
    return defaultConfigs
}


func Time(args ...string) (time.Duration) {
    value, err := Get(args...)
    if err != nil { return 0 }
    if value == "" { return 0 }
    t , err := time.ParseDuration(value)
    if err != nil { return 0 }
    return t
}


func Get(args ...string) (string , error) {
    lenArgs := len(args)
    if lenArgs < 1 {
        return "", errors.New("config func String args 1 is must")
    }
    var sectionName string
    if len(args) == 1 { 
        sectionName = defaultSection
    } else {
        sectionName = args[0]
    }
    name := args[lenArgs-1]
    sectionName = strings.ToLower(sectionName)
    name = strings.ToLower(name)

    if _, ok := defaultConfigs[sectionName]; !ok {
        return "", fmt.Errorf("config section [%s] is not exists" , sectionName)
    }
    if _, ok := defaultConfigs[sectionName][name]; !ok {
        return "", fmt.Errorf("config section [%s] key [%s] is not exists" , sectionName, name)
    }
    return defaultConfigs[sectionName][name], nil
}

func Set(args ...string) (error) {
    mutex.Lock()
    defer mutex.Unlock()
    
    lenArgs := len(args)
    if lenArgs < 2 {
        return errors.New("config func Set args 2 is must")
    }
    
    var sectionName string
    if lenArgs == 2 { 
        sectionName = defaultSection
    } else {
        sectionName = args[0]
    }
    name := args[lenArgs-2]
    value := args[lenArgs-1]

    name = strings.TrimSpace(name)
    name = strings.ToLower(name)

    value = strings.TrimSpace(value)

    sectionName = strings.TrimSpace(sectionName)
    sectionName = strings.ToLower(sectionName)

    _, sectionExists := defaultConfigs[sectionName]
    if sectionExists {
        defaultConfigs[sectionName][name] = value
        return nil
    }
    section := map[string]string{}
    section[name] = value
    defaultConfigs[sectionName] = section
    return nil
}

func Load(filename string) (error) {
    if len(filename) < 1 { 
        return errors.New("Config file is required")
    }
    if err := parseFile(filename); err != nil {
        return err
    }
    if err := checkConfig(); err != nil {
        return err
    }
    return nil
}
func checkConfig()(error){
    if len(defaultConfigs) < 1 {
        return errors.New("Config's data is NULL ")
    }
    return nil
}

func parseFile(filename string) error {
    currentSection := defaultSection
    fp, err := os.Open(filename)
    if err != nil { 
        return err 
    }
    reader := bufio.NewReader(fp)
    for {
        line, err := reader.ReadBytes(lineBreak)
        if err == io.EOF { 
            break 
        }
        line = bytes.TrimSpace(line)
        if len(line) == 0 || line[0] == lineBreak || line[0] == bCommentA  || line[0] == bCommentB  {
            continue
        }
        // parse section
        if line[0] == bSectionStart {
            if line[len(line)-1] == bSectionEnd {
                currentSection = string(line[1 : len(line)-1])
                continue
            }
        }
        // parse item
        err = parseLine(currentSection, line)
        if err != nil {
            return err
        }
    }
    return nil
}

func parseLine(sectionName string, line []byte) error {
    if line[0] == bCommentA || line[0] == bCommentB {
        return nil
    }
    parts := bytes.SplitN(line, bEqual, 2)
    name, value := parts[0], parts[1]
    name = bytes.TrimSpace(name)
    name = bytes.Trim(name, bQuote)
    name = bytes.Trim(name, bDoubleQuote)
    name = bytes.ToLower(name)

    value = bytes.TrimSpace(value)
    value = bytes.Trim(value, bQuote)
    value = bytes.Trim(value, bDoubleQuote)

    sectionName = strings.ToLower(sectionName)
    _, sectionExists := defaultConfigs[sectionName]
    if sectionExists {
        defaultConfigs[sectionName][string(name)] = string(value)
        return nil
    }
    section := map[string]string{}
    section[string(name)] = string(value)
    defaultConfigs[sectionName] = section

    return nil
}


