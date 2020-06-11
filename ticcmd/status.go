package ticcmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Status struct {
	Name                string
	Serial              string
	FirmwareVersion     string
	LastReset           string
	Uptime              string
	EncoderPosition     int
	InputState          string
	InputAfterAveraging string
	InputAfterHystersis string
	InputBeforeScaling  string
	InputAafterScaling  string
	ForwardLimitActive  bool
	ReverseLimitActive  bool
	VINVoltage          string
	OperationState      string
	Energized           bool
	HomingActive        bool
	Target              string
	CurrentPosition     int
	PositinUncertain    bool
	CurrentVelocity     int

	ErrorsStoppingMotor  []string
	ErrorsSinceLastCheck []string
	TargetPosition       int
	TargetVelocity       int
}

func ParseStatus(op []byte) (Status, error) {
	r := bytes.NewReader(op)
	sc := bufio.NewScanner(r)
	st := Status{}

	var err error
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if len(line) == 0 {
			continue
		}
		colonIdx := strings.IndexByte(line, ':')
		if colonIdx != -1 {
			fieldName := strings.TrimSpace(line[0:colonIdx])
			fieldValue := strings.TrimSpace(line[colonIdx+1:])
			switch fieldName {
			case "Label", "Name":
				st.Name = fieldValue
			case "Serial number":
				st.Serial = fieldValue
			case "Firmware version":
				st.FirmwareVersion = fieldValue
			case "Last reset":
				st.LastReset = fieldValue
			case "Up time":
				st.Uptime = fieldValue
			case "Encoder position":
				encPos, err := strconv.ParseInt(fieldValue, 10, 64)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing encoder position %s: %w", fieldValue, err)
				}
				st.EncoderPosition = int(encPos)
			case "Input state":
				st.InputState = fieldValue
			case "Input after averaging":
				st.InputAfterAveraging = fieldValue
			case "Input after hysteresis":
				st.InputAfterHystersis = fieldValue
			case "Input before scaling":
				st.InputBeforeScaling = fieldValue
			case "Input after scaling":
				st.InputAafterScaling = fieldValue
			case "Forward limit active":
				st.ForwardLimitActive, err = parseBool(fieldValue)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing forward limit active %s: %w", fieldValue, err)
				}
			case "Reverse limit active":
				st.ReverseLimitActive, err = parseBool(fieldValue)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing reverse limit active %s: %w", fieldValue, err)
				}
			case "VIN voltage":
				st.VINVoltage = fieldValue
			case "Operation state":
				st.OperationState = fieldValue
			case "Energized":
				st.Energized, err = parseBool(fieldValue)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing energized %s: %w", fieldValue, err)
				}
			case "Homing active":
				st.HomingActive, err = parseBool(fieldValue)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing homing active %s: %w", fieldValue, err)
				}
			case "Target":
				st.Target = fieldValue
			case "Current position":
				pos, err := strconv.ParseInt(fieldValue, 10, 64)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing current position %s: %w", fieldValue, err)
				}
				st.CurrentPosition = int(pos)
			case "Target position":
				pos, err := strconv.ParseInt(fieldValue, 10, 64)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing target position %s: %w", fieldValue, err)
				}
				st.TargetPosition = int(pos)
			case "Target velocity":
				vel, err := strconv.ParseInt(fieldValue, 10, 64)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing target velocity %s: %w", fieldValue, err)
				}
				st.TargetVelocity = int(vel)

			case "Position uncertain":
				st.PositinUncertain, err = parseBool(fieldValue)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing position uncertain %s: %w", fieldValue, err)
				}
			case "Current velocity":
				vel, err := strconv.ParseInt(fieldValue, 10, 64)
				if err != nil {
					return Status{}, fmt.Errorf("error parsing current velocity %s: %w", fieldValue, err)
				}
				st.CurrentVelocity = int(vel)
			case "Errors currently stopping the motor":
				stoppingMotorErrors := true
				for sc.Scan() {
					line := strings.TrimSpace(sc.Text())
					if strings.HasPrefix(line, "Errors that occurred since last check:") {
						stoppingMotorErrors = false
						continue
					}
					if strings.HasPrefix(line, "Last motor driver errors:") {
						stoppingMotorErrors = false
						continue
					}
					msg := line
					if len(msg) == 0 {
						continue
					}
					if msg[0] != '-' {
						return Status{}, fmt.Errorf("unexpected error format: %s", line)
					}
					msg = msg[2:]
					if stoppingMotorErrors {
						st.ErrorsStoppingMotor = append(st.ErrorsSinceLastCheck, msg)
					} else {
						st.ErrorsSinceLastCheck = append(st.ErrorsSinceLastCheck, msg)
					}
				}

			default:
				log.Println("unsupported field", fieldName, "=", fieldValue)
			}
		} else {
			log.Println("unexpected line in status:", line)
		}
	}
	return st, nil
}

func parseBool(value string) (bool, error) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	switch value {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	}
	return false, fmt.Errorf("error parsing '%s' as boolean", value)
}
