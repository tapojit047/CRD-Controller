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

## Controller:
- `controller.go` file contains the controller structure and all the necessary function
- Controller can be run by running `CRD-Controller` file
### Resources: 
  -  [Sample Controller](https://github.com/kubernetes/sample-controller/blob/master/controller.go#L353)
  - [Vivek Singh Controller](https://www.youtube.com/watch?v=lzoWSfvE2yA&list=PLh4KH3LtJvRQ43JAwwjvTnsVOMp0WKnJO)
  - [Vivek Singh Operator](https://www.youtube.com/watch?v=89PdRvRUcPU&list=PLh4KH3LtJvRTtFWz1WGlyDa7cKjj2Sns0)
  - [Blog](https://www.youtube.com/watch?v=89PdRvRUcPU&list=PLh4KH3LtJvRTtFWz1WGlyDa7cKjj2Sns0)

