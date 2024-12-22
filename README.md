# Kaniko Kubernetes operator
### Kubernetes operator for building docker image using kaniko as a tool for creating container images from Dockerfile

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
#   command:
#     - /bin/sh
  args: 
    - --cache=true
    - --cache-repo=registry.home.local/kaniko-cache
    - --skip-tls-verify
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
###


https://github.com/GoogleContainerTools/kaniko