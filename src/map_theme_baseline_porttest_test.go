//go:build porttest

package opennox

// Corrected C captures repeated in default/server/highres.
var themeCExpected = map[string]string{
	"algorithms":           "c5d425cc75b6e5c85ee575bbf8e8148b3aa577c6367a897c9887e718537beb39",
	"conditions":           "bf6756d1411b44a6069b9ef85ec494b62af5524892101b6da7a39dd8acc68e48",
	"filtered-tokens":      "9e3acea2cfcacc07bdb958e45ce4c504f4313558d902db176378a9de16eb5a35",
	"frequency-validation": "83ed56440a122c1076d8baab0965803b2f0ba1761d6462184d920338275a7fd0",
	"operators":            "9a661b17463131316dd6095933f073286a6710ada7161d4f1d6efb36b8f72f74",
	"raw-tokens":           "6c27914fd71312bc40e28b194ddb44b036508005cbbef495f555d70f1a4ff3cd",
	"skips":                "f82c246c5ed7d37638310626eab2560e79def86dd1cf67d3873a3aae0eebb8bc",
	"template-removals":    "34b1ac9e2920987c8196c0f5ea808957bf0d1778da79912f74ac0789526fcf69",
}
