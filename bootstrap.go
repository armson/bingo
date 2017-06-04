package bingo

import(
    "flag"
    "github.com/armson/bingo/config"
    "os"
    "fmt"
    "path/filepath"
    "strconv"
    "errors"
)

var usageStr = `Usage:%s [options]
Server Options:
    -c, --config <file>              Configuration file path
    -m, --mode [debug|release]  Set run mode (default:debug)
Common Options:
    -h, --help                       Show this message
    -v, --version                    Show version
`
func usage(){
    fmt.Printf(usageStr, os.Args[0])
    os.Exit(0)
}

func createPIDFile() error {
    pidPath, err := config.String("pid")
    if err != nil {
        return err
    }
     _, er := os.Stat(pidPath)
     pidOverride := true
    if os.IsNotExist(er) || pidOverride == true {
        currentPid := os.Getpid()
        if err := os.MkdirAll(filepath.Dir(pidPath), os.ModePerm); err != nil {
            return fmt.Errorf("Can't create PID folder on %v", err)
        }

        file, err := os.Create(pidPath)
        if err != nil {
            return fmt.Errorf("Can't create PID file: %v", err)
        }
        defer file.Close()
        if _, err := file.WriteString(strconv.FormatInt(int64(currentPid), 10)); err != nil {
            return fmt.Errorf("Can'write PID information on %s: %v", pidPath, err)
        }
    } else {
        return fmt.Errorf("%s already exists", pidPath)
    }
    return nil
}

func init(){
    var isShowVersion bool
    var configFile string
    var mode string
    flag.StringVar(&configFile, "c", "conf/app.conf", "Configuration file path.")
    flag.StringVar(&configFile, "config", "conf/app.conf", "Configuration file path.")
    flag.StringVar(&mode, "m", "", "Set run mode.")
    flag.StringVar(&mode, "mode", "", "Set run mode.")
    flag.BoolVar(&isShowVersion, "version", false, "Print version information.")
    flag.BoolVar(&isShowVersion, "v", false, "Print version information.")
    flag.Usage = usage
    flag.Parse()

    args := os.Args
    if len(args) < 2 { 
        usage() 
    }
    if isShowVersion {
        PrintVersion()
        os.Exit(0)  
    }
       
    err := config.Load(configFile)
    if err != nil {
        fmt.Println(err)
        os.Exit(0)
    }

    if err = createPIDFile(); err != nil {
        fmt.Println(err)
        os.Exit(0)
    }
   
    if mode == "" {
        if configMode, err := config.String("runMode"); err == nil {
            mode = configMode
        }
    }
    if mode != DebugMode && mode != ReleaseMode {
        fmt.Println(errors.New("run mode unknown ,it must be debug or release"))
        os.Exit(0)
    }
    SetMode(mode)



}





