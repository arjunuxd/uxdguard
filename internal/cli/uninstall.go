package root

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arjunuxd/uxdguard/internal/app"
)

func cmdUninstall(args []string) int {
	force := false
	for i := 0; i < len(args); i++ {
		if args[i] == "-y" || args[i] == "--yes" {
			force = true
			continue
		}
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Printf("Usage: uxdguard uninstall [-y]\n\nRemoves the globally installed uxdguard executable.\n")
			return app.ExitOK
		}
		fmt.Fprintf(os.Stderr, "unknown flag: %s\n", args[i])
		return app.ExitUsageOrConfig
	}

	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return app.ExitOperationalFail
	}
	exe = filepath.Clean(exe)

	fmt.Printf("Current executable: %s\n\n", exe)
	fmt.Println("Note: Go's 'go install' does not have a built-in registry. This command will attempt to remove the executable if safe to do so.")
	fmt.Println("This will NOT remove source repos, configs, hooks, or other files.")

	if !force {
		fmt.Print("Remove this executable? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		resp, _ := reader.ReadString('\n')
		resp = strings.TrimSpace(strings.ToLower(resp))
		if resp != "y" && resp != "yes" {
			fmt.Println("Cancelled.")
			return app.ExitOK
		}
	}

	if err := os.Remove(exe); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to remove: %v\n\n", err)
		fmt.Println("On Windows, you may need to close all running uxdguard processes first.")
		fmt.Println("Manual removal:")
		fmt.Printf("  Remove: %s\n", exe)
		return app.ExitOperationalFail
	}

	fmt.Println("Removed successfully.")
	fmt.Println("Note: If uxdguard is still accessible, check other copies on PATH.")
	return app.ExitOK
}
