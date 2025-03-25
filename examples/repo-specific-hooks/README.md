Repository Specific Configuration
===

In this example, we still use the globaly defined hooks in Git that reference the [centralised-hooks](../centralised-hooks/README.md) example.

But we want to extend the centralised hooks to make project specific tasks.

## Reference Central Configuration

Install a minimalistic configuration with a reference to the central repository

```bash
init -m --reference-repository https://github.com/ylallemant/githook-companion/tree/main/examples/centralised-hooks --reference-path ../githooks
```

## Extend Githooks Functionlity

Let say this repo is a Terraform module and we want to ensure some Terraform standards on commit and push.

We added more dependencies to the configuration:
- terraform: in order to init and validate the module (at push time, because it takes time)
- terraform-docs: to generate information about the module (at pre-commit time)
- tflint: to ensure code quality standards (at pre-commit time)

Once the configuration is in place, you can install the dependencies:

- locally if the tools may have different verions in different projects (used in the example)
```shell
githook-companion init
```
