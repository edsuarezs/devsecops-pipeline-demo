# devsecops-pipeline-demo

[![CI](https://github.com/edsuarezs/devsecops-pipeline-demo/actions/workflows/ci.yml/badge.svg)](https://github.com/edsuarezs/devsecops-pipeline-demo/actions/workflows/ci.yml)
[![CD](https://github.com/edsuarezs/devsecops-pipeline-demo/actions/workflows/cd.yml/badge.svg)](https://github.com/edsuarezs/devsecops-pipeline-demo/actions/workflows/cd.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Security: Trivy](https://img.shields.io/badge/security-trivy-blue.svg)](https://github.com/aquasecurity/trivy)

> Production-grade DevSecOps pipeline for a FastAPI microservice deployed on Amazon EKS.
> Built as a reference implementation covering CI, CD, IaC, and security best practices.

---

## Table of Contents

- [devsecops-pipeline-demo](#devsecops-pipeline-demo)
  - [Table of Contents](#table-of-contents)
  - [Architecture](#architecture)
  - [Tech Stack](#tech-stack)
  - [Repository Structure](#repository-structure)
  - [Prerequisites](#prerequisites)
  - [Getting Started](#getting-started)
    - [1. Clone the repository](#1-clone-the-repository)
    - [2. Configure AWS credentials](#2-configure-aws-credentials)
    - [3. Provision infrastructure](#3-provision-infrastructure)
    - [4. Run the app locally](#4-run-the-app-locally)
  - [Pipeline Overview](#pipeline-overview)
    - [CI — Triggered on every push and pull request](#ci--triggered-on-every-push-and-pull-request)
    - [CD — Triggered on merge to `main`](#cd--triggered-on-merge-to-main)
  - [Security Controls](#security-controls)
  - [Infrastructure](#infrastructure)
  - [Helm Chart](#helm-chart)
  - [Contributing](#contributing)
  - [License](#license)

---

## Architecture

┌────────────────────────────────────────────────────────────┐
│                        GitHub                              │
│                                                            │
│   Push / PR        ┌──────────┐      ┌──────────┐          │
│  ───────────────►  │  CI Flow │─────►│  CD Flow │          │
│                    └──────────┘      └──────────┘          │
│                         │                  │               │
└─────────────────────────┼──────────────────┼───────────────┘
│                  │
┌─────▼─────┐      ┌─────▼─────┐
│  Amazon   │      │  Amazon   │
│    ECR    │      │    EKS    │
└───────────┘      └───────────┘
│
┌────────▼────────┐
│   Helm Release  │
│  (staging/prod) │
└─────────────────┘

---

## Tech Stack

| Layer                | Tool                         | Purpose                     |
| -------------------- | ---------------------------- | --------------------------- |
| **Application**      | Python 3.12 + FastAPI        | REST API                    |
| **Containerization** | Docker                       | Image build                 |
| **Registry**         | Amazon ECR                   | Image storage               |
| **Orchestration**    | Amazon EKS (Kubernetes 1.29) | Container runtime           |
| **IaC**              | Terraform >= 1.7             | Infrastructure provisioning |
| **Packaging**        | Helm 3                       | Kubernetes deployment       |
| **CI/CD**            | GitHub Actions               | Pipeline orchestration      |
| **SAST**             | Bandit + Semgrep             | Static code analysis        |
| **Image Scanning**   | Trivy                        | Vulnerability scanning      |
| **Secret Scanning**  | Gitleaks                     | Secrets detection           |
| **SBOM**             | Syft                         | Software Bill of Materials  |

---

## Repository Structure

devsecops-pipeline-demo/
├── .github/
│   ├── workflows/
│   │   ├── ci.yml          # Lint, test, SAST, build, push
│   │   └── cd.yml          # Helm deploy, smoke test, rollback
│   └── pull_request_template.md
├── app/                    # FastAPI application
├── docker/                 # Dockerfile
├── terraform/              # EKS + ECR infrastructure
├── helm/                   # Helm chart (staging + prod)
├── docs/                   # Architecture & runbook
├── .gitignore
├── .pre-commit-config.yaml
├── CHANGELOG.md
└── README.md

---

## Prerequisites

| Tool        | Version | Install                                                                      |
| ----------- | ------- | ---------------------------------------------------------------------------- |
| `git`       | >= 2.40 | [git-scm.com](https://git-scm.com)                                           |
| `docker`    | >= 24.0 | [docs.docker.com](https://docs.docker.com/get-docker/)                       |
| `terraform` | >= 1.7  | [developer.hashicorp.com](https://developer.hashicorp.com/terraform/install) |
| `helm`      | >= 3.14 | [helm.sh](https://helm.sh/docs/intro/install/)                               |
| `kubectl`   | >= 1.29 | [kubernetes.io](https://kubernetes.io/docs/tasks/tools/)                     |
| `aws-cli`   | >= 2.15 | [aws.amazon.com](https://aws.amazon.com/cli/)                                |

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/edsuarezs/devsecops-pipeline-demo.git
cd devsecops-pipeline-demo
```

### 2. Configure AWS credentials

```bash
aws configure
# AWS Access Key ID:
# AWS Secret Access Key:
# Default region: us-east-1
```

### 3. Provision infrastructure

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
terraform init
terraform plan
terraform apply
```

### 4. Run the app locally

```bash
cd app
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn main:app --reload
# API available at http://localhost:8000
# Docs available at http://localhost:8000/docs
```

---

## Pipeline Overview

### CI — Triggered on every push and pull request

Checkout → Lint (ruff) → Test (pytest) → SAST (Bandit + Semgrep)
→ Secret Scan (Gitleaks) → Docker Build → Image Scan (Trivy)
→ Push to ECR

### CD — Triggered on merge to `main`

Checkout → Configure AWS → Helm Lint → Helm Diff
→ Deploy to Staging → Smoke Test → Deploy to Prod
→ Notify (on success or failure)

---

## Security Controls

| Control          | Tool                     | Stage          |
| ---------------- | ------------------------ | -------------- |
| Secret detection | Gitleaks                 | CI — pre-push  |
| Static analysis  | Bandit + Semgrep         | CI             |
| Dependency audit | pip-audit                | CI             |
| Image scanning   | Trivy (CRITICAL block)   | CI             |
| SBOM generation  | Syft                     | CI             |
| Least privilege  | IAM roles + IRSA         | Infrastructure |
| Network policies | Kubernetes NetworkPolicy | Helm chart     |

---

## Infrastructure

Managed with Terraform. See [`terraform/`](terraform/) for full details.

- **EKS cluster** — managed node group, private subnets
- **ECR repository** — image scanning enabled, lifecycle policy
- **IAM** — IRSA (IAM Roles for Service Accounts), least privilege
- **VPC** — dedicated VPC, public/private subnets, NAT gateway

---

## Helm Chart

Two environment overlays:

```bash
# Staging
helm upgrade --install devsecops-pipeline-demo ./helm/devsecops-pipeline-demo \
  -f helm/devsecops-pipeline-demo/values-staging.yaml \
  --namespace staging

# Production
helm upgrade --install devsecops-pipeline-demo ./helm/devsecops-pipeline-demo \
  -f helm/devsecops-pipeline-demo/values-prod.yaml \
  --namespace prod
```

---

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Commit using Conventional Commits: `feat:`, `fix:`, `docs:`, `chore:`
4. Open a Pull Request — the CI pipeline runs automatically

---

## License

MIT License — see [LICENSE](LICENSE) for details.
