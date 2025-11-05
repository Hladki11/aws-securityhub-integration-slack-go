package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"

	"github.com/cruxstack/aws-securityhub-integration-slack-go/internal/app"
)

func main() {
	envpath := filepath.Join(".env")
	log.Print(envpath)
	if _, err := os.Stat(envpath); err == nil {
		_ = godotenv.Load(envpath)
	}

	cfg, err := app.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	a := app.New(cfg)

	path := filepath.Join("fixtures", "samples.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var findings []json.RawMessage
	if err := json.Unmarshal(raw, &findings); err != nil {
		log.Fatal(err)
	}

	for i, finding := range findings {
		detail := map[string]interface{}{
			"findings": []json.RawMessage{finding},
		}
		detailBytes, err := json.Marshal(detail)
		if err != nil {
			log.Fatalf("marshal detail: %v", err)
		}

		evt := awsevents.CloudWatchEvent{
			Version:    "0",
			ID:         fmt.Sprintf("sample-%d", i),
			DetailType: "Findings Imported V2",
			Source:     "aws.securityhub",
			AccountID:  "123456789012",
			Region:     "us-east-1",
			Detail:     detailBytes,
		}

		if err := a.Process(evt); err != nil {
			log.Fatalf("process sample %d: %v", i, err)
		}
		log.Printf("processed sample %d successfully", i)
	}
}
