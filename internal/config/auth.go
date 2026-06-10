package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Path returns the path to the skpm config file.
func Path() (string, error) {
	return configPath()
}

// SetRegistryToken reads the config file, sets the token for the named registry,
// and writes the file back. If the registry does not exist, it is created with a
// minimal "local" placeholder entry so the user can fill in the remaining fields.
func SetRegistryToken(registryName, token string) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	return editConfigYAML(path, func(node *yaml.Node) error {
		mapping := yamlDocMapping(node)
		if mapping == nil {
			return fmt.Errorf("config file is not a YAML mapping")
		}
		regsNode := yamlGetOrCreate(mapping, "registries", yaml.MappingNode)
		regNode := yamlGetOrCreate(regsNode, registryName, yaml.MappingNode)
		yamlSetScalar(regNode, "token", token)
		return nil
	})
}

// ClearRegistryToken removes the token (and auth.token) for the named registry.
// Returns nil even when the registry or token does not exist.
func ClearRegistryToken(registryName string) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	return editConfigYAML(path, func(node *yaml.Node) error {
		mapping := yamlDocMapping(node)
		if mapping == nil {
			return nil
		}
		regsNode := yamlFind(mapping, "registries")
		if regsNode == nil {
			return nil
		}
		regNode := yamlFind(regsNode, registryName)
		if regNode == nil {
			return nil
		}
		yamlDeleteKey(regNode, "token")
		authNode := yamlFind(regNode, "auth")
		if authNode != nil {
			yamlDeleteKey(authNode, "token")
		}
		return nil
	})
}

// editConfigYAML reads the config at path as a yaml.Node, calls fn to mutate
// it, then writes the result back atomically.
func editConfigYAML(path string, fn func(*yaml.Node) error) error {
	var doc yaml.Node

	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read config %s: %w", path, err)
	}

	if len(data) == 0 {
		doc = yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{
			{Kind: yaml.MappingNode, Tag: "!!map"},
		}}
	} else {
		if err := yaml.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("parse config %s: %w", path, err)
		}
	}

	if err := fn(&doc); err != nil {
		return err
	}

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	return os.WriteFile(path, out, 0o600)
}

// yamlDocMapping returns the root mapping node from a document node.
func yamlDocMapping(doc *yaml.Node) *yaml.Node {
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		n := doc.Content[0]
		if n.Kind == yaml.MappingNode {
			return n
		}
	}
	return nil
}

// yamlFind returns the value node for key in a mapping node, or nil.
func yamlFind(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}

// yamlGetOrCreate returns the value node for key in a mapping node, creating
// a new node of the given kind if the key is absent.
func yamlGetOrCreate(mapping *yaml.Node, key string, kind yaml.Kind) *yaml.Node {
	if n := yamlFind(mapping, key); n != nil {
		return n
	}
	tagFor := map[yaml.Kind]string{
		yaml.MappingNode:  "!!map",
		yaml.SequenceNode: "!!seq",
		yaml.ScalarNode:   "!!str",
	}
	keyNode := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}
	valNode := &yaml.Node{Kind: kind, Tag: tagFor[kind]}
	mapping.Content = append(mapping.Content, keyNode, valNode)
	return valNode
}

// yamlSetScalar sets (or creates) a scalar key in a mapping node.
func yamlSetScalar(mapping *yaml.Node, key, value string) {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content[i+1].Kind = yaml.ScalarNode
			mapping.Content[i+1].Tag = "!!str"
			mapping.Content[i+1].Value = value
			return
		}
	}
	mapping.Content = append(mapping.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value},
	)
}

// yamlDeleteKey removes a key-value pair from a mapping node.
func yamlDeleteKey(mapping *yaml.Node, key string) {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content = append(mapping.Content[:i], mapping.Content[i+2:]...)
			return
		}
	}
}
