package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	args := os.Args[1:]

	cmdArgs := []string{
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("tern", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Command execution failed: ", err)
		fmt.Println("Output: ", string(output))
		panic(err)
	}

	fmt.Println("Command executed successfully: ", string(output))
}
