# ephemerup

A Helm chart for ephemerup

## Source Code

* <https://github.com/tlinden/ephemerup>

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://charts.bitnami.com/bitnami | common | 1.x.x |

## Values

| Key                                               | Type   | Description                                               | Default               |
|---------------------------------------------------|--------|-----------------------------------------------------------|-----------------------|
| kubeVersion                                       | string |                                                           | `""`                  |
| nameOverride                                      | string |                                                           | `""`                  |
| fullnameOverride                                  | string |                                                           | `""`                  |
| namespaceOverride                                 | string |                                                           | `""`                  |
| commonLabels                                      | object |                                                           | `{}`                  |
| commonAnnotations.app                             | string |                                                           | `"ephemerup"`         |
| clusterDomain                                     | string |                                                           | `"cluster.local"`     |
| logLevel                                          | string |                                                           | `"info"`              |
| image.registry                                    | string |                                                           | `"docker.io"`         |
| image.repository                                  | string |                                                           | `"tlinden/ephemerup"` |
| image.tag                                         | string |                                                           | `"latest"`            |
| image.pullPolicy                                  | string |                                                           | `"IfNotPresent"`      |
| image.pullSecrets                                 | list   |                                                           | `[]`                  |
| secrets                                           | object |                                                           | `{}`                  |
| mountSecrets                                      | list   |                                                           | `[]`                  |
| env                                               | list   |                                                           | `[]`                  |
| config                                            | object | Backup plans. For details, see [values.yaml](values.yaml) | `{}`                  |
| replicaCount                                      | int    |                                                           | `1`                   |
| sidecars                                          | list   |                                                           | `[]`                  |
| lifecycleHooks                                    | object |                                                           | `{}`                  |
| podAnnotations                                    | object |                                                           | `{}`                  |
| podLabels                                         | object |                                                           | `{}`                  |
| updateStrategy.type                               | string |                                                           | `"RollingUpdate"`     |
| podAffinityPreset                                 | string |                                                           | `""`                  |
| podAntiAffinityPreset                             | string |                                                           | `"soft"`              |
| nodeAffinityPreset.type                           | string |                                                           | `""`                  |
| nodeAffinityPreset.key                            | string |                                                           | `""`                  |
| nodeAffinityPreset.values                         | list   |                                                           | `[]`                  |
| affinity                                          | object |                                                           | `{}`                  |
| nodeSelector                                      | object |                                                           | `{}`                  |
| tolerations                                       | list   |                                                           | `[]`                  |
| resources.limits.cpu                              | string |                                                           | `"500m"`              |
| resources.limits.memory                           | string |                                                           | `"256Mi"`             |
| resources.requests.cpu                            | string |                                                           | `"100m"`              |
| resources.requests.memory                         | string |                                                           | `"128Mi"`             |
| podSecurityContext.fsGroup                        | int    |                                                           | `65534`               |
| containerSecurityContext.enabled                  | bool   |                                                           | `false`               |
| containerSecurityContext.allowPrivilegeEscalation | bool   |                                                           | `false`               |
| containerSecurityContext.capabilities.drop[0]     | string |                                                           | `"ALL"`               |
| containerSecurityContext.privileged               | bool   |                                                           | `false`               |
| containerSecurityContext.runAsUser                | int    |                                                           | `0`                   |
| containerSecurityContext.runAsNonRoot             | bool   |                                                           | `false`               |
| livenessProbe.enabled                             | bool   |                                                           | `true`                |
| livenessProbe.initialDelaySeconds                 | int    |                                                           | `5`                   |
| livenessProbe.timeoutSeconds                      | int    |                                                           | `1`                   |
| livenessProbe.periodSeconds                       | int    |                                                           | `20`                  |
| livenessProbe.failureThreshold                    | int    |                                                           | `6`                   |
| livenessProbe.successThreshold                    | int    |                                                           | `1`                   |
| readinessProbe.enabled                            | bool   |                                                           | `true`                |
| readinessProbe.initialDelaySeconds                | int    |                                                           | `5`                   |
| readinessProbe.timeoutSeconds                     | int    |                                                           | `1`                   |
| readinessProbe.periodSeconds                      | int    |                                                           | `20`                  |
| readinessProbe.failureThreshold                   | int    |                                                           | `6`                   |
| readinessProbe.successThreshold                   | int    |                                                           | `1`                   |
| startupProbe.enabled                              | bool   |                                                           | `true`                |
| startupProbe.initialDelaySeconds                  | int    |                                                           | `10`                  |
| startupProbe.timeoutSeconds                       | int    |                                                           | `1`                   |
| startupProbe.periodSeconds                        | int    |                                                           | `20`                  |
| startupProbe.failureThreshold                     | int    |                                                           | `6`                   |
| startupProbe.successThreshold                     | int    |                                                           | `1`                   |
| customLivenessProbe                               | object |                                                           | `{}`                  |
| customStartupProbe                                | object |                                                           | `{}`                  |
| customReadinessProbe                              | object |                                                           | `{}`                  |
| service.type                                      | string |                                                           | `"ClusterIP"`         |
| service.ports.http                                | int    |                                                           | `8090`                |
| service.nodePorts.http                            | string |                                                           | `""`                  |
| service.clusterIP                                 | string |                                                           | `""`                  |
| service.extraPorts                                | list   |                                                           | `[]`                  |
| service.loadBalancerIP                            | string |                                                           | `""`                  |
| service.loadBalancerSourceRanges                  | list   |                                                           | `[]`                  |
| service.externalTrafficPolicy                     | string |                                                           | `"Cluster"`           |
| service.annotations                               | object |                                                           | `{}`                  |
| service.sessionAffinity                           | string |                                                           | `"None"`              |
| service.sessionAffinityConfig                     | object |                                                           | `{}`                  |
| ingress.enabled                                   | bool   |                                                           | `false`               |
| ingress.pathType                                  | string |                                                           | `"Prefix"`            |
| ingress.apiVersion                                | string |                                                           | `""`                  |
| ingress.hostname                                  | string |                                                           | `"ephemerup.local"`   |
| ingress.path                                      | string |                                                           | `"/"`                 |
| ingress.annotations                               | object |                                                           | `{}`                  |
| ingress.tls                                       | bool   |                                                           | `false`               |
| ingress.tlsSecretName                             | string |                                                           | `""`                  |
| ingress.extraPaths                                | list   |                                                           | `[]`                  |
| ingress.selfSigned                                | bool   |                                                           | `false`               |
| ingress.ingressClassName                          | string |                                                           | `"nginx"`             |
| ingress.extraHosts                                | list   |                                                           | `[]`                  |
| ingress.extraTls                                  | list   |                                                           | `[]`                  |
| ingress.secrets                                   | list   |                                                           | `[]`                  |
| ingress.extraRules                                | list   |                                                           | `[]`                  |
metrics.serviceMonitor.enabled | bool | `true` |  |
| metrics.serviceMonitor.port | string | `"http"` |  |
| metrics.serviceMonitor.namespace | string | `""` |  |
| metrics.serviceMonitor.interval | string | `"30s"` |  |
| metrics.serviceMonitor.scrapeTimeout | string | `"10s"` |  |
| storage.longTerm | object | `{"name":"ephemerup-storage","spec":{"accessModes":["ReadWriteOnce"],"resources":{"requests":{"storage":"100Gi"}},"storageClassName":"standard"}}` | Persistent volume for backups, see `config.retention` |
| storage.tmp | object | `{"name":"ephemerup-tmp","spec":{"accessModes":["ReadWriteOnce"],"resources":{"requests":{"storage":"10Gi"}},"storageClassName":"standard"}}` | Persistent volume for temporary files |
| storage.restoreTmp.name | string | `"ephemerup-restore-tmp"`  |  |
| storage.restoreTmp.spec.accessModes[0] | string | `"ReadWriteOnce"` |  |
| storage.restoreTmp.spec.resources.requests.storage | string | `"100Gi"` |  |
| storage.restoreTmp.spec.storageClassName | string | `"standard"` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
