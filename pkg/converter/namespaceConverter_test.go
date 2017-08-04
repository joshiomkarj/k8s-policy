package converter_test

import (
	"encoding/json"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/projectcalico/k8s-policy/pkg/converter"
	"github.com/projectcalico/libcalico-go/lib/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

/*
1. check for empty object
2. send null object
3. check type of the returned object
4. check if name is correct and in the correct format
5,6. check that namespace rules are "allow" ingress and egress
7. check for correctness of label

*/
var _ = Describe("NamespaceConverter", func() {

	//1. check for empty object
	Context("With empty namespace object", func() {

		ns := v1.Namespace{
			Spec:   v1.NamespaceSpec{},
			Status: v1.NamespaceStatus{},
		}

		nsConverter := converter.NewNamespaceConverter()
		policyObject, _ := nsConverter.Convert(&ns)

		It("policy object must not be null", func() {
			Expect(policyObject).NotTo(Equal(nil))
		})

	})

	/*
		//	2. send null object  -- currently commented due to panic and fatal issue
			Context("With null namespace object", func() {

					var ns v1.Namespace
					nsConverter := converter.NewNamespaceConverter()
			//		policyObject,errorObject := nsConverter.Convert(ns)


					It("policy object must throw an error", func() {

						Ω(nsConverter.Convert(ns)).Should(Panic())
			//			Expect(policyObject).To(Equal(nil))
			//			Expect(errorObject).NotTo(Equal(nil))
					})
			        })

	*/

	//3. check type of the returned object
	Context("With the namespace object", func() {

		ns := v1.Namespace{
			Spec:   v1.NamespaceSpec{},
			Status: v1.NamespaceStatus{},
		}

		nsConverter := converter.NewNamespaceConverter()
		policyObject, _ := nsConverter.Convert(&ns)

		x := reflect.TypeOf(*api.NewProfile())
		y := reflect.TypeOf(policyObject)

		It("the returned object must be policy object", func() {
			Ω(x).Should(Equal(y))
		})
		//		fmt.Printf("%#v", policyObject)
	})

	//4. check if name is correct and in the correct format
	Context("With assigned names", func() {
		testString := "testObjectName"
		ns := v1.Namespace{
			Spec:       v1.NamespaceSpec{},
			Status:     v1.NamespaceStatus{},
			ObjectMeta: metav1.ObjectMeta{Name: testString},
		}

		nsConverter := converter.NewNamespaceConverter()
		temp, _ := nsConverter.Convert(&ns)
		policyObject := temp.(api.Profile)

		checkName := "ns.projectcalico.org/" + testString
		receivedName := policyObject.Metadata.Name

		It("the returned object must have the same names", func() {
			Ω(checkName).Should(Equal(receivedName))
		})
		fmt.Println(policyObject.Metadata.Name)

	})

	//5,6. check that namespace rules are "allow" ingress and egress
	Context("With the namespace object, check ingress rules in policy", func() {

		ns := v1.Namespace{
			Spec:   v1.NamespaceSpec{},
			Status: v1.NamespaceStatus{},
		}

		nsConverter := converter.NewNamespaceConverter()
		temp, _ := nsConverter.Convert(&ns)
		policyObject := temp.(api.Profile)

		ingressRuleReturned := policyObject.Spec.IngressRules[0].Action
		ingressRuleExpected := "allow"

		It("must be allowed", func() {
			Ω(ingressRuleExpected).Should(Equal(ingressRuleReturned))
		})
	})

	Context("With the namespace object, check egress rules in policy", func() {

		ns := v1.Namespace{
			Spec:   v1.NamespaceSpec{},
			Status: v1.NamespaceStatus{},
		}

		nsConverter := converter.NewNamespaceConverter()
		temp, _ := nsConverter.Convert(&ns)
		policyObject := temp.(api.Profile)

		egressRuleReturned := policyObject.Spec.EgressRules[0].Action
		egressRuleExpected := "allow"

		It("must be allowed", func() {
			Ω(egressRuleExpected).Should(Equal(egressRuleReturned))
		})
	})

	//7. check for correctness of label
	Context("With the namespace object, check labels", func() {
		var labelsMap map[string]string
		labelsMap = make(map[string]string)

		labelsMap["app"] = "nginx"

		ns := v1.Namespace{
			Spec:       v1.NamespaceSpec{},
			Status:     v1.NamespaceStatus{},
			ObjectMeta: metav1.ObjectMeta{Labels: labelsMap},
		}

		nsConverter := converter.NewNamespaceConverter()
		temp, _ := nsConverter.Convert(&ns)
		policyObject := temp.(api.Profile)

		var receivedLabel map[string]string

		receivedLabel = policyObject.Metadata.Labels

		It("labels must match", func() {
			Ω(labelsMap["app"]).Should(Equal(receivedLabel["k8s_ns/label/app"]))
		})

		a, _ := json.Marshal(policyObject)
		fmt.Println(string(a))
	})

})
