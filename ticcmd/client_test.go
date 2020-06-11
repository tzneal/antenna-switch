package ticcmd_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tzneal/antenna-switch/ticcmd"
)

func TestNewClient(t *testing.T) {
	c, err := ticcmd.NewClient("~/.bin/ticcmd")
	if err != nil {
		t.Fatalf("error constructing client: %s", err)
	}
	c.Debug = false
	st, err := c.Status()
	if err != nil {
		t.Fatalf("expected no error with status: %s", err)
	}
	fmt.Println(st)

	for i := 0; i < 10; i++ {
		c.Energize()
		err = c.SetPosition(50)
		if err != nil {
			t.Fatalf("error setting position: %s", err)
		}
		err = c.WaitForPosition(50, 1*time.Second)
		if err != nil {
			t.Fatalf("error wait position: %s", err)
		}
		c.Deenergize()
		time.Sleep(250 * time.Millisecond)
		c.Energize()
		err = c.SetPosition(0)
		if err != nil {
			t.Fatalf("error setting position: %s", err)
		}
		err = c.WaitForPosition(0, 1*time.Second)
		if err != nil {
			t.Fatalf("error wait position: %s", err)
		}
		c.Deenergize()
		time.Sleep(250 * time.Millisecond)

	}

}

func TestParseStatus(t *testing.T) {
	st := []byte(`Label:                         Tic T834 Stepper Motor Controller
Serial number:                00305428
Firmware version:             1.06
Last reset:                   Power-on reset
Up time:                      0:38:02

Encoder position:             0
Input state:                  Position
Input after averaging:        N/A
Input after hysteresis:       N/A
Input before scaling:         N/A
Input after scaling:          -44
Forward limit active:         No
Reverse limit active:         No

VIN voltage:                  5.3 V
Operation state:              De-energized
Energized:                    No
Homing active:                No

Target:                       No target
Current position:             -91
Position uncertain:           Yes
Current velocity:             0

Errors currently stopping the motor:
  - Intentionally de-energized
  - Command timeout
  - Safe start violation
Errors that occurred since last check:
  - Command timeout
  - Safe start violation
`)
	s, err := ticcmd.ParseStatus(st)
	if err != nil {
		t.Errorf("error parsing status: %s", err)
	}
	fmt.Println(s)
}
