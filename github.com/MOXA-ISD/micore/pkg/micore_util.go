package micore

import (
    "fmt"
    "time"
    "syscall"
    "os/exec"
)

func GetUTCTime() time.Time {
    t := time.Now()
    loc, err := time.LoadLocation("UTC")
    if err != nil {
        fmt.Println(err)
    }
    UTC := t.In(loc)
    return UTC
}

func GetTimeStamp() string {
    return GetUTCTime().Format(time.RFC3339)
}

func Exec(strCmd string) (int, string) {
    var exitStatus int = 200
    var exitOutput string = ""

    cmd := exec.Command("sh", "-c", strCmd)
    stdout, err := cmd.CombinedOutput()
    if exitOutput = string(stdout); err != nil {
	    if exitErr, ok := err.(*exec.ExitError); ok {
            if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
                switch {
                    case status.Signaled():
		                exitStatus = 500
                        exitOutput = "Command stopped by signal"
                        fmt.Printf("Return signal error: signal code=%d\n", status.Signal())
                }
	      }
        } else {
	        exitStatus = 500
            exitOutput = "Command stopped by unexpected error"
            fmt.Printf("Return other error: %s\n", err)
        }
    }
    return exitStatus, exitOutput
}
