# Contributing to configsManager

We ❤️ contributions! To keep the project clean, consistent, and maintainable, please follow these guidelines before creating an issue or submitting a PR.

---

## 1. Before creating an Issue

Please check if a similar issue already exists. If not, create a new one using this template:

### 📄 Summary
Provide a short 1-2 line description of the problem or proposal.

### 🧩 Problem / Motivation
Explain what’s wrong, why this issue is important, and why it should be solved.

### 💡 Proposed Solution
Share your idea for a solution—even if you’re not 100% sure.

### 📎 Additional Context
Include screenshots, logs, references, commits, or code examples to help explain the issue.

---

## 2. Before submitting a Pull Request (PR)

1. **Always link the PR to an issue.** If no issue exists, make sure the PR is fully self-contained and clearly describes the change.  
2. **Run all tests**:  

```bash
go test ./... --count=1
````

3. **Check formatting**:

```bash
gofmt -s -w .
```

4. **Run the linter**:

```bash
golangci-lint run
```

> Note: All of the above checks are automatically run in CI, so the PR will fail if any tests, linter checks, or formatting issues are found.

Use this template for your PR description:

### 📄 Summary

Briefly describe the problem or improvement.

### 🧩 What I’ve done

Explain what the PR actually implements or fixes.

### 💡 Why this approach

Explain why you chose this solution and what are its advantages—even if you’re not 100% sure.

### 📎 Additional

Screenshots, logs, links, commits, or relevant code snippets that help review your PR.

---

## 3. Code style and workflow

* Follow the Go idioms and project formatting (`gofmt -s -w`).
* Keep your commits small and focused.
* Make sure your branch is up-to-date with `master` before creating a PR.
* Ensure CI passes: tests, linter, and formatting checks.

By following these guidelines, you help keep `configsManager` clean, maintainable, and easy to review. Thank you for contributing! 🚀

