.. _ksctl_delete-cluster_aws:

ksctl delete-cluster aws
------------------------

Use to deletes a EKS cluster

Synopsis
~~~~~~~~


It is used to delete cluster of given provider

::

  ksctl delete-cluster aws [flags]

Examples
~~~~~~~~

::


  ksctl delete aws --name demo --region ap-south-1 --storage store-local


Options
~~~~~~~

::

  -h, --help             help for aws
  -n, --name string      Cluster Name (default "demo")
  -r, --region string    Region
  -s, --storage string   storage provider
  -v, --verbose int      for verbose output
  -y, --yes              approval to avoid showMsg (default true)

SEE ALSO
~~~~~~~~

* `ksctl delete-cluster <ksctl_delete-cluster.rst>`_ 	 - Use to delete a cluster

*Auto generated by spf13/cobra on 2-Aug-2024*
