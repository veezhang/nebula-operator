package e2e

import (
	"github.com/vesoft-inc/nebula-operator/tests/e2e/e2ematcher"
	"sigs.k8s.io/e2e-framework/third_party/helm"

	"github.com/vesoft-inc/nebula-operator/tests/e2e/envfuncsext"
)

const (
	LabelCustomConfig        = "custom config"
	LabelCustomConfigStatic  = "static"
	LabelCustomConfigDynamic = "dynamic"
)

var testCasesCustomConfig []ncTestCase

func init() {
	testCasesCustomConfig = append(testCasesCustomConfig, testCasesCustomConfigStatic...)
	testCasesCustomConfig = append(testCasesCustomConfig, testCasesCustomConfigDynamic...)
}

// test cases about static custom config
var testCasesCustomConfigStatic = []ncTestCase{
	{
		Name: "custom config for static",
		Labels: map[string]string{
			LabelKeyCategory: LabelCustomConfig,
			LabelKeyGroup:    LabelCustomConfigStatic,
		},
		InstallWaitNCOptions: []envfuncsext.NebulaClusterOption{
			envfuncsext.WithNebulaClusterReadyFuncs(
				envfuncsext.NebulaClusterReadyFuncForFields(false, map[string]any{
					"Spec": map[string]any{
						"Graphd": map[string]any{
							"Replicas": e2ematcher.ValidatorEq(2),
						},
						"Metad": map[string]any{
							"Replicas": e2ematcher.ValidatorEq(3),
						},
						"Storaged": map[string]any{
							"Replicas": e2ematcher.ValidatorEq(3),
						},
					},
				}),
				envfuncsext.DefaultNebulaClusterReadyFunc,
			),
		},
		LoadLDBC: true,
		UpgradeCases: []ncTestUpgradeCase{
			{
				Name:        "update configs",
				UpgradeFunc: nil,
				UpgradeNCOptions: []envfuncsext.NebulaClusterOption{
					envfuncsext.WithNebulaClusterHelmRawOptions(
						helm.WithArgs(
							"--set-string", "nebula.graphd.config.max_sessions_per_ip_per_user=100",
							"--set-string", "nebula.metad.config.default_parts_num=30",
							"--set-string", "nebula.storaged.config.minimum_reserved_bytes=134217728",
						),
					),
				},
				UpgradeWaitNCOptions: []envfuncsext.NebulaClusterOption{
					envfuncsext.WithNebulaClusterReadyFuncs(
						envfuncsext.NebulaClusterReadyFuncForFields(false, map[string]any{
							"Spec": map[string]any{
								"Graphd": map[string]any{
									"Replicas": e2ematcher.ValidatorEq(2),
									"Config": map[string]any{
										"max_sessions_per_ip_per_user": e2ematcher.ValidatorEq("100"),
									},
								},
								"Metad": map[string]any{
									"Replicas": e2ematcher.ValidatorEq(3),
									"Config": map[string]any{
										"default_parts_num": e2ematcher.ValidatorEq("30"),
									},
								},
								"Storaged": map[string]any{
									"Replicas": e2ematcher.ValidatorEq(3),
									"Config": map[string]any{
										"minimum_reserved_bytes": e2ematcher.ValidatorEq("134217728"),
									},
								},
							},
						}),
						// TODO: check whether the change is really successful via host:port/flags api
						envfuncsext.DefaultNebulaClusterReadyFunc,
					),
				},
			},
		},
	},
}

// test cases about dynamic custom config
var testCasesCustomConfigDynamic = []ncTestCase{
	// TODO
}
