package models
import "time"

type CAdvisorDockerList struct {
	Name string `json:"name"`
	Subcontainers []struct {
		Name string `json:"name"`
	} `json:"subcontainers"`
}


type CAdvisorDocker struct {
	Name string `json:"name"`
	Aliases []string `json:"aliases"`
	Namespace string `json:"namespace"`
	Spec struct {
		CreationTime time.Time `json:"creation_time"`
		HasCPU bool `json:"has_cpu"`
		CPU struct {
			Limit int `json:"limit"`
			MaxLimit int `json:"max_limit"`
			Mask string `json:"mask"`
		} `json:"cpu"`
		HasMemory bool `json:"has_memory"`
		Memory struct {
			Limit int64 `json:"limit"`
			SwapLimit int64 `json:"swap_limit"`
		} `json:"memory"`
		HasNetwork bool `json:"has_network"`
		HasFilesystem bool `json:"has_filesystem"`
		HasDiskio bool `json:"has_diskio"`
		HasCustomMetrics bool `json:"has_custom_metrics"`
		Image string `json:"image"`
	} `json:"spec"`
	Stats []struct {
		Timestamp time.Time `json:"timestamp"`
		CPU struct {
			Usage struct {
				Total int64 `json:"total"`
				PerCPUUsage []int `json:"per_cpu_usage"`
				User int64 `json:"user"`
				System int64 `json:"system"`
			} `json:"usage"`
			LoadAverage int `json:"load_average"`
		} `json:"cpu"`
		Diskio struct {
			IoServiceBytes []struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
				Stats struct {
					Async int `json:"Async"`
					Read int `json:"Read"`
					Sync int `json:"Sync"`
					Total int `json:"Total"`
					Write int `json:"Write"`
				} `json:"stats"`
			} `json:"io_service_bytes"`
			IoServiced []struct {
				Major int `json:"major"`
				Minor int `json:"minor"`
				Stats struct {
					Async int `json:"Async"`
					Read int `json:"Read"`
					Sync int `json:"Sync"`
					Total int `json:"Total"`
					Write int `json:"Write"`
				} `json:"stats"`
			} `json:"io_serviced"`
		} `json:"diskio"`
		Memory struct {
			Usage int `json:"usage"`
			WorkingSet int `json:"working_set"`
			ContainerData struct {
				Pgfault int `json:"pgfault"`
				Pgmajfault int `json:"pgmajfault"`
			} `json:"container_data"`
			HierarchicalData struct {
				Pgfault int `json:"pgfault"`
				Pgmajfault int `json:"pgmajfault"`
			} `json:"hierarchical_data"`
		} `json:"memory"`
		Network struct {
			Name string `json:"name"`
			RxBytes int `json:"rx_bytes"`
			RxPackets int `json:"rx_packets"`
			RxErrors int `json:"rx_errors"`
			RxDropped int `json:"rx_dropped"`
			TxBytes int `json:"tx_bytes"`
			TxPackets int `json:"tx_packets"`
			TxErrors int `json:"tx_errors"`
			TxDropped int `json:"tx_dropped"`
			Interfaces []struct {
				Name string `json:"name"`
				RxBytes int `json:"rx_bytes"`
				RxPackets int `json:"rx_packets"`
				RxErrors int `json:"rx_errors"`
				RxDropped int `json:"rx_dropped"`
				TxBytes int `json:"tx_bytes"`
				TxPackets int `json:"tx_packets"`
				TxErrors int `json:"tx_errors"`
				TxDropped int `json:"tx_dropped"`
			} `json:"interfaces"`
		} `json:"network"`
		Filesystem []struct {
			Device string `json:"device"`
			Capacity int64 `json:"capacity"`
			Usage int `json:"usage"`
			Available int `json:"available"`
			ReadsCompleted int `json:"reads_completed"`
			ReadsMerged int `json:"reads_merged"`
			SectorsRead int `json:"sectors_read"`
			ReadTime int `json:"read_time"`
			WritesCompleted int `json:"writes_completed"`
			WritesMerged int `json:"writes_merged"`
			SectorsWritten int `json:"sectors_written"`
			WriteTime int `json:"write_time"`
			IoInProgress int `json:"io_in_progress"`
			IoTime int `json:"io_time"`
			WeightedIoTime int `json:"weighted_io_time"`
		} `json:"filesystem"`
		TaskStats struct {
			NrSleeping int `json:"nr_sleeping"`
			NrRunning int `json:"nr_running"`
			NrStopped int `json:"nr_stopped"`
			NrUninterruptible int `json:"nr_uninterruptible"`
			NrIoWait int `json:"nr_io_wait"`
		} `json:"task_stats"`
	} `json:"stats"`
}
