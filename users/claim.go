package users

func GetClaimTypes() map[string]string {
	claimTypes := make(map[string]string)
	claimTypes["cpuCores"] = "CPU Cores"
	claimTypes["gpuCores"] = "GPU Cores"
	claimTypes["templateRead"] = "Template Read"
	claimTypes["templateWrite"] = "Template Write"
	claimTypes["templateDelete"] = "Template Delete"
	claimTypes["templateGroup"] = "Template Group"
	return claimTypes
}
