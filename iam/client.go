package iam

import (
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

	svc := &IAM{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
	}
	return svc, nil
}

// SetLogger sets logger.
func (svc *IAM) SetLogger(logger log.Logger) {
	svc.logger = logger
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
