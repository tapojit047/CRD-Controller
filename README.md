# CRD-Controller
## Steps of generating code:
- Created `types.go` defining the struct for CRD
- Create `register.go` to register the custom resource
- Don't forget to implicitly add `code-generator`. For that add this line in the ``import`` section of ``main.go``:
`_ "k8s.io/code-generator"`
- Then run `./hack/update-codegen.sh`.
- Oh! To solve the issue run this `chmod u+x vendor/k8s.io/code-generator/generate-groups.sh` and `chmod u+x hack/update-codegen.sh` and then run ``update-codegen.sh``
### Generate manifest(artifacts) yamls using the following command:
- `controller-gen rbac:roleName=controller-perms crd paths=./... output:crd:dir=./artifacts output:stdout`
