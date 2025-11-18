package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func List_file(path string, reg string, max int) ([]string, error) {
	filelist := []string{}

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return filelist, err
	}
	sep := "/"
	for _, fi := range dir {
		if fi.IsDir() {
			s, err := List_file(path+sep+fi.Name(), reg, max)
			if err == nil {
				for _, a := range s {
					filelist = append(filelist, a)
					if len(filelist) >= max {
						return filelist, nil
					}
				}
			}
		} else {
			if len(reg) > 0 {
				m, _ := regexp.MatchString(reg, fi.Name())
				if m {
					filelist = append(filelist, path+sep+fi.Name())
				}
			} else {
				if strings.HasSuffix(fi.Name(), ".tmp") || strings.HasSuffix(fi.Name(), ".temp") || strings.HasSuffix(fi.Name(), ".dealling") {
				} else {
					filelist = append(filelist, path+sep+fi.Name())
				}
			}
			if len(filelist) >= max {
				return filelist, nil
			}
		}
	}
	return filelist, nil
}

func Isexist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return true
		}
	} else {
		return true
	}
}

func Runcmd(cmd string) ([]byte, error) {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/C", cmd)
	} else {
		c = exec.Command("/bin/sh", "-c", cmd)
	}
	return c.Output()
}

func Decimal(value float64) float64 {
	val := int64(value * 100)
	if val/10 == 0 {
		val += 1
	}
	return float64(val) / 100
}
