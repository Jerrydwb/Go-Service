package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dryRun := flag.Bool("dry", false, "Realiza un simulacro sin borrar archivos")
	flag.Parse()

	basePath := "/var/www/v2/suite/storage/storage"
	foldersToDelete := []string{"agente", "prod", "demo"}

	if *dryRun {
		fmt.Println("⚠️  MODO PRUEBA (DRY RUN) ACTIVADO - No se borrará nada")
		fmt.Println("-------------------------------------------------------")
	}

	entries, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for _, entry := range entries {
		name := entry.Name()

		if entry.IsDir() && strings.HasPrefix(name, "I-") && name != "I-3" && name != "I-undefined" {
			for _, folder := range foldersToDelete {
				target := filepath.Join(basePath, name, "fe", folder)

				if _, err := os.Stat(target); err == nil {
					if *dryRun {
						fmt.Printf("[SIMULACRO] Se borraría: %s\n", target)
					} else {
						err := os.RemoveAll(target)
						if err != nil {
							fmt.Printf("[!] Error eliminando %s: %v\n", target, err)
						} else {
							fmt.Printf("[OK] Eliminado: %s\n", target)
						}
					}
				}
			}
		}
	}

	if *dryRun {
		fmt.Println("-------------------------------------------------------")
		fmt.Println("Simulacro finalizado. Para borrar realmente, quita el flag -dry")
	} else {
		fmt.Println("Limpieza real completada.")
	}
}
