package types

import (
	"errors"
	"strings"
)

var (
	Attackers = map[string]Attacker{}
)

type Attacker interface {
	// LoadAttackConfig loads a attack config
	LoadAttackConfig(ac AttackConfig) error
	// GetAttackConfig returns the attacker's configuration
	GetAttackConfig() AttackConfig
	// PrepareAttack generate an attack script (rc file or bash script)
	PrepareAttack() (attackScript string, err error)
	// LaunchAttack execute the attack script
	LaunchAttack(attackScript string) error
}

// LaunchNewAttack launches a new attack
func LaunchNewAttack(name string, ac AttackConfig) (err error) {
	var attacker Attacker
	var exists bool
	var attackScript string
	if len(name) == 0 {
		name = "metasploit"
	}

	name = strings.ToLower(name)

	if attacker, exists = Attackers[name]; !exists {
		return errors.New("unknown attack tool")
	}

	if err := attacker.LoadAttackConfig(ac); err != nil {
		return err
	}

	if attackScript, err = attacker.PrepareAttack(); err != nil {
		return err
	}

	if err = attacker.LaunchAttack(attackScript); err != nil {
		return err
	}

	return nil
}

// RegisterAttacker register a attacker tool
func RegisterAttacker(name string, a Attacker) {
	Attackers[name] = a
}
