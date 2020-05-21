package ticcmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	bin   string
	Debug bool
}

func (c Client) runCmd(args ...string) ([]byte, error) {
	cmd := exec.Command(c.bin, args...)
	op, err := cmd.Output()
	if c.Debug {
		msg := []string{"running", c.bin}
		msg = append(msg, args...)
		log.Println(msg)
		log.Println(string(op))
	}
	return op, err
}

func (c Client) Status() (Status, error) {
	op, err := c.runCmd("--status")
	if err != nil {
		return Status{}, fmt.Errorf("error retrieving status: %w", err)
	}
	return ParseStatus(op)
}

func (c Client) SetPosition(pos int) error {
	if err := c.exitSafeStart(); err != nil {
		return err
	}
	st, err := c.Status()
	if err != nil {
		return err
	}
	if !st.Energized {
		return errors.New("motor not energized")
	}
	op, err := c.runCmd("--position", strconv.Itoa(pos))
	if err != nil {
		return fmt.Errorf("error setting position: %s (%w)", string(op), err)
	}
	return nil
}

func (c Client) WaitForPosition(pos int, max time.Duration) error {
	end := time.Now().Add(max)
	for {
		if time.Now().After(end) {
			return errors.New("wait time expired")
		}
		st, err := c.Status()
		if err != nil {
			return err
		}
		fmt.Println("tgt=", pos, st.CurrentPosition, st.TargetPosition)
		if st.CurrentPosition == pos {
			return nil
		}
		time.Sleep(max / 10)
	}
}
func (c Client) exitSafeStart() error {
	// start with an enter safe start command
	op, err := c.runCmd("--exit-safe-start")
	if err != nil {
		return fmt.Errorf("error entering safe start: %s (%w)", string(op), err)
	}
	return nil
}

func (c Client) Energize() error {
	op, err := c.runCmd("--energize")
	if err != nil {
		return fmt.Errorf("error energizing: %s (%w)", string(op), err)
	}
	return nil
}

func (c Client) Deenergize() error {
	op, err := c.runCmd("--deenergize")
	if err != nil {
		return fmt.Errorf("error de-energizing: %s (%w)", string(op), err)
	}
	return nil
}

func NewClient(bin string) (*Client, error) {
	if len(bin) > 0 && strings.HasPrefix(bin, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error looking up home dir: %w", err)
		}
		bin = strings.Replace(bin, "~/", home+"/", 1)
	}
	// not passed in, so try to look it up
	if bin == "" {
		bin, _ = exec.LookPath("ticcmd")
	}
	c := &Client{
		bin: bin,
	}

	return c, nil
}
