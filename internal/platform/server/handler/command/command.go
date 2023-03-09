package command

import (
    "bufio"
    "github.com/gin-gonic/gin"
    "net/http"
    "os"
    "os/exec"
)

func CmdHandler() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        //var cmd *exec.Cmd
        com := ctx.Query("cmd")
        commandOutput, err := execute(ctx.Request)

        if err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "code":    -2001,
                "message": err.Error(),
            })
            return
        }

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

func execute(r *http.Request) ([]byte, error) {
    //var commd []string
    //
    //commd = strings.Split(inputCommd, " ")
    var out []byte
    //if len(commd) > 1 {
    //    out, _ = exec.Command(commd[0], commd[1:]...).Output()
    //} else {
    //    out, _ = exec.Command(commd[0]).Output()
    //}
    //
    //fmt.Printf("%s", out)
    //
    //return out
    cmd := exec.Command("bash")
    cmd.Stdin = r.Body
    cmd.Stderr = os.Stderr
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return out, err
    }

    if err := cmd.Start(); err != nil {
        return out, err
    }

    defer cmd.Wait()
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        //fmt.Fprintln(w, scanner.Text())
        out = append(out, scanner.Bytes()...)
    }

    return out, nil
}
