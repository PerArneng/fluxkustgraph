package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Kustomization represents the YAML structure we are interested in.
type Kustomization struct {
	ApiVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Interval           string `yaml:"interval"`
		ServiceAccountName string `yaml:"serviceAccountName"`
		DependsOn          []struct {
			Name string `yaml:"name"`
		} `yaml:"dependsOn"`
		SourceRef struct {
			Kind string `yaml:"kind"`
			Name string `yaml:"name"`
		} `yaml:"sourceRef"`
		Path  string `yaml:"path"`
		Prune bool   `yaml:"prune"`
	} `yaml:"spec"`
}

func main() {
	var sourceDir, outputFile string
	flag.StringVar(&sourceDir, "source", "", "Source directory to search for YAML files")
	flag.StringVar(&outputFile, "output", "", "Output file path for the Mermaid diagram")
	flag.Parse()

	if sourceDir == "" || outputFile == "" {
		fmt.Println("Source directory and output file path are required.")
		os.Exit(1)
	}

	// Find and parse YAML files.
	kustomizations, err := findAndParseYAMLs(sourceDir)
	if err != nil {
		fmt.Printf("Error processing YAML files: %v\n", err)
		os.Exit(1)
	}

	// Generate Mermaid class diagram.
	diagram := generateMermaidDiagram(kustomizations)

	// Write the diagram to the specified output file.
	if err := ioutil.WriteFile(outputFile, []byte(diagram), 0644); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Mermaid class diagram generated successfully.")
}

// findAndParseYAMLs recursively finds and parses YAML files in the given directory.
func findAndParseYAMLs(dir string) (map[string]Kustomization, error) {
	kustomizations := make(map[string]Kustomization)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			docs := strings.Split(string(content), "---")
			for _, doc := range docs {
				var k Kustomization
				if yaml.Unmarshal([]byte(doc), &k) == nil && k.Metadata.Name != "" {
					if strings.HasPrefix(k.ApiVersion, "kustomize.toolkit.fluxcd.io") {
						kustomizations[k.Metadata.Name] = k
					}
				}
			}
		}
		return nil
	})
	return kustomizations, err
}

// sanitizeName prepares a string to be used as an ID in Mermaid diagrams
// by replacing non-alphanumeric characters with underscores.
func sanitizeName(name string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, name)
}

// generateMermaidDiagram generates the Mermaid class diagram from the kustomizations.
func generateMermaidDiagram(kustomizations map[string]Kustomization) string {
	var diagram strings.Builder
	diagram.WriteString("classDiagram\n")
	for _, k := range kustomizations {
		safeName := sanitizeName(k.Metadata.Name)
		diagram.WriteString(fmt.Sprintf("    class %s {\n", safeName))
		diagram.WriteString(fmt.Sprintf("        +string name %s\n", k.Metadata.Name))
		diagram.WriteString(fmt.Sprintf("        +string namespace %s\n", k.Metadata.Namespace))
		diagram.WriteString("    }\n")
		//if k.Spec != nil {
		for _, dep := range k.Spec.DependsOn {
			safeDepName := sanitizeName(dep.Name)
			diagram.WriteString(fmt.Sprintf("    %s --> %s : depends on\n", safeName, safeDepName))
		}
		//}
	}
	return diagram.String()
}
