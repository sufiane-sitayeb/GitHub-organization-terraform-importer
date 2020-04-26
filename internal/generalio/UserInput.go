package generalio

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// InputFromUser waits for an imput from the user before continuing the terraform input
func InputFromUser(resource string, autoApprove bool) bool {
	if !autoApprove {
		// Read input to make sure we want to continue and import everything
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\n\n*** Are you sure you want to continue with the import of the resource %s? [Y/N] ***\n", resource)
		input, _ := reader.ReadString('\n')
		// Normalizing the input, removing the \n that ReadString leaves in the string and setting everything lowercase
		inputNormal := strings.TrimSuffix(strings.ToLower(input), "\n")

		// Check input from user we only accept Y/YES or N/NO
		switch inputNormal {
		case "y", "yes":
			return true
		case "n", "no":
			fmt.Printf("\nWon't import, your terraform files are in the ./terraform/%s folder. Moving to next resource\n\n", resource)
		default:
			fmt.Printf("\nInvalid answer only Y/N, YES/NO are accepted your answer was: %s\n", input)
			os.Exit(1)
		}
		return false
	}
	return true
}

// EnvExist Check if an env var exists
func EnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}
