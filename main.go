package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UpsStatus struct {
	Model string
	Firmware string
	VoltageRating int
	PowerRating int
	State string
	UtilityVoltage int
	OutputVoltage int
	BatteryCapacity int
	RemainingRuntime time.Duration
	Load int
	LineInteraction string
	TestResult string
	LastPowerEvent string
}

func main() {
	// create the pwrstat status command
	var cmd = exec.Command("pwrstat", "-status")

	// execute the command
	var output, err = cmd.Output()

	// check for failure
	if err != nil {
		// parse command output into a status struct
		var status = parsePwrstatOutput(string(output))
		fmt.Printf("%+v\n", status)
	}

	// expose via snmp
}

func parsePwrstatOutput(output string) UpsStatus {
	var status UpsStatus

	var statusRegex = regexp.MustCompile("\\s+([^\\.]+)\\.+ (.+)")

	/* example output:
	The UPS information shows as following:

	Properties:
		Model Name................... CP1500PFCLCDa
		Firmware Number.............. CXXJV2020538
		Rating Voltage............... 120 V
		Rating Power................. 1000 Watt(1500 VA)

	Current UPS status:
		State........................ Normal
		Power Supply by.............. Utility Power
		Utility Voltage.............. 119 V
		Output Voltage............... 119 V
		Battery Capacity............. 100 %
		Remaining Runtime............ 98 min.
		Load......................... 50 Watt(5 %)
		Line Interaction............. None
		Test Result.................. In progress
		Last Power Event............. Blackout at 2019/12/15 16:31:55
	*/

	var splitLine = strings.Split(output, "\n")

	var statusMap = make(map[string]string)

	for _, x := range splitLine {
		matches := statusRegex.FindStringSubmatch(x)

		if len(matches) > 0 {
			// create a map of status type : value
			statusMap[matches[1]] = matches[2]
		}
	}

	status = UpsStatus{
		Model: statusMap["Model Name"],
		Firmware: statusMap["Firmware Number"],
		VoltageRating: parseNumber(statusMap["Rating Voltage"]),
		PowerRating: parseNumber(statusMap["Rating Power"]),
		State: statusMap["State"],
		UtilityVoltage: parseNumber(statusMap["Utility Voltage"]),
		OutputVoltage: parseNumber(statusMap["Output Voltage"]),
		BatteryCapacity: parseNumber(statusMap["Battery Capacity"]),
		RemainingRuntime: time.Duration(parseNumber(statusMap["Remaining Runtime"])) * time.Minute,
		Load: parseNumber(statusMap["Load"]),
		LineInteraction: statusMap["Line Interaction"],
		TestResult: statusMap["Test Result"],
		LastPowerEvent: statusMap["Last Power Event"],
	}

	return status
}

func parseNumber(output string) int {
	var num int

	var numberRegex = regexp.MustCompile("(\\d+).+")

	var matches = numberRegex.FindStringSubmatch(output)

	if len(matches) > 0 {
		num, _ = strconv.Atoi(matches[1])
	}

	return num
}
