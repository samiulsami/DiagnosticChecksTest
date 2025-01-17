package work_postgres

import (
	"context"
	"fmt"

	_ "database/sql"

	_ "github.com/lib/pq"
	utils "github.com/shn27/Test/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	kmapi "kmodules.xyz/client-go/api/v1"
	"kubedb.dev/apimachinery/apis/kubedb"
	dbapi "kubedb.dev/apimachinery/apis/kubedb/v1"
	"kubedb.dev/db-client-go/postgres"
)

const (
	postgresResourceMemoryKey = "memory"
)

func PrimaryServiceDNS(db *dbapi.Postgres) string {
	return fmt.Sprintf("%v.%v.svc", db.ServiceName(), db.Namespace)
}

func GetPostgresClient(kbClient client.Client, db *dbapi.Postgres) (*postgres.Client, error) {
	kubeDBClient, err := postgres.NewKubeDBClientBuilder(kbClient, db).
		WithContext(context.Background()).
		WithURL(PrimaryServiceDNS(db)).
		GetPostgresClient()
	if err != nil {
		return nil, fmt.Errorf("failed to build kubedb postgres client : %v", err)
	}

	return kubeDBClient, nil
}

func GetPostgresDB(kbClient client.Client) (*dbapi.Postgres, error) {
	ref := kmapi.ObjectReference{
		Name:      "postgres",
		Namespace: "demo",
	}
	gvk := schema.GroupVersionKind{
		Version: "v1",
		Group:   "kubedb.com",
		Kind:    "Postgres",
	}

	obj, err := utils.GetK8sObject(gvk, ref, kbClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s object : %v", err)
	}

	db := &dbapi.Postgres{}
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(obj.UnstructuredContent(), db)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unstructured object to a concrete type: %w", err)
	}

	return db, nil
}

func GetTotalMemory(postgresClient *postgres.Client, db *dbapi.Postgres) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("db is nil")
	}

	totalMemory := int64(0)
	var pgContainer *corev1.Container
	for _, v := range db.Spec.PodTemplate.Spec.Containers {
		if v.Name == kubedb.PostgresContainerName {
			pgContainer = &v
			break
		}
	}

	if pgContainer == nil {
		return 0, fmt.Errorf("postgres container not found")
	}

	if qv, exists := pgContainer.Resources.Requests[postgresResourceMemoryKey]; exists {
		totalMemory += int64(qv.Value())
	} else if qv, exists := pgContainer.Resources.Limits[postgresResourceMemoryKey]; exists {
		totalMemory += int64(qv.Value())
	} else {
		return 0, fmt.Errorf("neither memory request nor limit found")
	}

	return totalMemory, nil
}

func GetSharedBuffers(postgresClient *postgres.Client) (string, error) {
	var sharedBuffers string
	if err := postgresClient.DB.QueryRow("SHOW shared_buffers").Scan(&sharedBuffers); err != nil {
		return "", fmt.Errorf("failed to get shared buffers: %w", err)
	}
	return sharedBuffers, nil
}

func GetEffectiveCacheSize(postgresClient *postgres.Client) (string, error) {
	var effectiveCacheSize string
	if err := postgresClient.DB.QueryRow("SHOW effective_cache_size").Scan(&effectiveCacheSize); err != nil {
		return "", fmt.Errorf("failed to get effective cache size: %w", err)
	}
	return effectiveCacheSize, nil
}
