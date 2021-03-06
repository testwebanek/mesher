package egress_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-mesh/mesher/cmd"
	mesherconfig "github.com/go-mesh/mesher/config"
	_ "github.com/go-mesh/mesher/control/istio"
	"github.com/go-mesh/mesher/pkg/egress"
	"github.com/go-mesh/mesher/pkg/egress/archaius"
	"gopkg.in/yaml.v2"
)

func BenchmarkMatch(b *testing.B) {
	lager.Initialize("", "DEBUG", "",
		"size", true, 1, 10, 7)

	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-mesh/mesher")
	cmd.Init()
	config.Init()
	mesherconfig.Init()
	egress.Init()
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	control.Init(opts)
	var yamlContent = `---
egress:
  infra: cse # pilot or cse
  address: http://istio-pilot.istio-system:15010
egressRule:
  google-ext:
    - hosts:
        - "www.google.com"
        - "*.yahoo.com"
      ports:
        - port: 80
          protocol: HTTP
  facebook-ext:
    - hosts:
        - "www.facebook.com"
      ports:
        - port: 80
          protocol: HTTP`

	ss := mesherconfig.EgressConfig{}
	err := yaml.Unmarshal([]byte(yamlContent), &ss)
	if err != nil {
		fmt.Println("unmarshal failed")
	}
	archaius.SetEgressRule(ss.Destinations)

	myString := "www.google.com"
	for i := 0; i < b.N; i++ {
		egress.Match(myString)
	}
}

func TestMatch(t *testing.T) {
	lager.Initialize("", "DEBUG", "",
		"size", true, 1, 10, 7)

	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-mesh/mesher")
	cmd.Init()
	config.Init()
	mesherconfig.Init()
	egress.Init()
	//control.Init()
	var yamlContent = `---
egress:
  infra: cse # pilot or cse
  address: http://istio-pilot.istio-system:15010
egressRule:
  google-ext:
    - hosts:
        - "www.google.com"
        - "*.yahoo.com"
      ports:
        - port: 80
          protocol: HTTP
  facebook-ext:
    - hosts:
        - "www.facebook.com"
      ports:
        - port: 80
          protocol: HTTP`

	ss := mesherconfig.EgressConfig{}
	err := yaml.Unmarshal([]byte(yamlContent), &ss)
	if err != nil {
		fmt.Println("unmarshal failed")
	}
	archaius.SetEgressRule(ss.Destinations)

	myString := "www.google.com"
	b, _ := egress.Match(myString)
	if b == false {
		t.Errorf("Expected true but got false")
	}
	myString = "*.yahoo.com"
	b, _ = egress.Match(myString)
	if b == false {
		t.Errorf("Expected true but got false")
	}
}
