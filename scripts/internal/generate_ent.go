package internal

import (
	"log"
	"os"
	"os/exec"
)

// GenerateEnt runs the Ent code generator.
func GenerateEnt() error {
	cmd := exec.Command("go", "run", "-mod=mod", "entgo.io/ent/cmd/ent", "generate", "./ent/schema")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Running Ent code generation...")
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Println("Ent code generation completed successfully.")
	return nil
}
