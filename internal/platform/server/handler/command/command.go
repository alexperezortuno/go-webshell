package command

import (
	"database/sql"
	"github.com/alexperezortuno/go-webshell/internal/platform/storage/data_base"
	"github.com/alexperezortuno/go-webshell/tools/environment"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var params = environment.Server()

func CmdHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("command to execute: %s", ctx.Query("cmd"))
		rowsCommands, err := data_base.GetBlackList()

		if err != nil {
			log.Printf("Error getting blacklist: %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Printf("Error closing rows: %s", err)
			}
		}(rowsCommands)

		forbiddenCommand := "forbidden command"
		var forbiddenCommands []string
		var command string
		for rowsCommands.Next() {
			err := rowsCommands.Scan(&command)
			if err != nil {
				log.Printf("Error scanning row: %s", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			forbiddenCommands = append(forbiddenCommands, command)
		}

		for _, str := range strings.Split(forbiddenCommands[0], ",") {
			if strings.Contains(ctx.Query("cmd"), str) {
				ctx.JSON(http.StatusForbidden, gin.H{
					"code":    -2002,
					"message": forbiddenCommand,
				})
				return
			}
		}

		if ctx.Query("cmd") == "" {
			_, err := data_base.Insert(ctx, http.StatusBadRequest)
			if err != nil {
				log.Printf("error inserting command in database: %s", err)
			}

			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    -2000,
				"message": "command is empty",
			})
			return
		}

		commandOutput, err := execute(ctx.Request)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    -2001,
				"message": err.Error(),
			})
			return
		}

		_, err = data_base.Insert(ctx, http.StatusOK)
		if err != nil {
			log.Printf("error inserting command in database: %s", err)
		}

		rows, err := data_base.GetAll()
		if err != nil {
			log.Printf("error getting all commands from database: %s", err)
		}
		defer rows.Close()
		log.Printf("rows: %v", rows)
		for rows.Next() {
			var id int
			var header string
			var ip string
			var command string
			var status int
			var created_at string
			var updated_at string
			err = rows.Scan(&id, &ip, &header, &command, &status, &created_at, &updated_at)
			if err != nil {
				log.Printf("error scanning rows: %s", err)
			}
			log.Printf("id: %d, header: %s, ip: %s, command: %s, status: %d, created_at: %s, updated_at: %s", id, header, ip, command, status, created_at, updated_at)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code":    2000,
			"message": commandOutput,
		})
	}
}

func execute(r *http.Request) (string, error) {
	cmd := r.URL.Query().Get("cmd")
	var out []byte

	actualDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	newDir := params.SafetyZone
	err = os.Chdir(newDir)
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		sh := "cmd.exe"
		out, err := exec.Command(sh, "/K", cmd).Output()
		if err != nil {
			return "", err
		}

		return string(out), nil
	}

	sh := "sh"
	out, err = exec.Command(sh, "-c", cmd).Output()

	if err != nil {
		return "", err
	}

	err = os.Chdir(actualDir)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
