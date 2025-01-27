package work_postgres

import (
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var PgCmdTest2 = &cobra.Command{
	Use:   "pgTestAll",
	Short: "This is a simple CLI application",
	Long:  `Test postgres all`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("Testing 'TestPostgresServerStatus()':\n")
		TestPostgresServerStatus()
		klog.Info("\n========================\n")

		klog.Info("Testing 'TestClientFuncs()':\n")
		TestClientFuncs()
		klog.Info("\n========================\n")

		klog.Info("Testing 'TestCheckAvailableSharedBuffers()':\n")
		TestCheckAvailableSharedBuffers()
		klog.Info("\n========================\n")

		klog.Info("Testing 'TestCheckEffectiveCacheSize()':\n")
		TestCheckEffectiveCacheSize()
		klog.Info("\n========================\n")
	},
}

var PgCmdTestSharedBuffers = &cobra.Command{
	Use:   "pgTestSharedBuffers",
	Short: "This is a simple CLI Application",
	Long:  `Test postgres shared buffers`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("Testing `TestSharedBuffers()`:\n")
		TestCheckAvailableSharedBuffers()
	},
}

var PgCmdTestRequestMethods = &cobra.Command{
	Use:   "pgTestRequestMethods",
	Short: "This is a simple CLI Application",
	Long:  `Test postgres request methods`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("Testing `TestSharedBuffers()`:\n")
		TestCheckRequestMethods()
	},
}
