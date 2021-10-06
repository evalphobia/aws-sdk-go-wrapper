package iam

import (
	"github.com/aws/aws-sdk-go/aws/session"
	SDK "github.com/aws/aws-sdk-go/service/iam"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	serviceName = "IAM"
)

// IAM has IAM client.
type IAM struct {
	client *SDK.IAM

	logger log.Logger
}

// New returns initialized *IAM.
func New(conf config.Config) (*IAM, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	return NewFromSession(sess), nil
}

// NewFromSession returns initialized *IAM from aws.Session.
func NewFromSession(sess *session.Session) *IAM {
	return &IAM{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
	}
}

// SetLogger sets logger.
func (svc *IAM) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// ListUsers fetches all of the user list.
func (svc *IAM) ListUsers() ([]User, error) {
	return svc.listUsers(&SDK.ListUsersInput{})
}

// listUsers executes listUsers operation.
func (svc *IAM) listUsers(input *SDK.ListUsersInput) ([]User, error) {
	// set default limit
	if input.MaxItems == nil {
		input.MaxItems = pointers.Long64(1000)
	}

	output, err := svc.client.ListUsers(input)
	if err != nil {
		svc.Errorf("error on `ListUsers` operation; error=%s;", err.Error())
		return nil, err
	}
	return NewUsers(output.Users), nil
}

// ListUserPolicies fetches inline policies of the user.
func (svc *IAM) ListUserPolicies(userName string) ([]string, error) {
	return svc.listUserPolicies(&SDK.ListUserPoliciesInput{
		UserName: pointers.String(userName),
	})
}

// listUserPolicies executes ListUserPolicies operation.
func (svc *IAM) listUserPolicies(input *SDK.ListUserPoliciesInput) ([]string, error) {
	// set default limit
	if input.MaxItems == nil {
		input.MaxItems = pointers.Long64(1000)
	}

	output, err := svc.client.ListUserPolicies(input)
	if err != nil {
		svc.Errorf("error on `ListUserPolicies` operation; error=%s;", err.Error())
		return nil, err
	}

	list := make([]string, 0, len(output.PolicyNames))
	for _, p := range output.PolicyNames {
		if p != nil {
			list = append(list, *p)
		}
	}

	return list, nil
}

// GetGroup executes GetGroup operation.
func (svc *IAM) GetGroup(groupName string) (*SDK.GetGroupOutput, error) {
	output, err := svc.client.GetGroup(&SDK.GetGroupInput{
		GroupName: pointers.String(groupName),
	})
	if err != nil {
		svc.Errorf("error on `GetGroup` operation; error=[%s]; groupName=[%s];", err.Error(), groupName)
		return nil, err
	}
	return output, nil
}

// ListGroups fetches all of the group list.
func (svc *IAM) ListGroups() ([]Group, error) {
	return svc.listGroups(&SDK.ListGroupsInput{})
}

// listGroups executes listGroups operation.
func (svc *IAM) listGroups(input *SDK.ListGroupsInput) ([]Group, error) {
	// set default limit
	if input.MaxItems == nil {
		input.MaxItems = pointers.Long64(1000)
	}

	output, err := svc.client.ListGroups(input)
	if err != nil {
		svc.Errorf("error on `ListGroups` operation; error=%s;", err.Error())
		return nil, err
	}
	return NewGroups(output.Groups), nil
}

// ListGroupPolicies fetches inline policies of the user.
func (svc *IAM) ListGroupPolicies(groupName string) ([]string, error) {
	return svc.listGroupPolicies(&SDK.ListGroupPoliciesInput{
		GroupName: pointers.String(groupName),
	})
}

// listGroupPolicies executes ListGroupPolicies operation.
func (svc *IAM) listGroupPolicies(input *SDK.ListGroupPoliciesInput) ([]string, error) {
	// set default limit
	if input.MaxItems == nil {
		input.MaxItems = pointers.Long64(1000)
	}

	output, err := svc.client.ListGroupPolicies(input)
	if err != nil {
		svc.Errorf("error on `ListGroupPolicies` operation; error=%s;", err.Error())
		return nil, err
	}

	list := make([]string, 0, len(output.PolicyNames))
	for _, p := range output.PolicyNames {
		if p != nil {
			list = append(list, *p)
		}
	}

	return list, nil
}

// ListRoles fetches all of the role list.
func (svc *IAM) ListRoles() ([]Role, error) {
	return svc.listRoles(&SDK.ListRolesInput{})
}

// listRoles executes listRoles operation.
func (svc *IAM) listRoles(input *SDK.ListRolesInput) ([]Role, error) {
	// set default limit
	if input.MaxItems == nil {
		input.MaxItems = pointers.Long64(1000)
	}

	output, err := svc.client.ListRoles(input)
	if err != nil {
		svc.Errorf("error on `ListRoles` operation; error=%s;", err.Error())
		return nil, err
	}
	return NewRoles(output.Roles), nil
}

// ListRolePolicies fetches inline policies of the user.
func (svc *IAM) ListRolePolicies(roleName string) ([]string, error) {
	return svc.listRolePolicies(&SDK.ListRolePoliciesInput{
		RoleName: pointers.String(roleName),
	})
}

// listRolePolicies executes ListRolePolicies operation.
func (svc *IAM) listRolePolicies(input *SDK.ListRolePoliciesInput) ([]string, error) {
	// set default limit
	if input.MaxItems == nil {
		input.MaxItems = pointers.Long64(1000)
	}

	output, err := svc.client.ListRolePolicies(input)
	if err != nil {
		svc.Errorf("error on `ListRolePolicies` operation; error=%s;", err.Error())
		return nil, err
	}

	list := make([]string, 0, len(output.PolicyNames))
	for _, p := range output.PolicyNames {
		if p != nil {
			list = append(list, *p)
		}
	}

	return list, nil
}

// GetPolicyVersion executes GetPolicyVersion operation.
func (svc *IAM) GetPolicyVersion(arn, versionID string) (*SDK.PolicyVersion, error) {
	output, err := svc.client.GetPolicyVersion(&SDK.GetPolicyVersionInput{
		PolicyArn: pointers.String(arn),
		VersionId: pointers.String(versionID),
	})
	if err != nil {
		svc.Errorf("error on `GetPolicyVersion` operation; error=[%s]; arn=[%s];", err.Error())
		return nil, err
	}
	return output.PolicyVersion, nil
}

// GetUserPolicyDocument fetched Statement from user's inline policy.
func (svc *IAM) GetUserPolicyDocument(userName, policyName string) (*PolicyDocument, error) {
	output, err := svc.getUserPolicy(userName, policyName)
	switch {
	case err != nil:
		return nil, err
	case output == nil,
		output.PolicyDocument == nil:
		return nil, nil
	}

	doc, err := NewPolicyDocumentFromDocument(*output.PolicyDocument)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// getUserPolicy executes GetUserPolicy operation.
func (svc *IAM) getUserPolicy(userName, policyName string) (*SDK.GetUserPolicyOutput, error) {
	output, err := svc.client.GetUserPolicy(&SDK.GetUserPolicyInput{
		UserName:   pointers.String(userName),
		PolicyName: pointers.String(policyName),
	})
	if err != nil {
		svc.Errorf("error on `GetUserPolicy` operation; error=[%s]; arn=[%s];", err.Error())
	}
	return output, err
}

// GetGroupPolicyDocument fetched Statement from user's inline policy.
func (svc *IAM) GetGroupPolicyDocument(groupName, policyName string) (*PolicyDocument, error) {
	output, err := svc.getGroupPolicy(groupName, policyName)
	switch {
	case err != nil:
		return nil, err
	case output == nil,
		output.PolicyDocument == nil:
		return nil, nil
	}

	doc, err := NewPolicyDocumentFromDocument(*output.PolicyDocument)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// getGroupPolicy executes GetGroupPolicy operation.
func (svc *IAM) getGroupPolicy(groupName, policyName string) (*SDK.GetGroupPolicyOutput, error) {
	output, err := svc.client.GetGroupPolicy(&SDK.GetGroupPolicyInput{
		GroupName:  pointers.String(groupName),
		PolicyName: pointers.String(policyName),
	})
	if err != nil {
		svc.Errorf("error on `GetGroupPolicy` operation; error=[%s]; arn=[%s];", err.Error())
	}
	return output, err
}

// GetRolePolicyDocument fetched Statement from user's inline policy.
func (svc *IAM) GetRolePolicyDocument(roleName, policyName string) (*PolicyDocument, error) {
	output, err := svc.getRolePolicy(roleName, policyName)
	switch {
	case err != nil:
		return nil, err
	case output == nil,
		output.PolicyDocument == nil:
		return nil, nil
	}

	doc, err := NewPolicyDocumentFromDocument(*output.PolicyDocument)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// getRolePolicy executes GetRolePolicy operation.
func (svc *IAM) getRolePolicy(roleName, policyName string) (*SDK.GetRolePolicyOutput, error) {
	output, err := svc.client.GetRolePolicy(&SDK.GetRolePolicyInput{
		RoleName:   pointers.String(roleName),
		PolicyName: pointers.String(policyName),
	})
	if err != nil {
		svc.Errorf("error on `GetGroupPolicy` operation; error=[%s]; arn=[%s];", err.Error())
	}
	return output, err
}

// ListAllPolicies fetches all of the policies list.
func (svc *IAM) ListAllPolicies() ([]Policy, error) {
	return svc.listPolicies(&SDK.ListPoliciesInput{})
}

// ListAttachedPolicies fetches attached policy list.
func (svc *IAM) ListAttachedPolicies() ([]Policy, error) {
	return svc.listPolicies(&SDK.ListPoliciesInput{
		OnlyAttached: pointers.Bool(true),
	})
}

// listPolicies executes ListPolicies operation.
func (svc *IAM) listPolicies(input *SDK.ListPoliciesInput) ([]Policy, error) {
	// set default limit
	if input.MaxItems == nil {
		input.MaxItems = pointers.Long64(1000)
	}

	output, err := svc.client.ListPolicies(input)
	if err != nil {
		svc.Errorf("error on `ListPolicies` operation; error=%s;", err.Error())
		return nil, err
	}
	return NewPolicies(output.Policies), nil
}

// ListEntitiesForPolicy executes ListEntitiesForPolicy operation.
func (svc *IAM) ListEntitiesForPolicy(arn string) ([]PolicyEntity, error) {
	output, err := svc.client.ListEntitiesForPolicy(&SDK.ListEntitiesForPolicyInput{
		PolicyArn: pointers.String(arn),
	})
	if err != nil {
		svc.Errorf("error on `ListEntitiesForPolicy` operation; error=%s;", err.Error())
		return nil, err
	}
	return NewPolicyEntityList(output), nil
}

// Infof logging information.
func (svc *IAM) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *IAM) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}
