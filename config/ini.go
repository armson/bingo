package config

import (
    "bufio"
    "bytes"
    "io"
    "io/ioutil"
    "os"
    "strings"
    "path"
    "errors"
    "strconv"
    "sync"
)

type Container struct {
    data        map[string]map[string]string
    files       []string
}

var (
    bCommentA       byte = '#'
    bCommentB       byte = ';'
    bEqual          []byte = []byte{'='}
    bQuote          string = "'"
    bDoubleQuote    string = "\""
    bSectionStart   byte = '['
    bSectionEnd     byte = ']'
    lineBreak       byte = '\n'
    Ini *Container
    mutex sync.Mutex
)
const (
    INIFILEDIR  = "conf/"
    INIFILESUFFIX = ".conf"
    DEFAULTSECTION  = "default"
)

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
    
    if _, ok := Ini.data[sectionName]; !ok {
        return "", errors.New("config section is not exists")
    }
    if _, ok := Ini.data[sectionName][name]; !ok {
        return "", errors.New("config key is not exists")
    }
    return Ini.data[sectionName][name], nil
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

    _, sectionExists := Ini.data[sectionName]
    if sectionExists {
        Ini.data[sectionName][string(name)] = string(value)
        return nil
    }
    section := map[string]string{}
    section[string(name)] = string(value)
    Ini.data[sectionName] = section
    
    return nil
}

func init() {
    Ini = &Container{
        data:map[string]map[string]string{},
        files:[]string{},
    }
    Ini = Ini.load(INIFILEDIR)
}
func (this *Container) load(dirname string) *Container {
    files := this.getIniFiles(dirname)
    if len(files) < 1 { return this }
    for _, file := range files {
        this.parseFile(file)
    }
    return this
}

func (this *Container) parseFile(filename string) *Container {
    currentSection := DEFAULTSECTION
    fp, err := os.Open(filename)
    if err != nil {
        return this
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
        err = this.parseLine(currentSection, line)
        if err != nil {
            panic("Load conf file is wrong")
        }
    }
    return this
}

func (this *Container) parseLine(sectionName string, line []byte) error {
    if line[0] == bCommentA || line[0] == bCommentB {
        return nil
    }
    parts := bytes.SplitN(line, bEqual, 2)
    name, value := parts[0], parts[1]
    name = bytes.TrimSpace(name)
    name = bytes.Trim(name, bQuote)
    name = bytes.Trim(name, bDoubleQuote)

    value = bytes.TrimSpace(value)
    value = bytes.Trim(value, bQuote)
    value = bytes.Trim(value, bDoubleQuote)

    _, sectionExists := this.data[sectionName]
    if sectionExists {
        this.data[sectionName][string(name)] = string(value)

        return nil
    }
    section := map[string]string{}
    section[string(name)] = string(value)
    this.data[sectionName] = section

    return nil
}

func (this *Container) getIniFiles(dirname string) []string {
    files , err := ioutil.ReadDir(dirname)
    if err != nil || len(files) < 1 { 
        return this.files 
    }
    for _, file := range files {
        filename := strings.ToLower(file.Name())
        if !file.IsDir() && path.Ext(filename) == INIFILESUFFIX {
            this.files = append(this.files, INIFILEDIR+filename)
        }
    }
    return this.files
}


