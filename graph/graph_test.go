package graph

func nodeInfosEqual(expected, actual NodeInfo) bool {
	return expected.Name == actual.Name && expected.Version == actual.Version && expected.Timestamp == actual.Timestamp
}

func createTestNodeInfo(pi PackageInfo, version string) NodeInfo {
	return NodeInfo{
		id:        -1,
		stringID:  "invalid",
		Name:      pi.Name,
		Version:   version,
		Timestamp: pi.Versions[version].Timestamp,
	}
}
