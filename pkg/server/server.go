package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/terraform-cli/pkg/client"
	"github.com/terraform-cli/pkg/terraform"
)

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) handleTerraformProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request: %v", err), http.StatusBadRequest)
		return
	}

	// Parse input JSON
	var input client.InputConfig
	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Create temporary directory
	tempDir, err := ioutil.TempDir("", "terraform-*")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating temp directory: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	// Create HTTP client
	c := client.NewClient(input.Token)

	// Download and extract config
	configPath := filepath.Join(tempDir, "config.tar.gz")
	if err := c.DownloadFile(input.ConfigURL, configPath); err != nil {
		http.Error(w, fmt.Sprintf("Error downloading config: %v", err), http.StatusInternalServerError)
		return
	}

	extractPath := filepath.Join(tempDir, "config")
	if err := os.MkdirAll(extractPath, 0755); err != nil {
		http.Error(w, fmt.Sprintf("Error creating extract directory: %v", err), http.StatusInternalServerError)
		return
	}

	if err := terraform.ExtractTar(configPath, extractPath); err != nil {
		http.Error(w, fmt.Sprintf("Error extracting tar: %v", err), http.StatusInternalServerError)
		return
	}

	// Run terraform init
	if err := terraform.RunTerraformInit(extractPath); err != nil {
		http.Error(w, fmt.Sprintf("Error running terraform init: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate terraform graph
	if err := terraform.GenerateGraph(extractPath); err != nil {
		http.Error(w, fmt.Sprintf("Error generating terraform graph: %v", err), http.StatusInternalServerError)
		return
	}

	// Download plan.json
	planPath := filepath.Join(extractPath, "plan.json")
	if err := c.DownloadFile(input.PlanURL, planPath); err != nil {
		http.Error(w, fmt.Sprintf("Error downloading plan: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare upload request
	uploadReq := &client.UploadRequest{
		TFPlan:      planPath,
		ProductID:   input.ProductID,
		Name:        input.Name,
		MappingFile: "iriusrisk.yaml",
		TFGraphFile: filepath.Join(extractPath, "graph.gv"),
	}

	// Make final POST request
	endpoint := fmt.Sprintf("https://bayatest-dev.aexp.com/api/v1/tfplan/%s", input.ProductID)
	if err := c.UploadFiles(endpoint, uploadReq); err != nil {
		http.Error(w, fmt.Sprintf("Error uploading files: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Process completed successfully",
	})
}

func (s *Server) Start() error {
	http.HandleFunc("/api/terraform/process", s.handleTerraformProcess)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
}
