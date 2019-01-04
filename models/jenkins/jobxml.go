package jobconfig

import (
	"encoding/xml"
)

type JobConfig struct {
	XMLName xml.Name `xml:"project"`
	//Properties      properties `xml:"properties"`
	Disabled        string `xml:"disabled"`
	ConcurrentBuild string `xml:"concurrentBuild"`
	AssignedNode    string `xml:"assignedNode"`
	Builders        builders
	Publishers      publishers
}

type properties struct {
	ParametersDefinitionProperty parametersDefinitionProperty `xml:"hudson.model.ParametersDefinitionProperty"`
}

type parametersDefinitionProperty struct {
	ParameterProperty parameterProperty
}

type parameterProperty struct {
	XMLName          xml.Name          `xml:"parameterDefinitions"`
	StringParameters []stringParameter `xml:"hudson.model.StringParameterDefinition"`
	ChoiceParameters []choiceParameter `xml:"hudson.model.ChoiceParameterDefinition"`
}

type stringParameter struct {
	Name         string `xml:"name"`
	Description  string `xml:"description"`
	DefaultValue string `xml:"defaultValue"`
}

type choiceParameter struct {
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Choices     choices
}

type choices struct {
	XMLName xml.Name `xml:"choices"`
	Class   string   `xml:"class,attr"`
	Aarray  aarray
}

type aarray struct {
	XMLName  xml.Name `xml:"a"`
	Class    string   `xml:"class,attr"`
	Children []string `xml:"string"`
}

type builders struct {
	XMLName xml.Name `xml:"builders"`
	Shell   shell    `xml:"hudson.tasks.Shell"`
}

type shell struct {
	Command string `xml:"command"`
}

type publishers struct {
	XMLName xml.Name `xml:"publishers"`
	Script  script   `xml:"org.jenkinsci.plugins.postbuildscript.PostBuildScript"`
}

type script struct {
	Plugin              string `xml:"plugin,attr"`
	BuildSteps          buildSteps
	ScriptOnlyIfSuccess string `xml:"scriptOnlyIfSuccess"`
	ScriptOnlyIfFailure string `xml:"scriptOnlyIfFailure"`
	MarkBuildUnstable   string `xml:"markBuildUnstable"`
}

type buildSteps struct {
	XMLName xml.Name `xml:"buildSteps"`
	Shell   shell    `xml:"hudson.tasks.Shell"`
}
