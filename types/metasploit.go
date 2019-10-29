package types

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/attacker/utils"
)

const (
	resourceOutputDir = "/metasploit/"
)

func init() {
	RegisterAttacker("metasploit", &Metasploit{})
}

type Metasploit struct {
	AttackConfig
}

func (m *Metasploit) LaunchAttack(attackScript string) (err error) {
	var stdin io.WriteCloser

	if len(attackScript) == 0 {
		return errors.New("attack script is not specified")
	}

	// launch attack
	cmd := exec.Command("msfconsole", "-r", attackScript)

	if m.IsReverseShell {
		stdin, err = cmd.StdinPipe()
		if err != nil {
			return errors.New("errors in reverse shell")
		}
	}

	cmd.Start()

	if m.IsReverseShell {
		time.Sleep(time.Second * 30)
		io.WriteString(stdin, "exit\n")
	}

	cmd.Wait()
	if !cmd.ProcessState.Success() {
		return errors.New(cmd.ProcessState.String())
	}
	return nil
}

func (m *Metasploit) PrepareAttack() (string, error) {
	var stmt string

	if err := m.Validate(); err != nil {
		return "", err
	}

	// update resource file when reverse shell is specified
	if m.IsReverseShell {
		if len(m.LPORT) == 0 {
			localPort, err := utils.GetFreePort()
			if err != nil {
				return "", err
			}
			m.LPORT = localPort
		}

		if len(m.LHOST) == 0 {
			eth0, _ := net.InterfaceByName("eth0")
			addrs, _ := eth0.Addrs()
			localHostIP := strings.Split(addrs[0].String(), "/")[0]
			if localHostIP == "" {
				return "", errors.New("empty local IP address")
			}
			m.LHOST = localHostIP
		}
	}

	fileName := resourceOutputDir + m.Name + ".rc"
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	if len(m.LHOST) > 0 {
		stmt = fmt.Sprintf("set LHOST %s\n", m.LHOST)
		writer.WriteString(stmt)
	}

	if len(m.LPORT) > 0 {
		stmt = fmt.Sprintf("set LPORT %s\n", m.LPORT)
		writer.WriteString(stmt)
	}

	if len(m.RHOST) > 0 {
		stmt = fmt.Sprintf("set RHOST %s\n", m.RHOST)
		writer.WriteString(stmt)
	}

	if len(m.RPORT) > 0 {
		stmt = fmt.Sprintf("set RPORT %s\n", m.RPORT)
		writer.WriteString(stmt)
	}

	if len(m.Exploit) > 0 {
		stmt = fmt.Sprintf("use %s\n", m.Exploit)
		writer.WriteString(stmt)
	}

	if len(m.Payload) > 0 {
		stmt = fmt.Sprintf("set PAYLOAD %s\n", m.Payload)
		writer.WriteString(stmt)
	}

	if len(m.TargetURI) > 0 {
		stmt = fmt.Sprintf("set TARGETURI %s\n", m.TargetURI)
		writer.WriteString(stmt)
	}

	if len(m.Database) > 0 {
		stmt = fmt.Sprintf("set DATABASE %s\n", m.Database)
		writer.WriteString(stmt)
	}

	if len(m.Username) > 0 {
		stmt = fmt.Sprintf("set USERNAME %s\n", m.Username)
		writer.WriteString(stmt)
	}

	if len(m.Password) > 0 {
		stmt = fmt.Sprintf("set PASSWORD %s\n", m.Password)
		writer.WriteString(stmt)
	}

	if len(m.SRVHOST) > 0 {
		stmt = fmt.Sprintf("set SRVHOST %s\n", m.SRVHOST)
		writer.WriteString(stmt)
	}

	// add run and exit to the end of the resource file
	writer.WriteString("\nrun\nexit\n")
	writer.Flush()

	return fileName, nil
}

func (m *Metasploit) GetAttackConfig() AttackConfig {
	return m.AttackConfig
}

func (m *Metasploit) LoadAttackConfig(ac AttackConfig) error {
	if err := ac.Validate(); err != nil {
		return err
	}

	m.AttackConfig = ac
	return nil
}
