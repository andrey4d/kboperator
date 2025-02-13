# Build in kubernetes
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
#### Run operator 
```shell
kubectl run -it --rm --image=fedora -- bash
[root@bash /]# dnf install git make golang kubectl
[root@bash /]# git clone https://github.com/andrey4d/kboperator.git
[root@bash /]#
[root@bash kboperator]# make genereate
[root@bash kboperator]# make manifests
[root@bash kboperator]# make install
[root@bash kboperator]# make run
```

cat <<EOF | kubectl apply -f -
apiVersion: kbo.k8s.dav.io/v1alpha1
kind: KanikoBuild
metadata:
  labels:
    app.kubernetes.io/name: builder-operator
    app.kubernetes.io/managed-by: kustomize
  name: kbopertor
  namespace: kaniko-builder
spec:
  name: kbo-builder
  image: gcr.io/kaniko-project/executor:debug
  command:
    - /kaniko/executor
  args:
    - --context=git://github.com/andrey4d/kboperator.git#refs/heads/main
    - --dockerfile=/kaniko/buildcontext/Dockerfile
    - --destination=registry-docker-registry.registry:5000/kboperator:v0.0.6
    - --insecure-registry=registry-docker-registry.registry:5000
  docker_config:
    registry: registry-docker-registry.registry:5000
    auth: dXNlcjpwYXNzd29yZA==
EOF

make deploy IMG=registry.dvampere.lab/kboperator:v0.0.6
kubectl delete ClusterRoleBinding default-manager-rolebinding
```