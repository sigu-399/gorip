// author  			sigu-399
// author-github 	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description		Server's router implementation.
// 
// created      	04-03-2013

package rip

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	const_route_node_part         = ``
	const_route_element_separator = `/`
)

type router struct {
	rootNode                routerNode
	routeVariableValidators map[string]RouteVariableValidator
}

func newRouter() *router {
	r := &router{}
	r.routeVariableValidators = make(map[string]RouteVariableValidator)
	r.rootNode = newRouterNodeInvariable(r, const_route_node_part)
	return r
}

// Adds a endpoint to the router tree
func (r *router) RegisterEndpoint(endp *endpoint) error {

	routeString := endp.GetRoute()

	if !strings.HasPrefix(routeString, const_route_element_separator) {
		return errors.New(fmt.Sprintf(`A route must start with '%s'`, const_route_element_separator))
	}

	currentRouterNode := r.rootNode

	splitRouteString := strings.Split(routeString, const_route_element_separator)

	// Start parsing parts ( ommit root ( part : ``, route : `/` ) with 1: )
	for _, v := range splitRouteString[1:] {

		findChild := currentRouterNode.GetChildByPart(v, false)
		if findChild != nil {
			// found : move to next tree leaf
			currentRouterNode = findChild
		} else {
			// not found : have to create it and then move to the next tree leaf

			var newChild routerNode

			// Detecting routerNodeVariable
			if isRouteVariable(v) {
				rvIdentifier, rvKind, err := getRouteVariableParts(v)
				if err != nil {
					return err
				} else {
					if r.GetRouteVariableValidatorByKind(rvKind) == nil {
						return errors.New(fmt.Sprintf("Given route uses a route variable kind '%s' that was not registered", rvKind))
					} else {
						newChild = newRouterNodeVariable(r, v, rvIdentifier, rvKind)
					}
				}

			} else { // Otherwise routerNodeInvariable

				newChild = newRouterNodeInvariable(r, v)
			}

			// Add the new node to the tree

			err := currentRouterNode.AddChild(newChild)
			if err != nil {
				return err
			}

			currentRouterNode = newChild

		}

	}

	// Link the endpoint to the last node

	if currentRouterNode.GetEndpoint() != nil {
		return errors.New(fmt.Sprintf(`Endpoint already exists on route %s`, routeString))
	} else {
		currentRouterNode.SetEndpoint(endp)
	}

	return nil

}

// Adds a route variable validator to the router
func (r *router) RegisterRouteVariableValidator(kind string, validator RouteVariableValidator) error {

	log.Printf("Registering route variable validator of kind '%s'\n", kind)

	if r.GetRouteVariableValidatorByKind(kind) != nil {
		return errors.New(fmt.Sprintf(`Route variable validator of kind '%s' already exists`, kind))
	} else {
		r.routeVariableValidators[kind] = validator
	}

	return nil
}

func (r *router) GetRouteVariableValidatorByKind(kind string) RouteVariableValidator {

	if _, ok := r.routeVariableValidators[kind]; ok {
		return r.routeVariableValidators[kind]
	}

	return nil

}

// Find a matching route given url
func (r *router) FindNodeByRoute(routeString string) (routerNode, map[string]string, error) {

	routeVariableMap := make(map[string]string)

	// Check root
	if routeString == const_route_element_separator {
		return r.rootNode, routeVariableMap, nil
	}

	splitRouteString := strings.Split(routeString, const_route_element_separator)

	// Start parsing parts ( ommit root ( part : ``, route : `/` ) with 1: )
	currentRouterNode := r.rootNode
	for _, v := range splitRouteString[1:] {

		foundChild := currentRouterNode.GetChildByPart(v, true)

		if foundChild != nil {

			// Extract route variable if any
			switch foundChild.(type) {
			case *routerNodeVariable:
				variable := foundChild.(*routerNodeVariable)
				routeVariableMap[variable.identifier] = v
			}

			currentRouterNode = foundChild
		} else {
			return nil, nil, errors.New(fmt.Sprintf(`Could not find a route given the part '%s'`, v))
		}
	}

	return currentRouterNode, routeVariableMap, nil
}

// Displays the resulting router tree in the log

func (r *router) PrintRouterTree() {

	log.Printf("Router tree : \n")
	r.printRouterTreeRecursive(r.rootNode, "", 0)

}

func (r *router) printRouterTreeRecursive(node routerNode, text string, level int) {

	indent := ``
	for l := 0; l != level; l++ {
		indent += ` `
	}

	log.Printf("%s/%s\n", indent, text)

	children := node.GetChildren()

	for _, value := range children {
		r.printRouterTreeRecursive(value, value.GetPart(), level+1)
	}

}

type routerNode interface {
	GetRouter() *router
	GetPart() string
	GetChildren() map[string]routerNode
	AddChild(routerNode) error
	GetChildByPart(part string, invariableMode bool) routerNode
	SetEndpoint(endp *endpoint)
	GetEndpoint() *endpoint
}

type routerNodeImplementation struct {
	part         string
	children     map[string]routerNode
	parentRouter *router
	endp         *endpoint
}

// Initialize the interface implementation of routerNode, 
// must be called as a constructor function by all the sub structs implementing the routerNode interface
func (rni *routerNodeImplementation) Initialize(r *router, part string, isVariable bool) *routerNodeImplementation {
	rni.part = part
	rni.children = make(map[string]routerNode)
	rni.parentRouter = r
	return rni
}

func (rni *routerNodeImplementation) GetPart() string {
	return rni.part
}

func (rni *routerNodeImplementation) SetEndpoint(endp *endpoint) {
	rni.endp = endp
}

func (rni *routerNodeImplementation) GetEndpoint() *endpoint {
	return rni.endp
}

func (rni *routerNodeImplementation) GetNodeImplementation() *routerNodeImplementation {
	return rni
}

func (rni *routerNodeImplementation) GetChildren() map[string]routerNode {
	return rni.children
}

func (rni *routerNodeImplementation) AddChild(child routerNode) error {

	if rni.GetChildByPart(child.GetPart(), false) != nil {
		return errors.New(fmt.Sprintf(`A child '%s' already exists`, child.GetPart()))
	} else {
		rni.children[child.GetPart()] = child
	}

	return nil

}

func (rni *routerNodeImplementation) GetRouter() *router {
	return rni.parentRouter
}

func (rni *routerNodeImplementation) GetChildByPart(part string, invariableMode bool) routerNode {

	var nodeFound routerNode

	if invariableMode {
		// Check variable ones first
		for k := range rni.children {
			child := rni.children[k]

			switch child.(type) {
			case *routerNodeVariable:
				variable := child.(*routerNodeVariable)
				validator := child.GetRouter().GetRouteVariableValidatorByKind(variable.kind)
				if validator.Matches(part) {
					if nodeFound != nil {
						log.Printf("Warning : Multiple routings for a given route")
					}
					nodeFound = child
				}
			default:
			}
		}
	}

	// Check invariable ones
	if _, ok := rni.children[part]; ok {
		if nodeFound != nil {
			log.Printf("Warning : Multiple routings for a given route")
		}
		return rni.children[part]
	}

	return nodeFound

}

type routerNodeInvariable struct {
	routerNodeImplementation
}

func newRouterNodeInvariable(r *router, part string) *routerNodeInvariable {
	rniv := &routerNodeInvariable{}
	rniv.routerNodeImplementation.Initialize(r, part, false)
	return rniv
}

type routerNodeVariable struct {
	routerNodeImplementation
	identifier string
	kind       string
}

func newRouterNodeVariable(r *router, part string, identifier string, kind string) *routerNodeVariable {
	rnva := &routerNodeVariable{identifier: identifier, kind: kind}
	rnva.routerNodeImplementation.Initialize(r, part, false)
	return rnva
}

type RouteVariableValidator interface {
	Matches(string) bool
}

const (
	const_regexp_route_variable_pattern       = "\\{(.*?)\\}"
	const_regexp_route_variable_parts_pattern = "\\{([0-9a-zA-Z_]*)\\:([0-9a-zA-Z_]*)\\}"
)

var regexpRouteVariable *regexp.Regexp      // anything like {...}
var regexpRouteVariableParts *regexp.Regexp // {identifier:kind} they both accepts alpha and numerals (a-z A-Z 0-9) with optional _

func isRouteVariable(part string) bool {
	return regexpRouteVariable.MatchString(part)
}

func getRouteVariableParts(part string) (string, string, error) {

	matches := regexpRouteVariableParts.FindAllStringSubmatch(part, 2)

	// It is valid
	if len(matches) == 1 {
		return matches[0][1], matches[0][2], nil
	}

	// Otherwise throws an error
	return "", "", errors.New(fmt.Sprintf(`Part %s is not a valid route variable definition`, part))

}

func init() {

	var err error

	regexpRouteVariable, err = regexp.Compile(const_regexp_route_variable_pattern)
	if err != nil {
		panic("Could not compile regexpRouteVariable")
	}

	regexpRouteVariableParts, err = regexp.Compile(const_regexp_route_variable_parts_pattern)
	if err != nil {
		panic("Could not compile regexpRouteVariableParts")
	}

}
