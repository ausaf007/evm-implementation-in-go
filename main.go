package main

import (
	"ethereum-vm/evm"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const DISCLAIMER = "\nDisclaimer(For Windows users only):" +
	"\nPlease note Windows CMD or Powershell might have a input limit of 255 characters." +
	"\nIf you are on windows, consider using WSL or Git Bash Emulator to run this program."

// loadFlags that specifies the verbosity level of the program based on user specified flag
func loadFlags() {
	isVerbose := flag.Bool("verbose", false, "Specifies verbosity of logs. True means Info Level. "+
		"False means Warn Level")
	flag.Parse()

	if *isVerbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

// printError prints the error and asks user to try again
func printError(err error) {
	fmt.Println("\nError Encountered:", err,
		"\nPlease try again.")
}

// Driver code
// This shows how to run the EVM by providing the list of 
// instructions as a bytecode
func main() {

	loadFlags()
	fmt.Println(DISCLAIMER)

	// This keeps iterating the loop until a valid bytecode is entered,
	// and the bytecode is successfully executed
	for true {

		// Assuming the string entered is in hex format, with no 0x prefix
		fmt.Println("\nEnter bytecode to be executed: ")

		// Taking input bytecode
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			printError(err)
			continue
		}

		hash, gasCost, err := evm.RunByteCode(input)
		if err != nil {
			printError(err)
			continue
		}

		fmt.Println("keccak256 Hash of memory =", hash)
		fmt.Println("Gas Cost =", gasCost)
		break
	}
	fmt.Println("Execution Successful!")

}
