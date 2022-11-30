# Contributing to Jupiter

Welcome to Jupiter! This document is a guideline about how to contribute to Jupiter.
If you find something incorrect or missing, please leave comments / suggestions.

## Before you get started

### Setting up your development environment

You should have Go 1.13+ installed in your operating system.

## Contributing

We are always very happy to have contributions, whether for typo fix, bug fix or big new features.
Please do not ever hesitate to ask a question or send a pull request.

We strongly value documentation and integration with other projects.
We are very glad to accept improvements for these aspects.

### GitHub workflow

We use the `master` branch as the development branch, which indicates that this is an unstable branch.

Here are the workflow for contributors:

1. Fork to your own
1. Clone fork to local repository
1. Create a new branch and work on it
1. Keep your branch in sync
1. Commit your changes (make sure your commit message concise)
1. Push your commits to your forked repository
1. Create a pull request

Please follow [the pull request template](./.github/PULL_REQUEST_TEMPLATE.md).
Please make sure the PR has a corresponding issue.

After creating a PR, one or more reviewers will be assigned to the pull request.
The reviewers will review the code.

Before merging a PR, squash any fix review feedback, typo, merged, and rebased sorts of commits.
The final commit message should be clear and concise.

### Open an issue / PR

We use [GitHub Issues](https://github.com/douyu/jupiter/issues) and [Pull Requests](https://github.com/douyu/jupiter/pulls) for trackers.

If you find a typo in a document, find a bug in code, or want new features, or want to give suggestions,
you can [open an issue on GitHub](https://github.com/douyu/jupiter/issues/new) to report it.
Please follow the guideline message in the issue template.

If you want to contribute, please follow the [contribution workflow](#github-workflow) and create a new pull request.
If your PR contains large changes, e.g. component refactor or new components, please write detailed documents
about its design and usage.

Note that a single pull request should not be too large. If heavy changes are required, it's better to separate the changes
to a few individual PRs.

### Code review

All code should be well reviewed by one or more committers. Some principles:

- Readability: Important code should be well-documented. Comply with our code style.
- Elegance: New functions, classes or components should be well-designed.
- Testability: Important code should be well-tested (high unit test coverage).

## Community

### Contact us
