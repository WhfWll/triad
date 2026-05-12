package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// TreeNode represents a node in the tree structure
type TreeNode struct {
	Label    string      `json:"label"`
	Children []*TreeNode `json:"children"`
}

// ParsePath extracts and splits the path from the URL
func ParsePath(urlStr string) []string {
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil
	}
	path := strings.Split(u.Path, "/")
	// Filter out empty strings
	var filteredPath []string
	for _, part := range path {
		if part != "" {
			filteredPath = append(filteredPath, part)
		}
	}
	return filteredPath
}

// findOrCreateNode finds or creates a node in the tree
func findOrCreateNode(root *TreeNode, parts []string) *TreeNode {
	currentNode := root
	for _, part := range parts {
		found := false
		for _, child := range currentNode.Children {
			if child.Label == part {
				currentNode = child
				found = true
				break
			}
		}
		if !found {
			newNode := &TreeNode{Label: part}
			currentNode.Children = append(currentNode.Children, newNode)
			currentNode = newNode
		}
	}
	return currentNode
}

// BuildTree constructs the tree structure from the list of paths
func BuildTree(paths [][]string) *TreeNode {
	root := &TreeNode{Label: "root"}
	for _, path := range paths {
		findOrCreateNode(root, path)
	}
	return root
}

// ConvertTreeToSlice converts the tree to a slice, excluding the root node
func ConvertTreeToSlice(root *TreeNode) []*TreeNode {
	var result []*TreeNode
	for _, child := range root.Children {
		result = append(result, child)
	}
	return result
}

func main() {
	urls := []string{
		"http://192.168.3.3:8770/vul/xss/xss_reflected_get.php?message=crawler-uDElC&submit=submit",
		"http://192.168.3.3:8770/vul/fileinclude/fi_local.php?filename=file1.php&submit=crawler-fPHIf",
		"http://192.168.3.3:8770/assets/font-awesome/4.5.0/fonts/",
		"http://192.168.3.3:8770/vul/unsafedownload/execdownload.php?filename=pj.png",
		"http://192.168.3.3:8770/vul/xss/xss_02.php",
	}

	paths := make([][]string, len(urls))
	for i, urlStr := range urls {
		paths[i] = ParsePath(urlStr)
	}

	root := BuildTree(paths)
	tree := ConvertTreeToSlice(root)

	// Marshal the tree to JSON and print it
	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
