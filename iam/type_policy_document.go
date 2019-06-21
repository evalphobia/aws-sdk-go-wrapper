package iam

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	effectAllow = "Allow"
	effectDeny  = "Deny"
)

// PolicyDocument contains permission data of a policy.
type PolicyDocument struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

// used for multiple statements.
type documentForUnmarshal struct {
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

// used for single statement.
type documentForUnmarshalSingle struct {
	Version   string    `json:"Version"`
	Statement Statement `json:"Statement"`
}

// NewPolicyDocumentFromDocument returns initialized PolicyDocument from response data.
func NewPolicyDocumentFromDocument(document string) (PolicyDocument, error) {
	s, err := url.QueryUnescape(document)
	if err != nil {
		return PolicyDocument{}, err
	}
	return NewPolicyDocumentFromJSONString(s)
}

// NewPolicyDocumentFromJSONString returns initialized PolicyDocument from JSON data.
func NewPolicyDocumentFromJSONString(data string) (PolicyDocument, error) {
	p := PolicyDocument{}
	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		fmt.Printf("[%s]\n\n", string(data))
	}
	return p, err
}

// UnmarshalJSON converts from json to *PolicyDocument.
func (p *PolicyDocument) UnmarshalJSON(data []byte) error {
	d1 := documentForUnmarshal{}
	err := json.Unmarshal(data, &d1)
	if err == nil {
		p.Version = d1.Version
		p.Statement = d1.Statement
		return nil
	}

	d2 := documentForUnmarshalSingle{}
	err = json.Unmarshal(data, &d2)
	if err != nil {
		return err
	}

	p.Version = d2.Version
	p.Statement = []Statement{d2.Statement}
	return nil
}

// Statement represents statement of iam policy.
type Statement struct {
	Sid      string   `json:"Sid"`
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource []string `json:"Resource"`
}

// IsAllow checks that effect is allow.
func (s *Statement) IsAllow() bool {
	return s.Effect == effectAllow
}

// IsDeny checks that effect is deny.
func (s *Statement) IsDeny() bool {
	return s.Effect == effectDeny
}

// UnmarshalJSON converts from json to *Statement.
func (s *Statement) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	json.Unmarshal(data, &m)
	return s.setFromMap(m)
}

// used for converting from json data to *Statement
func (s *Statement) setFromMap(m map[string]interface{}) error {
	if v, ok := m["Sid"].(string); ok {
		s.Sid = v
	}
	if v, ok := m["Effect"].(string); ok {
		s.Effect = v
	}
	if v, ok := m["Action"]; ok {
		switch v := v.(type) {
		case []interface{}:
			s.Action = toStringList(v)
		case string:
			s.Action = []string{v}
		}
	}
	if v, ok := m["Resource"]; ok {
		switch v := v.(type) {
		case []interface{}:
			s.Resource = toStringList(v)
		case string:
			s.Resource = []string{v}
		}
	}
	return nil
}

func toStringList(list []interface{}) []string {
	result := make([]string, 0, len(list))
	for _, v := range list {
		if v, ok := v.(string); ok {
			result = append(result, v)
		}
	}
	return result
}
