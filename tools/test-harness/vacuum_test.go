package tests

import (
	"testing"
)

func TestVacuum(t *testing.T) {
	t.Parallel()
	t.Log("Testing the 'vacuum' example...")

	harness := setup(t, timeout, true)

	t.Log("Starting primary...")
	env := []string{}
	_, err := harness.runExample("examples/kube/primary/run.sh", env, t)
	if err != nil {
		t.Fatalf("Could not run example: %s", err)
	}
	if harness.Cleanup {
		defer harness.Client.DeleteNamespace(harness.Namespace)
		defer harness.runExample("examples/kube/primary/cleanup.sh", env, t)
	}

	pods := []string{"primary"}
	t.Log("Checking if pods are ready to use...")
	if err := harness.Client.CheckPods(harness.Namespace, pods); err != nil {
		t.Fatal(err)
	}

	t.Log("Starting vacuum job...")
	_, err = harness.runExample("examples/kube/vacuum/run.sh", env, t)
	if err != nil {
		t.Fatalf("Could not run example: %s", err)
	}
	if harness.Cleanup {
		defer harness.runExample("examples/kube/vacuum/cleanup.sh", env, t)
	}

        t.Log("Checking if job has completed...")
        job, err := harness.Client.GetJob(harness.Namespace, "vacuum")
        if err != nil {
                t.Fatal(err)
        }

        if err := harness.Client.IsJobComplete(harness.Namespace, job); err != nil {
                t.Fatal(err)
        }

	report, err := harness.createReport()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(report)
}
