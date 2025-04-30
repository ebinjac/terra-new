package terraform

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// ExtractTar extracts a tar file to the specified destination
func ExtractTar(tarPath, destPath string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("error opening tar file: %v", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %v", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar: %v", err)
		}

		target := filepath.Join(destPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("error creating directory: %v", err)
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("error creating file: %v", err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("error writing file: %v", err)
			}
			f.Close()
		}
	}

	return nil
}

// RunTerraformInit runs 'terraform init' in the specified directory
func RunTerraformInit(dir string) error {
	cmd := exec.Command("terraform", "init")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running terraform init: %v", err)
	}

	return nil
}

// GenerateGraph runs 'terraform graph' and saves the output to graph.gv
func GenerateGraph(dir string) error {
	outputFile := filepath.Join(dir, "graph.gv")
	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating graph file: %v", err)
	}
	defer out.Close()

	cmd := exec.Command("terraform", "graph")
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error generating terraform graph: %v", err)
	}

	return nil
}
