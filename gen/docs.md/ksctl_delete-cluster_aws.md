## ksctl delete-cluster aws

Use to deletes a EKS cluster

### Synopsis

[105;90m[0;0m
[104;90m░  ░░░░  ░░░      ░░░░      ░░░        ░░  ░░░░░░░[0;0m
[106;90m▒  ▒▒▒  ▒▒▒  ▒▒▒▒▒▒▒▒  ▒▒▒▒  ▒▒▒▒▒  ▒▒▒▒▒  ▒▒▒▒▒▒▒[0;0m
[102;90m▓     ▓▓▓▓▓▓      ▓▓▓  ▓▓▓▓▓▓▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓▓▓[0;0m
[103;90m▓  ▓▓▓  ▓▓▓▓▓▓▓▓▓  ▓▓  ▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓  ▓▓▓▓▓▓▓[0;0m
[101;90m█  ████  ███      ████      ██████  █████        █[0;0m

[103;30mIt is used to delete cluster of given provider[0;0m

```
ksctl delete-cluster aws [flags]
```

### Examples

```

ksctl delete aws --name demo --region ap-south-1 --storage store-local

```

### Options

```
  -h, --help             help for aws
  -n, --name string      Cluster Name (default "demo")
  -r, --region string    Region
  -s, --storage string   storage provider
  -v, --verbose int      for verbose output
  -y, --yes              approval to avoid showMsg (default true)
```

### SEE ALSO

* [ksctl delete-cluster](ksctl_delete-cluster.md)	 - Use to delete a cluster

###### Auto generated by spf13/cobra on 24-Aug-2024
