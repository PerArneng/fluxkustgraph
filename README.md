# FluxKustGraph

FluxKustGraph is a CLI tool designed to analyze Kustomization YAML files within a
specified directory, generating a Mermaid class diagram to visualize the dependencies
and relationships between the defined resources. This utility is invaluable for
developers and operators utilizing Flux CD and Kustomize in Kubernetes, offering a
clear, visual understanding of configuration interdependencies.

## Features

- **Recursive Search**: Traverses a directory to locate Kustomization YAML files.
- **YAML Parsing**: Extracts and interprets the relevant parts of Kustomization files.
- **Diagram Generation**: Creates a Mermaid class diagram to represent resource relationships.
- **Flexible Output**: Saves the generated diagram to a user-defined file.

## Prerequisites

To use FluxKustGraph, ensure you have Go version 1.15 or later installed on your system.

## Installation

First, clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/fluxkustgraph.git
cd fluxkustgraph
```
Then, compile the source code:

```bash
go build -o fluxkustgraph
```

## Usage
Execute FluxKustGraph by specifying the source directory containing your Kustomization
YAML files and the output file path for the diagram:

```bash
./fluxkustgraph -source /path/to/yaml/files -output /path/to/output/diagram.md
```

This command will scan the specified directory for YAML files, parse them to identify
kustomization dependencies, and generate a Mermaid diagram that is saved to the
specified output file.