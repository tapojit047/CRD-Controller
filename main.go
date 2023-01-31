package main

import (
	"fmt"
	myv1 "github.com/tapojit047/CRD-Controller/pkg/apis/fullmetal.com/v1"
	_ "k8s.io/code-generator"
)

func main() {
	a := myv1.Alchemist{}
	fmt.Println(a)
}
