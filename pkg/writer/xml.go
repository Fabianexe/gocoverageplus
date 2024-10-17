package writer

import (
	"encoding/xml"
)

type Coverage struct {
	XMLName xml.Name `xml:"coverage"`

	Sources  *Sources  `xml:"sources"`
	Packages *Packages `xml:"packages"`

	LineRate        string `xml:"line-rate,attr"`
	BranchRate      string `xml:"branch-rate,attr"`
	LinesCovered    string `xml:"lines-covered,attr"`
	LinesValid      string `xml:"lines-valid,attr"`
	BranchesCovered string `xml:"branches-covered,attr"`
	BranchesValid   string `xml:"branches-valid,attr"`
	Complexity      string `xml:"complexity,attr"`
	Version         string `xml:"version,attr"`
	Timestamp       string `xml:"timestamp,attr"`
}

type Sources struct {
	Sources []*Source `xml:"source"`
}

type Source struct {
	Path string `xml:",chardata"`
}

type Packages struct {
	Packages []*Package `xml:"package"`
}

type Package struct {
	Classes *Classes `xml:"classes"`

	Name       string `xml:"name,attr"`
	LineRate   string `xml:"line-rate,attr"`
	BranchRate string `xml:"branch-rate,attr"`
	Complexity string `xml:"complexity,attr"`
}

type Classes struct {
	Classes []*Class `xml:"class"`
}

type Class struct {
	Methods *Methods `xml:"methods"`
	Lines   *Lines   `xml:"lines"`

	Name       string `xml:"name,attr"`
	Filename   string `xml:"filename,attr"`
	LineRate   string `xml:"line-rate,attr"`
	BranchRate string `xml:"branch-rate,attr"`
	Complexity string `xml:"complexity,attr"`
}

type Methods struct {
	Methods []*Method `xml:"method"`
}

type Method struct {
	Lines *Lines `xml:"lines"`

	Name       string `xml:"name,attr"`
	Signature  string `xml:"signature,attr"`
	LineRate   string `xml:"line-rate,attr"`
	BranchRate string `xml:"branch-rate,attr"`
	Complexity string `xml:"complexity,attr"`
}

type Lines struct {
	Lines []*Line `xml:"line"`
}

type Line struct {
	Conditions *Conditions `xml:"conditions"`

	Number            string `xml:"number,attr"`
	Hits              string `xml:"hits,attr"`
	Branch            string `xml:"branch,attr"`
	ConditionCoverage string `xml:"condition-coverage,attr"`
}

type Conditions struct {
	Methods []*Condition `xml:"condition"`
}

type Condition struct {
	Number   string `xml:"number,attr"`
	Type     string `xml:"type,attr"`
	Coverage string `xml:"coverage,attr"`
}
