# Kubernetes condition pretty printer
A command line tool to print status.conditions that is used frequently in Kubernetes object.

## Installation
```
go get github.com/jsafrane/condition-printer/cmd/cnd
```

## Usage
**Single object:**
```
$ kubeclt get pod mypod -o yaml | cnd

Last transition time: 2020-08-11 12:19:12 +0200 CEST

ContainersReady | True | 
Initialized     | True | 
PodScheduled    | True | 
Ready           | True | 
```
Both `-o json` and `-o yaml` work. And any object that has `.status.conditions` can be used. The output is nicely wrapped and colorful if stdout is a terminal.

**Watch an object:**
```
$ kubectl get pod mypod -w -o json | cnd

Last transition time: 2020-08-11 12:19:12 +0200 CEST

ContainersReady | False | 
Initialized     | True | 
PodScheduled    | True | 
Ready           | False | 

--- [ new screen ] --- 

Last transition time: 2020-08-11 12:19:21 +0200 CEST

ContainersReady | True | 
Initialized     | True | 
PodScheduled    | True | 
Ready           | True | 
```

Use `shift+PageUp` / `shift+PageDown` in a Linux terminal to move between each state.

`-w -o json` is strongly recommended to get the screen up to date. You may get late refreshes with `-w -o yaml`, as it is not possible to find end of an instance ("`---`") before the next instance is printed in YAML output.

**Works great when watching OpenShift operators:**

```
$ kubectl get kubeapiserver -o json | cnd

Last transition time: 2020-08-13 16:34:46 +0200 CEST

BackingResourceControllerDegraded                      | False | 
CertRotationTimeUpgradeable                            | True  | 
CertRotation_AggregatorProxyClientCert_Degraded        | False | 
CertRotation_ExternalLoadBalancerServing_Degraded      | False | 
CertRotation_InternalLoadBalancerServing_Degraded      | False | 
CertRotation_KubeAPIServerToKubeletClientCert_Degraded | False | 
CertRotation_KubeControllerManagerClient_Degraded      | False | 
CertRotation_KubeSchedulerClient_Degraded              | False | 
CertRotation_LocalhostRecoveryServing_Degraded         | False | 
CertRotation_LocalhostServing_Degraded                 | False | 
CertRotation_ServiceNetworkServing_Degraded            | False | 
ConfigObservationDegraded                              | False | 
Encrypted                                              | False | Encryption is not enabled
EncryptionKeyControllerDegraded                        | False | 
EncryptionMigrationControllerDegraded                  | False | 
EncryptionMigrationControllerProgressing               | False | 
EncryptionPruneControllerDegraded                      | False | 
EncryptionStateControllerDegraded                      | False | 
FeatureGatesUpgradeable                                | True  | 
InstallerControllerDegraded                            | False | 
InstallerPodContainerWaitingDegraded                   | False | 
InstallerPodNetworkingDegraded                         | False | 
InstallerPodPendingDegraded                            | False | 
KubeAPIServerStaticResourcesDegraded                   | False | 
NodeControllerDegraded                                 | True  | The master nodes not ready: node "jsafrane-master-1" not ready since 202
                                                                 0-08-12 20:17:43 +0000 UTC because NodeStatusUnknown (Kubelet stopped po
                                                                 sting node status.)
NodeInstallerDegraded                                  | False | 
NodeInstallerProgressing                               | False | 3 nodes are at revision 8
ResourceSyncControllerDegraded                         | False | 
RevisionControllerDegraded                             | False | 
StaticPodsAvailable                                    | True  | 3 nodes are active; 3 nodes are at revision 8
StaticPodsDegraded                                     | False | 
TargetConfigControllerDegraded                         | False | 
UnsupportedConfigOverridesUpgradeable                  | True  | 
```
