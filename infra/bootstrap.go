package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const stackName = "gymshark-saddiqs1-bootstrap"

type config struct {
	awsProfile         string
	awsRegion          string
	githubOwner        string
	githubOwnerID      string
	githubRepository   string
	githubRepositoryID string
	githubBranch       string
}

func run(input string, name string, args ...string) (string, error) {
	command := exec.Command(name, args...)
	if input != "" {
		command.Stdin = strings.NewReader(input)
	}

	var stdout bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return "", fmt.Errorf("command failed (%s %s): %w", name, strings.Join(args, " "), err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func stackOutput(cfg config, key string) (string, error) {
	value, err := run(
		"",
		"aws",
		"cloudformation",
		"describe-stacks",
		"--stack-name", stackName,
		"--region", cfg.awsRegion,
		"--profile", cfg.awsProfile,
		"--query", fmt.Sprintf("Stacks[0].Outputs[?OutputKey=='%s'].OutputValue", key),
		"--output", "text",
	)
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", fmt.Errorf("CloudFormation output %s was empty", key)
	}
	return value, nil
}

func templatePath() (string, error) {
	fromRepositoryRoot := filepath.Join("infra", "bootstrap.yml")
	if _, err := os.Stat(fromRepositoryRoot); err == nil {
		return fromRepositoryRoot, nil
	}

	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	besideExecutable := filepath.Join(filepath.Dir(executable), "bootstrap.yml")
	if _, err := os.Stat(besideExecutable); err == nil {
		return besideExecutable, nil
	}

	return "", errors.New("could not find infra/bootstrap.yml; run this command from the repository root")
}

func bootstrap(cfg config) error {
	for _, executable := range []string{"aws", "gh"} {
		if _, err := exec.LookPath(executable); err != nil {
			return fmt.Errorf("required executable not found: %s", executable)
		}
	}

	template, err := templatePath()
	if err != nil {
		return err
	}

	fmt.Printf("Checking AWS authentication for profile %q...\n", cfg.awsProfile)
	if _, err := run(
		"",
		"aws",
		"sts",
		"get-caller-identity",
		"--profile", cfg.awsProfile,
		"--region", cfg.awsRegion,
	); err != nil {
		return fmt.Errorf("AWS authentication failed; run 'aws sso login --profile %s': %w", cfg.awsProfile, err)
	}

	fmt.Printf("Deploying bootstrap stack %q...\n", stackName)
	if _, err := run(
		"",
		"aws",
		"cloudformation",
		"deploy",
		"--stack-name", stackName,
		"--template-file", template,
		"--capabilities", "CAPABILITY_NAMED_IAM",
		"--region", cfg.awsRegion,
		"--profile", cfg.awsProfile,
		"--parameter-overrides",
		"GitHubOwner="+cfg.githubOwner,
		"GitHubOwnerId="+cfg.githubOwnerID,
		"GitHubRepository="+cfg.githubRepository,
		"GitHubRepositoryId="+cfg.githubRepositoryID,
		"GitHubBranch="+cfg.githubBranch,
	); err != nil {
		return err
	}

	stateBucket, err := stackOutput(cfg, "TerraformStateBucketName")
	if err != nil {
		return err
	}
	deployRoleARN, err := stackOutput(cfg, "GitHubDeployRoleArn")
	if err != nil {
		return err
	}
	planRoleARN, err := stackOutput(cfg, "GitHubPlanRoleArn")
	if err != nil {
		return err
	}

	repository := cfg.githubOwner + "/" + cfg.githubRepository
	fmt.Printf("Configuring GitHub Actions secrets for %q...\n", repository)
	if _, err := run(stateBucket, "gh", "secret", "set", "TF_STATE_BUCKET", "--repo", repository); err != nil {
		return err
	}
	if _, err := run(deployRoleARN, "gh", "secret", "set", "AWS_DEPLOY_ROLE_ARN", "--repo", repository); err != nil {
		return err
	}
	if _, err := run(planRoleARN, "gh", "secret", "set", "AWS_PLAN_ROLE_ARN", "--repo", repository); err != nil {
		return err
	}

	fmt.Println("Bootstrap complete. GitHub Actions can now be used to provision infra via terraform.")
	return nil
}

func main() {
	cfg := config{}
	flag.StringVar(&cfg.awsProfile, "aws-profile", "", "AWS CLI profile to use (required)")
	flag.StringVar(&cfg.awsRegion, "aws-region", "eu-west-1", "AWS region")
	flag.StringVar(&cfg.githubOwner, "github-owner", "saddiqs1", "GitHub repository owner")
	flag.StringVar(&cfg.githubOwnerID, "github-owner-id", "18683829", "Immutable GitHub repository owner ID")
	flag.StringVar(&cfg.githubRepository, "github-repository", "gymshark-saddiqs1", "GitHub repository name")
	flag.StringVar(&cfg.githubRepositoryID, "github-repository-id", "1333598191", "Immutable GitHub repository ID")
	flag.StringVar(&cfg.githubBranch, "github-branch", "main", "Branch allowed to deploy")
	flag.Parse()

	if cfg.awsProfile == "" {
		fmt.Fprintln(os.Stderr, "--aws-profile is required")
		flag.Usage()
		os.Exit(2)
	}

	if err := bootstrap(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
