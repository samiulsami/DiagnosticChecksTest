package work_postgres

import (
	"encoding/json"
	"fmt"

	gohumanize "github.com/dustin/go-humanize"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"kubedb.dev/apimachinery/apis/kubedb"
)

func TestPostgresServerStatus() {
	_, _, pgClient, err := GetPostgresClientsAndDB()
	if err != nil {
		klog.Error(err, "failed to get postgres clients and db")
		return
	}

	stats := pgClient.DB.Stats()
	klog.Info("=======Test postgres server stats=======")

	prettyData, err := json.MarshalIndent(stats, "  ", "   ")
	if err != nil {
		klog.Error(err, "failed to marshal db stats")
	}

	klog.Info(string(prettyData))
}

func TestClientFuncs() {
	_, _, pgClient, err := GetPostgresClientsAndDB()
	if err != nil {
		klog.Error(err, "failed to get postgres clients and db")
		return
	}

	err = pgClient.DB.Ping()
	if err != nil {
		klog.Error(err, "failed to ping postgres")
		return
	}

	klog.Info("Pinged postgres\n")
	klog.Infof("pgClient.DB.Stats().InUse : %d", pgClient.DB.Stats().InUse)
}

func TestCheckAvailableSharedBuffers() {
	_, db, pgClient, err := GetPostgresClientsAndDB()
	if err != nil {
		klog.Error(err, "failed to get postgres clients and db")
		return
	}

	totalMemory, err := GetTotalMemory(db)
	if err != nil {
		klog.Error(err, "failed to get total memory")
		return
	}

	sharedBuffersStr, err := GetSharedBuffers(pgClient)
	if err != nil {
		klog.Error(err, "failed to get shared buffers")
		return
	}

	sharedBuffers, err := gohumanize.ParseBytes(sharedBuffersStr)
	if err != nil {
		klog.Error(err, "failed to parse shared buffers")
		return
	}
	klog.Infof("Total memory: %s\n", gohumanize.IBytes(uint64(totalMemory)))
	klog.Infof("Shared buffers: %s\n", gohumanize.IBytes(uint64(sharedBuffers)))

	percentage := float64(sharedBuffers) / float64(totalMemory)
	klog.Infof("Shared buffers percentage: %.2f%%\n", percentage*float64(100))
}

func TestCheckEffectiveCacheSize() {
	_, db, pgClient, err := GetPostgresClientsAndDB()
	if err != nil {
		klog.Error(err, "failed to get postgres clients and db")
		return
	}

	totalMemory, err := GetTotalMemory(db)
	if err != nil {
		klog.Error(err, "failed to get total memory")
		return
	}

	effectiveCacheSizeStr, err := GetEffectiveCacheSize(pgClient)
	if err != nil {
		klog.Error(err, "failed to get shared buffers")
		return
	}

	effectiveCacheSize, err := gohumanize.ParseBytes(effectiveCacheSizeStr)
	if err != nil {
		klog.Error(err, "failed to parse shared buffers")
		return
	}
	klog.Infof("Total memory: %s\n", gohumanize.IBytes(uint64(totalMemory)))
	klog.Infof("effective cache size: %s\n", gohumanize.IBytes(uint64(effectiveCacheSize)))

	percentage := float64(effectiveCacheSize) / float64(totalMemory)
	klog.Infof("effective cache size percentage: %.2f%%\n", percentage*float64(100))
}

func TestCheckRequestMethods() {
	_, db, _, err := GetPostgresClientsAndDB()
	if err != nil {
		klog.Error(err, "failed to get postgres clients and db")
		return
	}

	if db == nil {
		klog.Error("db is nil")
		return
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
	}

	memoryMethod := pgContainer.Resources.Requests.Memory()
	cpuMethod := pgContainer.Resources.Requests.Cpu()

	klog.Infof("MemoryMethod: %v", memoryMethod.Value())
	klog.Infof("CPUMethod: %v", cpuMethod.Value())
}
