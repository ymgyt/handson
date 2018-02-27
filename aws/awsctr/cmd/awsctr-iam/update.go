package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/YmgchiYt/aws-handson/awsctr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/google/subcommands"
	"github.com/juju/errors"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

type UpdateCmd struct {
	log                 *zap.Logger
	generateCfgTemplate bool
}

func (u *UpdateCmd) Name() string     { return "update" }
func (u *UpdateCmd) Synopsis() string { return "update iam principal(user, group, role)" }
func (u *UpdateCmd) Usage() string {
	return fmt.Sprintf("Usage: %s update config.yml", os.Args[0])
}

func (u *UpdateCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&u.generateCfgTemplate, "generate-config", false, "generate update config template")
}

type UpdateCfg struct {
	Principal *PrincipalCfg `yaml:"Principal"`
	Policy    *PolicyCfg    `yaml:"Policy"`
}

type PrincipalCfg struct {
	Type string `yaml:"Type"` // user, group, role
	Path string `yaml:"Path"`
	Name string `yaml:"Name"`
}

type PolicyCfg struct {
	Description    string `yaml:"Description"`
	Path           string `yaml:"Path"`
	PolicyName     string `yaml:"PolicyName"`
	PolicyDocument string `yaml:"PolicyDocument"`
}

var updateCfgTemplate = `---
Principal:
  # user, group, role
  Type: 
  Path:
  Name:

Policy:
  Description:
  Path:
  PolicyName:
  PolicyDocument:
`

func (u *UpdateCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if u.generateCfgTemplate {
		fmt.Print(updateCfgTemplate)
		return subcommands.ExitSuccess
	}

	if f.NArg() != 1 {
		fmt.Fprintln(os.Stderr, u.Usage())
		return subcommands.ExitUsageError
	}

	var (
		cfg    UpdateCfg
		client *iam.IAM
		err    error
	)

	if err = awsctr.ReadCfgFromYAMLFile(f.Arg(0), &cfg); err != nil {
		u.log.Error("update", zap.Error(err))
		return subcommands.ExitFailure
	}

	client, err = awsctr.NewIAMClient()
	if err != nil {
		u.log.Error("update", zap.Error(err))
		return subcommands.ExitFailure
	}

	if err = u.updatePrincipal(client, cfg.Principal); err != nil {
		u.log.Error("update", zap.Error(err))
		return subcommands.ExitFailure
	}

	if err = u.updateAndAttachPolicy(client, &cfg); err != nil {
		u.log.Error("update", zap.Error(err))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (u *UpdateCmd) updateAndAttachPolicy(client *iam.IAM, cfg *UpdateCfg) (err error) {
	switch principal := strings.ToLower(cfg.Principal.Type); principal {
	case "user":
		err = u.updateAndAttachUserPolicy(client, cfg.Principal.Name, cfg.Policy)
	case "group":
		err = errors.NotImplementedf("update group")
	case "role":
		err = errors.NotImplementedf("update role")
	default:
		err = errors.Errorf("unexpected principal type %q", principal)
	}
	return errors.Trace(err)
}

func (u *UpdateCmd) updateAndAttachUserPolicy(client *iam.IAM, username string, cfg *PolicyCfg) error {
	createInput := &iam.CreatePolicyInput{
		Description:    aws.String(cfg.Description),
		Path:           aws.String(cfg.Path),
		PolicyName:     aws.String(cfg.PolicyName),
		PolicyDocument: aws.String(cfg.PolicyDocument),
	}
	_, err := client.CreatePolicy(createInput)
	if err == nil {
		u.log.Info("update",
			zap.String("msg", "iam policy successfully created"),
			zap.String("policy_name", cfg.PolicyName),
		)
	}
	if err != nil {
		if awsErr, ok := errors.Cause(err).(awserr.Error); ok {
			switch code := awsErr.Code(); code {
			case "EntityAlreadyExists":
				u.log.Info("update",
					zap.String("msg", "iam policy alredy exists"),
					zap.String("policy_name", cfg.PolicyName),
				)
				err = nil
			}
		}
		if err != nil {
			return errors.Trace(err)
		}
	}

	// policyを単体でupdateするAPIがみつからないので、Attachする際に更新する
	// CreatePolicyVersion !!!
	putInput := &iam.PutUserPolicyInput{
		PolicyDocument: aws.String(cfg.PolicyDocument),
		PolicyName:     aws.String(cfg.PolicyName),
		UserName:       aws.String(username),
	}
	putOutput, err := client.PutUserPolicy(putInput)
	if err != nil {
		return errors.Trace(err)
	}
	pp.Println(putOutput)

	return nil
}

func (u *UpdateCmd) createUser(client *iam.IAM, cfg *PrincipalCfg) error {
	input := &iam.CreateUserInput{
		Path:     aws.String(cfg.Path),
		UserName: aws.String(cfg.Name),
	}
	output, err := client.CreateUser(input)
	if err != nil {
		return errors.Trace(err)
	}
	u.log.Info("update",
		zap.String("msg", "iam user successfully created"),
		zap.String("arn", aws.StringValue(output.User.Arn)),
		zap.String("path", aws.StringValue(output.User.Path)),
		zap.String("user_id", aws.StringValue(output.User.UserId)),
		zap.String("username", aws.StringValue(output.User.UserName)),
	)
	return nil
}

func (u *UpdateCmd) updatePrincipal(client *iam.IAM, cfg *PrincipalCfg) (err error) {
	switch principal := strings.ToLower(cfg.Type); principal {
	case "user":
		err = u.createUser(client, cfg)
		if awsErr, ok := errors.Cause(err).(awserr.Error); ok {
			switch code := awsErr.Code(); code {
			case "EntityAlreadyExists":
				u.log.Info("update",
					zap.String("msg", "iam user alredy exists"),
					zap.String("username", cfg.Name),
				)
				err = nil
			}
		}
	case "group":
		err = errors.NotImplementedf("update group")
	case "role":
		err = errors.NotImplementedf("update role")
	default:
		err = errors.Errorf("unexpected principal type %q", principal)
	}
	return err
}
