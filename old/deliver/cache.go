package deliver

import (
	"bytes"
	"log"
	"os/exec"
	"static-proxy/settings"
	"strconv"
	"strings"
)

func fullCache() bool {
	out, err := exec.Command("du", "-sk", "./cache").Output()
	if err != nil {
		log.Println(err)
		return false
	}
	output := string(out)
	arr := strings.Split(output, "\t")
	if len(arr) > 1 {
		if s, err := strconv.ParseInt(arr[0], 10, 32); err == nil {
			if s > settings.Config.CacheSize {
				return true
			}
		}
		return false
	}

	return false
}

func flushOldest() error {
	var result, stdErr bytes.Buffer

	findC := exec.Command("find", "./cache", "-type", "f", "-print0")
	findStdout, err := findC.StdoutPipe()
	if err != nil {
		return err
	}
	xargsC := exec.Command("xargs", "-0", "stat", "-f", "\"%a %N\"")
	xargsC.Stdin = findStdout
	xargsStdout, err := xargsC.StdoutPipe()
	if err != nil {
		return err
	}
	sortC := exec.Command("sort")
	sortC.Stdin = xargsStdout
	sortStdout, err := sortC.StdoutPipe()
	if err != nil {
		return err
	}
	headC := exec.Command("head", "-20")
	headC.Stdin = sortStdout
	headC.Stdout = &result
	headC.Stderr = &stdErr

	if err = headC.Start(); err != nil {
		return err
	}
	if err = sortC.Start(); err != nil {
		return err
	}
	if err = xargsC.Start(); err != nil {
		return err
	}
	if err = findC.Run(); err != nil {
		log.Println(stdErr.String())
		return err
	}
	if err = xargsC.Wait(); err != nil {
		return err
	}
	if err = sortC.Wait(); err != nil {
		return err
	}
	if err = headC.Wait(); err != nil {
		return err
	}

	lines := strings.Split(result.String(), "\n")
	if len(lines) > 0 {
		var (
			files  = []string{"-f"}
			stdErr bytes.Buffer
		)

		for _, val := range lines {
			arr := strings.Split(val, " ")
			if len(arr) > 1 {
				file := strings.Replace(arr[1], "\"", "", -1)
				files = append(files, file)
			}
		}
		if len(files) > 1 {
			rmC := exec.Command("rm", files...)
			rmC.Stderr = &stdErr
			if err := rmC.Run(); err != nil {
				log.Println(stdErr.String())
				return err
			}
		}

	}
	return nil
}
