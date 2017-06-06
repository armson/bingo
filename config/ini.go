package config

import (
    "bufio"
    "bytes"
    "io"
    "path/filepath"
    "os"
    "strings"
    //"path"
    "errors"
    "strconv"
    "sync"
    "fmt"
)

var iniContainer = map[string]map[string]string{}
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
const DEFAULTSECTION  = "default"

func Int(args ...string) (int64 , error) {
    value, err := String(args...)
    if err != nil { return -1, err }
    i, err := strconv.ParseInt(value, 10, 64)
    if err != nil {
        return -1, err
    } 
    return i, nil
}

func Bool(args ...string) (bool , error) {
    value, err := String(args...)
    if err != nil { return false, err }
    b, err := strconv.ParseBool(value)
    if err != nil {
        return false, err
    } 
    return b, nil
}

func String(args ...string) (string , error) {
    lenArgs := len(args)
    if lenArgs < 1 {
        return "", errors.New("config func String args 1 is must")
    }
    var sectionName string
    if len(args) == 1 { 
        sectionName = DEFAULTSECTION  
    } else {
        sectionName = args[0]
    }
    name := args[lenArgs-1]
    sectionName = strings.ToLower(sectionName)
    name = strings.ToLower(name)

    if _, ok := iniContainer[sectionName]; !ok {
        return "", fmt.Errorf("config section [%s] is not exists" , sectionName)
    }
    if _, ok := iniContainer[sectionName][name]; !ok {
        return "", fmt.Errorf("config section [%s] key [%s] is not exists" , sectionName, name)
    }
    return iniContainer[sectionName][name], nil
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
        sectionName = DEFAULTSECTION
    } else {
        sectionName = args[0]
    }
    name := args[lenArgs-2]
    value := args[lenArgs-1]

    name = strings.TrimSpace(name)
    name = strings.ToLower(name)

    value = strings.TrimSpace(value)
    value = strings.ToLower(value)

    sectionName = strings.TrimSpace(sectionName)
    sectionName = strings.ToLower(sectionName)

    _, sectionExists := iniContainer[sectionName]
    if sectionExists {
        iniContainer[sectionName][name] = value
        return nil
    }
    section := map[string]string{}
    section[name] = value
    iniContainer[sectionName] = section
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
    if len(iniContainer) < 1 {
        return errors.New("Config's data is NULL ")
    }
    return nil
}
func SaveConfig(){
    fileName , _ := String("accessLog")
    fileName = fileName+".configure"
    if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
        fmt.Errorf("Can't create accessLog.configure folder on %v", err)
    }
    file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.ModePerm)
    if err != nil {
        fmt.Errorf("Can't create accessLog.configure file: %v", err)
    }
    defer file.Close()
    runMode,_ := String("runMode")
    pid,_ := String("pid")
    currentPid,_ := String("currentPid")
    configure,_ := String("configFile")

    file.WriteString("[runMode: "+runMode+"]\n")
    file.WriteString("[configure: "+configure+"]\n")
    file.WriteString("[pid: "+pid+"]\n")
    file.WriteString("[currentPid: "+currentPid+"]\n")
    file.WriteString("\n")
    for section, items := range iniContainer {
        file.WriteString("["+section+"]\n")
        for k , item := range items {
            file.WriteString(k+": "+item+"\n")
        }
        file.WriteString("\n")
    }
}
func PrintConfig() {
    fileName , _ := String("accessLog")
    fileName = fileName+".configure"
    file, err := os.Open(fileName)
    defer file.Close()
    if err == nil { 
        configure := make([]byte,409600)
        file.Read(configure)
        fmt.Println(configure)
    }
    fmt.Println(err,"configure")
}

func parseFile(filename string) error {
    currentSection := DEFAULTSECTION
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
    _, sectionExists := iniContainer[sectionName]
    if sectionExists {
        iniContainer[sectionName][string(name)] = string(value)
        return nil
    }
    section := map[string]string{}
    section[string(name)] = string(value)
    iniContainer[sectionName] = section

    return nil
}


