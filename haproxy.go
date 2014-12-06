package main

import(
    "text/template"
    "os"
    "os/exec"
    "log"
    "strconv"
)

var tpl *template.Template = nil
var pid int = -1

func createConfigFile(backends []Backend, templateFile, outputFile string)(error){
    cfgFile,_ := os.Create(outputFile)
    defer cfgFile.Close()

    if(tpl == nil){
        var err error = nil
        tpl, err = template.ParseFiles(templateFile)
        if (err != nil){
            return err
        }
    }

    return tpl.Execute(cfgFile, backends)
}

func reloadHAproxy(command, configFile string)(error){
    var cmd *exec.Cmd = nil
    if pid == -1{
        log.Println("Start HAproxy")
        cmd = exec.Command(command,"-f",configFile)
        go cmd.Wait()
    } else{
        log.Println("Restart HAproxy")
        cmd = exec.Command(command,"-f",configFile,"-sf",strconv.Itoa(pid))
    }

    err := cmd.Run()
    if (err == nil){
        pid = cmd.Process.Pid
        log.Println("New pid: ",pid)
    }
    return err
}

