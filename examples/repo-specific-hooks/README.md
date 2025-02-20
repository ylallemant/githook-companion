Repository Specific Configuration
===

In this example, we still use the globaly defined hooks in Git that reference the [centralised-hooks](../centralised-hooks//README.md) example.

But we want to extend the centralised hooks to make project specific tasks.

Let say this repo is a Terraform module and we want to ensure some Terraform standards on commit and push.

We added more dependencies to the configuration:
- terraform: in order to init and validate the module (at push time, because it takes time)
- terraform-docs: to generate information about the module (at pre-commit time)
- tflint: to ensure code quality standards (at pre-commit time)
