package main

import (
	"fmt"
	"path/filepath"

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
	settings ProjectSettings
	ROOT_DIR string
)

func main() {
	// This increases logging verbosity (optional)
	path := "/root/dev/go/dbt-language-server/out/log.txt"
	commonlog.Configure(1, &path)

	handler = protocol.Handler{
		Initialize:                     initialize,
		Initialized:                    initialized,
		Shutdown:                       shutdown,
		SetTrace:                       setTrace,
		TextDocumentDefinition:         definitionHandler,
		TextDocumentHover:              hoverHandler,
		WorkspaceDidChangeWatchedFiles: fileChanged,
	}

	server := server.NewServer(&handler, lsName, false)
	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	initLog := commonlog.GetLoggerf("%s.init", lsName)

	ROOT_DIR = params.WorkspaceFolders[0].URI
	settings, err := LoadSettings(ROOT_DIR)
	if err != nil {
		initLog.Errorf("ERROR %v", err)
		return nil, err
	}

	schemas, err := settings.GetSchemaFiles()
	if err != nil {
		initLog.Errorf("Could not load schema files %v", err)
	} else {
		initLog.Infof("Loaded Schema file %v", schemas)
	}

	manifest, err = settings.LoadManifestFile()
	if err != nil {
		initLog.Errorf("could not load manifest file %v", err)
		return nil, err
	}

	capabilities := handler.CreateServerCapabilities()
	initLog.Info("Returning initialized")
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

func highlighHandler(_ *glsp.Context, _ *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
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

	file := getModelNameFromFilePath(params.TextDocument.URI)
	key := fmt.Sprintf("model.%s.%s", manifest.Metadata.ProjectName, file)

	val, ok := manifest.Nodes[key]
	if !ok {
		definitionLog.Infof("could not find initial key %v", key)
		return nil, nil
	}

	model, err := val.GetDefinition(DefinitionRequest{
		FileUri:     params.TextDocument.URI,
		Position:    params.Position,
		Manifest:    manifest,
		ProjectName: manifest.Metadata.ProjectName,
	})
	if err != nil {
		definitionLog.Infof("getting the definition failed %v", err)
		return nil, err
	}

	return protocol.Location{
		URI: filepath.Join(ROOT_DIR, model.FileName),
		Range: protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End:   protocol.Position{Line: 0, Character: 0},
		},
	}, nil
}

func hoverHandler(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	definitionLog := commonlog.GetLoggerf("%s.hover", lsName)

	file := getModelNameFromFilePath(params.TextDocument.URI)
	key := fmt.Sprintf("model.%s.%s", manifest.Metadata.ProjectName, file)

	val, ok := manifest.Nodes[key]
	if !ok {
		definitionLog.Infof("could not find initial key %v", key)
		return nil, nil
	}

	model, err := val.GetDefinition(DefinitionRequest{
		FileUri:  params.TextDocument.URI,
		Position: params.Position,
	})
	if err != nil {
		definitionLog.Infof("getting the definition failed %v", err)
		return nil, err
	}

	key = fmt.Sprintf("model.%s.%s", manifest.Metadata.ProjectName, model)

	referencedNode, ok := manifest.Nodes[key]
	if !ok {
		definitionLog.Infof("could not referenced key %v", key)
		return nil, nil
	}

	return &protocol.Hover{Contents: referencedNode.Description}, nil
}

func fileChanged(context *glsp.Context, params *protocol.DidChangeWatchedFilesParams) error {
	// for _, uri := range params.Changes {
	// 	file := strings.ReplaceAll(uri.URI, "file://", "")
	//
	//
	// }
	return nil
}
