package types

import "time"

type CommandResult struct {
	CommandID     string
	Status        string
	ExitCode      int32
	StandardOut   string
	StandardErr   string
	ExecutionTime time.Duration
}

const EC2UserData = `#!/bin/bash
apt-get update -y
snap install amazon-ssm-agent --classic
systemctl enable snap.amazon-ssm-agent.amazon-ssm-agent.service
systemctl start snap.amazon-ssm-agent.amazon-ssm-agent.service`
