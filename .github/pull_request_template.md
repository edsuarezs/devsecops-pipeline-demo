# Pull Request

## Description

<!-- Describe what this PR does and why. Link the related issue if applicable. -->

Closes #<!-- issue number -->

---

## Type of change

- [ ] `feat` — New feature
- [ ] `fix` — Bug fix
- [ ] `docs` — Documentation only
- [ ] `chore` — Maintenance, dependencies, config
- [ ] `refactor` — Code restructure without behavior change
- [ ] `test` — Adding or updating tests
- [ ] `ci` — CI/CD pipeline changes
- [ ] `sec` — Security improvement

---

## Checklist

### Code quality

- [ ] My code follows the project style (golangci-lint passes locally)
- [ ] I have added/updated tests for my changes
- [ ] All existing tests pass locally (`go test ./...`)
- [ ] I have updated documentation if needed

### Security

- [ ] I have not hardcoded secrets, tokens, or credentials
- [ ] I ran Gitleaks locally and found no secrets (`gitleaks detect`)
- [ ] I ran gosec locally and addressed findings (`gosec ./...`)
- [ ] New dependencies have been reviewed for known vulnerabilities

### Infrastructure (if applicable)

- [ ] `terraform plan` runs without errors
- [ ] `terraform fmt` has been applied
- [ ] `terraform validate` passes
- [ ] No sensitive values are hardcoded in `.tf` files

### Helm (if applicable)

- [ ] `helm lint` passes
- [ ] `values.yaml` does not contain environment-specific secrets

---

## How to test

<!-- Step-by-step instructions for the reviewer to validate this PR -->

1.
2.
3.

---

## Screenshots / Evidence (if applicable)

<!-- Paste pipeline output, test results, or screenshots here -->
