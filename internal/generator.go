package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"protocoll/internal/domain/collection"
	"protocoll/internal/domain/generator"
	pParser "protocoll/internal/domain/parser"
)

type Generator struct {
}

func (g *Generator) Generate(folder string, name string) error {
	items := make([]*collection.Node, 0)
	rootNode := &collection.Node{
		Item: &items,
	}

	refs := make(map[string]*collection.Node)
	refs[""] = rootNode

	variables := make([]string, 0)

	// todo: move this to a separate entity
	walk := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() && strings.HasSuffix(path, ".proto") {
			parseResult, err := g.parseFile(path)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Could not parse file: %s\n", path)
				return err
			}

			if len(parseResult.Package.Services) > 0 {
				path := make([]string, 0)

				for index, pathElement := range parseResult.Package.Path {
					path = append(path, pathElement)
					parentPath := path[:len(path)-1]
					pathKey := strings.Join(path, ".")
					parentPathKey := strings.Join(parentPath, ".")

					// if element was not created before
					if _, ok := refs[pathKey]; !ok {
						items := make([]*collection.Node, 0)
						newNode := collection.Node{
							Name: pathElement,
							Item: &items,
						}

						if index == len(parseResult.Package.Path)-1 {
							// adding a service here
							for _, service := range parseResult.Package.Services {
								itemsDeref := *newNode.Item
								itemsDeref = append(itemsDeref, g.createServiceNode(service))
								newNode.Item = &itemsDeref
							}
						}

						toNode := refs[parentPathKey]
						toNodeItemsDeref := *toNode.Item
						toNodeItemsDeref = append(toNodeItemsDeref, &newNode)
						toNode.Item = &toNodeItemsDeref

						refs[pathKey] = &newNode
					}
				}

				variables = append(variables, parseResult.Variables...)
			}
		}

		return nil
	}

	err := filepath.Walk(folder, walk)
	if err != nil {
		return err
	}

	coll := &collection.Collection{
		Info: &collection.Schema{
			Name:   name,
			Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item:     rootNode.Item,
		Variable: g.createVariables(variables),
		Auth:     g.createAuth(),
	}

	jsonData, _ := json.MarshalIndent(coll, "", "  ")
	fmt.Printf("%s\n", jsonData)

	return nil
}

func (g *Generator) createServiceNode(service *pParser.Service) (node *collection.Node) {
	items := make([]*collection.Node, 0)
	for _, endpoint := range service.Endpoints {
		endpointNode := collection.Node{
			Name: endpoint.Name,
			Request: &collection.Request{
				Method: endpoint.Method,
				Url: collection.URL{
					Raw:  endpoint.URL,
					Host: []string{"{{host}}"},
					Path: endpoint.Path,
				},
			},
		}

		if endpoint.Body != "" {
			endpointNode.Request.Body = &collection.Body{
				Mode: "raw",
				Raw:  endpoint.Body,
				Options: collection.Options{
					Raw: collection.Raw{
						Language: "json",
					},
				},
			}
		}

		items = append(items, &endpointNode)
	}

	return &collection.Node{
		Name: service.Name,
		Item: &items,
	}
}

func (g *Generator) createVariables(parsedVariables []string) (variables *[]collection.Variable) {
	result := make([]collection.Variable, 2)
	result[0] = collection.Variable{
		Key:   "host",
		Value: "",
	}
	result[1] = collection.Variable{
		Key:   "token",
		Value: "",
	}

	seenVariables := make(map[string]bool)

	for _, variable := range parsedVariables {
		if _, ok := seenVariables[variable]; ok {
			continue
		}

		result = append(result, collection.Variable{
			Key:   variable,
			Value: "",
		})
		seenVariables[variable] = true
	}

	return &result
}

func (g *Generator) createAuth() (auth *collection.Auth) {
	result := collection.Auth{
		Type: "bearer",
		Bearer: &collection.Bearer{
			Key:   "token",
			Value: "{{token}}",
			Type:  "string",
		},
	}

	return &result
}

// todo: move these functions to a separate entity

func (g *Generator) parseFile(file string) (*generator.Result, error) {
	reader, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s, err %v\n", file, err)
	}
	defer reader.Close()

	got, err := protoparser.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse, err %v\n", err)
	}

	packageName := got.ProtoBody[0].(*parser.Package).Name
	packageNameParts := strings.Split(packageName, ".")

	services, variables := g.traverse(got.ProtoBody, packageName)

	return &generator.Result{
		Package: &pParser.Package{
			Path:     packageNameParts,
			Key:      strings.Join(packageNameParts, ""),
			Services: services,
		},
		Variables: variables,
	}, nil
}

func (g *Generator) traverse(nodes []parser.Visitee, packageName string) (services []*pParser.Service, variables []string) {
	services = make([]*pParser.Service, 0)
	variables = make([]string, 0)

	for _, node := range nodes {
		var iNode interface{} = node
		switch v := iNode.(type) {
		case *parser.Service:
			newService, newVariables := g.traverseServiceNode(v, packageName)
			services = append(services, newService)
			variables = append(variables, newVariables...)
		default:
		}
	}

	return services, variables
}

func (g *Generator) traverseServiceNode(node *parser.Service, packageName string) (service *pParser.Service, variables []string) {
	variables = make([]string, 0)

	serviceName := node.ServiceName

	endpoints := make([]*pParser.Endpoint, 0)
	restEndpoints := make([]*pParser.Endpoint, 0)

	for _, bodyNode := range node.ServiceBody {
		RPCNode := bodyNode.(*parser.RPC)
		methodName := RPCNode.RPCName
		methodVerb := "POST"
		url := fmt.Sprintf("/%s.%s/%s", packageName, serviceName, methodName)
		path := strings.Split(url, "/")

		endpoints = append(endpoints, &pParser.Endpoint{
			Name:   methodName,
			Method: methodVerb,
			URL:    url,
			Path:   path,
			Body:   "{\"country_code\": \"{{country_code}}\"}",
		})

		for _, option := range RPCNode.Options {
			if option.OptionName == "(google.api.http)" {

				optionBodyRegex := regexp.MustCompile(`(get|post|patch|put|delete)\s*:\s*["'](.+)["']`)
				matches := optionBodyRegex.FindAllStringSubmatch(option.Constant, -1)

				if len(matches) == 1 && len(matches[0]) == 3 {
					methodVerb = matches[0][1]
					url = matches[0][2]

					// find all variables
					variableRegex := regexp.MustCompile(`\{(\S+?)\}`)
					variableMatches := variableRegex.FindAllStringSubmatch(url, -1)

					if variableMatches != nil {
						for _, match := range variableMatches {
							variables = append(variables, match[1])
						}
					}

					url = strings.Replace(url, "{", "{{", -1)
					url = strings.Replace(url, "}", "}}", -1)
					path = strings.Split(url, "/")

					body := ""
					if methodVerb != "GET" {
						body = "{}"
					}

					restEndpoints = append(restEndpoints, &pParser.Endpoint{
						Name:   "REST_" + methodName,
						Method: strings.ToUpper(methodVerb),
						URL:    url,
						Path:   path,
						Body:   body,
					})
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "Could not parse option body for method %s.%s, skipping adding the REST method\n", serviceName, methodName)
				}

				break
			}
		}
	}

	endpoints = append(endpoints, restEndpoints...)

	return &pParser.Service{
		Name:      serviceName,
		Endpoints: endpoints,
	}, variables
}
