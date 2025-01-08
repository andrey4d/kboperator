# Kaniko Kubernetes operator
### Kubernetes operator for building docker image using kaniko as a tool for creating container images.

### kboperator build manifest "build context git://"

```yaml
apiVersion: kbo.k8s.dav.io/v1alpha1
kind: KanikoBuild
metadata:
  labels:
    app.kubernetes.io/name: builder-operator
    app.kubernetes.io/managed-by: kustomize
  name: kbopertor
  namespace: kaniko-test
spec:
  name: builder
  image: registry.home.local/kaniko/executor:v1.23.2-debug
  command:
    - /kaniko/executor
  args:
    - --context=git://github.com/andrey4d/kboperator.git#refs/heads/main
    - --dockerfile=/kaniko/buildcontext/Dockerfile
    - --destination=registry.home.local/kboperator:v0.0.1
    - --cache=true
    - --cache-repo=registry.home.local/kaniko-cache
    - --skip-tls-verify
  docker_config:
    registry: registry.home.local
    auth: dXNlcjpwYXNzd29yZA==
```
```shell
make deploy IMG=registry.home.local/kboperator:v0.0.1
```
### Make installer
```shell
make build-installer IMG=registry.home.local/kboperator:v0.0.1
kubectl apply -f  dist/install.yaml
```

#### manifest
```yaml
apiVersion: kbo.k8s.dav.io/v1alpha1
kind: KanikoBuild
metadata:
  labels:
    app.kubernetes.io/name: builder-operator
    app.kubernetes.io/managed-by: kustomize
  name: kanikobuild-sample
  namespace: kaniko-test
spec:
  name:  builder
  image: gcr.io/kaniko-project/executor:latest
  context: /workspace
  destination: registry.home.local/kboperator-curl:v0.0.1
  args: 
    - --cache=true
    - --cache-repo=registry.home.local/kaniko-cache
    - --skip-tls-verify
  persistence:
    enabled: true
    storageClass: local-path
    volumeSize: 1Gi    
  docker_config:
    registry: registry.home.local
    auth: dXNlcjpwYXNzd29yZA==
  dockerfile: |-
    FROM library/alpine:latest
    RUN apk add --no-cache curl
    CMD ["curl", "http://www.google.com"]
```

#### minimal manifest
```yaml
apiVersion: kbo.k8s.dav.io/v1alpha1
kind: KanikoBuild
metadata:
  name: kanikobuild-sample
  namespace: kaniko-test
spec:
  destination: registry.home.local/kboperator-curl:v0.0.1
  docker_config:
    registry: registry.home.local
    auth: dXNlcjpwYXNzd29yZA==
  dockerfile: |-
    FROM library/alpine:latest
    RUN apk add --no-cache curl
    CMD ["curl", "http://www.google.com"]
```
### build in kubernetes
#### default RBAC
```shell
kubectl apply -f config/rbac/role.yaml

cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  labels:
    app.kubernetes.io/name: kboperator
    app.kubernetes.io/managed-by: kustomize
  name: default-manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: manager-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: default
EOF  
```
```shell
kubectl run -it --rm --image=fedora -- bash
[root@bash /]# dnf install git make golang kubectl
[root@bash /]#git clone https://github.com/andrey4d/kboperator.git
[root@bash kboperator]# make genereate
[root@bash kboperator]# make manifests
[root@bash kboperator]# make install
[root@bash kboperator]# make run
```

https://github.com/GoogleContainerTools/kaniko

