package repository

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// NewSystem repository for system
func NewSystem(ipAddress string) System {
	return &systemRepo{
		ipAddress: ipAddress,
	}
}

// LogOpenFile send current open file value to datadog
func (r *systemRepo) LogOpenFile() {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
	if err == nil {
		count := bytes.Count(out, []byte("\n"))
		log.Println(count)
	}
}
