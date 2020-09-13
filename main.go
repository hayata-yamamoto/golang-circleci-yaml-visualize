package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/goccy/go-yaml"
)

// Config is ...
type Config struct {
	Version   string              `yaml:"version"`
	Jobs      map[string]Job      `yaml:"jobs"`
	Workflows map[string]Workflow `yaml:"workflows"`
}

// Job is ...
type Job struct {
	Docker interface{} `yaml:"docker"`
	Steps  interface{} `yaml:"steps"`
}

// Workflow is ...
type Workflow struct {
	Jobs []string `yaml:"jobs"`
}

func readYaml(filename string) Config {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		fmt.Println(err)
	}
	return c
}

func newGraph() (*graphviz.Graphviz, *cgraph.Graph) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	return g, graph
}

func writeDot(g *graphviz.Graphviz, graph *cgraph.Graph, filename string) {
	var buf bytes.Buffer
	if err := g.Render(graph, "dot", &buf); err != nil {
		log.Fatal(err)
	}

	_, err := g.RenderImage(graph)
	if err != nil {
		log.Fatal(err)
	}

	if err := g.RenderFilename(graph, graphviz.PNG, filename); err != nil {
		log.Fatal(err)
	}

}

func main() {
	data := readYaml("test.yml")
	g, graph := newGraph()

	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	nodes := make(map[string]*cgraph.Node)
	// Generate Nodes
	for k := range data.Jobs {
		e, err := graph.CreateNode(k)
		if err != nil {
			fmt.Println(err)
		}
		nodes[k] = e
	}

	for _, v := range data.Workflows {
		for i := 0; i < len(v.Jobs)-1; i++ {
			graph.CreateEdge("e", nodes[v.Jobs[i]], nodes[v.Jobs[i+1]])
		}
	}

	writeDot(g, graph, "graph.png")

}
