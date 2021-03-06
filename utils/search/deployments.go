package search

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	client "github.com/guessi/kubectl-search/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	deploymentsFields = "NAMESPACE\tNAME\tDESIRED\tCURRENT\tUP-TO-DATE\tAVAILABLE\tAGE"
)

// Deployments - a public function for searching pods with keyword
func Deployments(namespace string, allNamespaces bool, selector, fieldSelector, keyword string) {
	clientset := client.InitClient()

	if len(namespace) <= 0 {
		namespace = "default"
	}

	if allNamespaces {
		namespace = ""
	}

	listOptions := &metav1.ListOptions{
		LabelSelector: selector,
		FieldSelector: fieldSelector,
	}

	deployments, err := clientset.AppsV1().Deployments(namespace).List(*listOptions)
	if err != nil {
		panic(err.Error())
	}

	buf := bytes.NewBuffer(nil)
	w := tabwriter.NewWriter(buf, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, deploymentsFields)
	for _, d := range deployments.Items {
		// return all deployments under namespace if no keyword specific
		if len(keyword) > 0 {
			match := strings.Contains(d.Name, keyword)
			if !match {
				continue
			}
		}

		age, ageUnit := getAge(time.Since(d.CreationTimestamp.Time).Seconds())

		dInfo := fmt.Sprintf("%s\t%s\t%d\t%d\t%d\t%d\t%d%s",
			d.Namespace,
			d.Name,
			d.Status.Replicas,
			d.Status.ReadyReplicas,
			d.Status.UpdatedReplicas,
			d.Status.AvailableReplicas,
			age, ageUnit,
		)
		fmt.Fprintln(w, dInfo)
	}
	w.Flush()

	fmt.Printf("%s", buf.String())
}
