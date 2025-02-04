## How to Add a Configuration

To add a new configuration for a project and environment, create a file in this folder using the format:

```
{project}/{environment}.yml
```

For example, to add a configuration for the `my-project` project in the `production` environment, create a file named:

```bash
my-project/production.yml
```

The file should contain the configuration settings for the project and environment in YAML format. For example:

```yaml
configs:
    key: value
    key2: value2
```

> Note: Nested configurations are not supported.

Example Configuration File: [`sample/development.yml`](sample/development.yml)

