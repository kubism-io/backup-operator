/*
Copyright 2020 Backup Operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	backupv1alpha1 "github.com/kubism/backup-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type UpdateMongoDBBackupPlanFunc = func(spec *backupv1alpha1.MongoDBBackupPlanSpec)

func newMongoDBBackupPlan(namespace string, updates ...UpdateMongoDBBackupPlanFunc) *backupv1alpha1.MongoDBBackupPlan {
	plan := &backupv1alpha1.MongoDBBackupPlan{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      newTestName(),
		},
		Spec: backupv1alpha1.MongoDBBackupPlanSpec{
			Schedule:              "* * * * *",
			ActiveDeadlineSeconds: 3600,
			Retention:             2,
			URI:                   "mongodb://localhost:27017",
			Destination: &backupv1alpha1.Destination{
				S3: &backupv1alpha1.S3{
					Endpoint:        "localhost:8000",
					Bucket:          "test",
					UseSSL:          false,
					AccessKeyID:     "A",
					SecretAccessKey: "B",
				},
			},
		},
	}
	for _, f := range updates {
		f(&plan.Spec)
	}
	return plan
}

func mustCreateNewMongoDBBackupPlan(namespace string, updates ...UpdateMongoDBBackupPlanFunc) *backupv1alpha1.MongoDBBackupPlan {
	plan := newMongoDBBackupPlan(namespace, updates...)
	Expect(k8sClient.Create(context.Background(), plan)).Should(Succeed())
	return plan
}

var _ = Describe("VaultSecretReconciler", func() {
	ctx := context.Background()
	namespace := ""

	BeforeEach(func() {
		namespace = mustCreateNamespace()
	})
	AfterEach(func() {
		mustDeleteNamespace(namespace)
	})

	It("can create MongoDBBackupPlans", func() {
		Context("with missing data", func() {
			Expect(k8sClient.Create(ctx, &backupv1alpha1.MongoDBBackupPlan{})).ShouldNot(Succeed())
		})
		Context("with valid data", func() {
			plan := mustCreateNewMongoDBBackupPlan(namespace)
			defer mustRemoveFinalizers(plan)
		})
	})
	It("can process MongoDBBackupPlans", func() {
		Context("which are just created", func() {
			plan := mustCreateNewMongoDBBackupPlan(namespace)
			defer mustRemoveFinalizers(plan)
			res := mustReconcile(plan)
			Expect(res.Requeue).To(Equal(false))
		})
		Context("which were deleted", func() {
			plan := mustCreateNewMongoDBBackupPlan(namespace)
			defer func() {
				// If this test fails, we need to make sure the finalizers are removed
				if err := k8sClient.Get(ctx, namespacedName(plan), plan); err == nil {
					mustRemoveFinalizers(plan)
				}
			}()
			res := mustReconcile(plan)
			Expect(res.Requeue).To(Equal(false))
			Expect(k8sClient.Delete(ctx, plan)).Should(Succeed())
			res = mustReconcile(plan)
			Expect(res.Requeue).To(Equal(false))
		})
	})

})