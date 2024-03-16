package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	// Must include a backend implementation
	// See CommonLog for other options: https://github.com/tliron/commonlog
	_ "github.com/tliron/commonlog/simple"
)

const lsName = "dbt_lsp"

var (
	version  string = "0.0.1"
	handler  protocol.Handler
	manifest Manifest
	ROOT_DIR string
)

func main() {
	// This increases logging verbosity (optional)
	path := "/root/dev/go/dbt-language-server/out/log.txt"
	commonlog.Configure(1, &path)

	handler = protocol.Handler{
		Initialize:             initialize,
		Initialized:            initialized,
		Shutdown:               shutdown,
		SetTrace:               setTrace,
		TextDocumentDefinition: definitionHandler,
	}

	server := server.NewServer(&handler, lsName, false)
	server.Log.Info("Hello mark")

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	initLog := commonlog.GetLoggerf("%s.init", lsName)
	initLog.Infof("has got manifest %v", manifest)

	ROOT_DIR = params.WorkspaceFolders[0].URI

	rootDir := strings.ReplaceAll(ROOT_DIR, "file://", "")
	manifestName := fmt.Sprintf("%s/target/manifest.json", rootDir)
	file, err := os.ReadFile(manifestName)
	if err != nil {
		initLog.Errorf("could not open file %s %v", manifestName, err)
	} else {
		initLog.Info("opened manifest file")
	}

	err = json.Unmarshal(file, &manifest)
	if err != nil {
		initLog.Errorf("could not convert to json %v", err)
		return nil, err
	}

	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func highlighHandler(context *glsp.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	highlightLog := commonlog.GetLoggerf("%s.highlighter", lsName)
	text := protocol.DocumentHighlightKindText
	documentHighlights := []protocol.DocumentHighlight{
		{
			Range: protocol.Range{Start: protocol.Position{Line: 5, Character: 13}, End: protocol.Position{Line: 5, Character: 30}},
			Kind:  &text,
		},
	}
	highlightLog.Info("Highlighting shit")

	return documentHighlights, nil
}

func definitionHandler(context *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	definitionLog := commonlog.GetLoggerf("%s.definition", lsName)
	definitionLog.Infof("text document %v", params.TextDocument)

	fileContent, err := os.ReadFile(strings.ReplaceAll(params.TextDocument.URI, "file://", ""))
	if err != nil {
		definitionLog.Info("cannot read file")
		return nil, err
	}

	textDocumentFilePath := strings.Split(params.TextDocument.URI, "/")

	modelName := strings.ReplaceAll(textDocumentFilePath[len(textDocumentFilePath)-1], ".sql", "")
	definitionLog.Infof("modelName %s", modelName)

	key := fmt.Sprintf("model.%s.%s", manifest.Metadata.ProjectName, modelName)

	definitionLog.Infof("firstKey %s", key)
	val, ok := manifest.Nodes[key]
	if !ok {
		definitionLog.Infof("key not found %s", key)
	} else {
	}

	ok, reference := val.DoThing2(string(fileContent), params.Position)
	if !ok {
		definitionLog.Infof("reference could not be found %d %s", params.Position, reference)
		return nil, nil
	}

	key = fmt.Sprintf("model.%s.%s", manifest.Metadata.ProjectName, reference)
	definitionLog.Infof("going to find key %s", key)
	originalPath := manifest.Nodes[key].OriginalPath

	filePath := fmt.Sprintf("%s/%s", ROOT_DIR, originalPath)

	return protocol.Location{
		URI: filePath,
		Range: protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End:   protocol.Position{Line: 0, Character: 0},
		},
	}, nil
}
