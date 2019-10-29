package types

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"time"
)

const (
	scriptOutputDir = "/curl-scripts/"
	outputDir       = "/curl-scripts/output/"
)

var (
	re = regexp.MustCompile(`\s+`)
)

func init() {
	RegisterAttacker("curl", &Curl{})
}

type Curl struct {
	AttackConfig
}

func (b Curl) ConstructOutputFile() string {
	return ">> " + outputDir + re.ReplaceAllString(b.Name, "-") + ".log"
}

func (b *Curl) PrepareAttack() (string, error) {
	if err := b.Validate(); err != nil {
		return "", err
	}

	fileName := scriptOutputDir + b.Name + ".sh"
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return "", err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	writer.WriteString("#!/bin/bash\n\n")
	writer.WriteString("set -eux\n\n")

	writer.WriteString("# exploit\n")
	for _, e := range b.HTTPExploit {
		curlString := fmt.Sprintf("curl %s %s %s %s %s",
			e.ConstructMethod(), e.ConstructHeader(), e.ConstructData(), e.ConstructURI(b.AttackConfig), b.ConstructOutputFile())
		writer.WriteString(curlString)
		writer.WriteString("\n")
	}
	writer.WriteString("\n")

	writer.WriteString("# payload\n")
	for _, e := range b.HTTPPayload {
		curlString := fmt.Sprintf("curl %s %s %s %s %s",
			e.ConstructMethod(), e.ConstructHeader(), e.ConstructData(), e.ConstructURI(b.AttackConfig), b.ConstructOutputFile())
		writer.WriteString(curlString)
		writer.WriteString("\n")
	}
	writer.WriteString("\n")

	writer.Flush()

	return fileName, nil
}

func (b *Curl) LaunchAttack(attackScript string) (err error) {
	var stdin io.WriteCloser

	if len(attackScript) == 0 {
		return errors.New("attack script is not specified")
	}

	// launch attack
	cmd := exec.Command(attackScript)

	if b.IsReverseShell {
		stdin, err = cmd.StdinPipe()
		if err != nil {
			return errors.New("errors in reverse shell")
		}
	}

	cmd.Start()

	if b.IsReverseShell {
		time.Sleep(time.Second * 30)
		io.WriteString(stdin, "exit\n")
	}

	cmd.Wait()
	if !cmd.ProcessState.Success() {
		return errors.New(cmd.ProcessState.String())
	}
	return nil
}

func (b *Curl) GetAttackConfig() AttackConfig {
	return b.AttackConfig
}

func (b *Curl) LoadAttackConfig(ac AttackConfig) error {
	if err := ac.Validate(); err != nil {
		return err
	}
	b.AttackConfig = ac
	return nil
}
