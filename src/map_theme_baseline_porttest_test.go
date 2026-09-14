//go:build porttest

package opennox

// Corrected C: 3,392 cases, repeated in default/server/highres before conversion.
var themeCExpected = map[string]string{
	"algorithms":           "c5d425cc75b6e5c85ee575bbf8e8148b3aa577c6367a897c9887e718537beb39",
	"choices":              "0b2883ce7c5747a72c8e2ebaf78a09661074c7cb78c8eeae5afa96d6b3092b58",
	"conditions":           "bf6756d1411b44a6069b9ef85ec494b62af5524892101b6da7a39dd8acc68e48",
	"decor-copies":         "6103f157a2692f0b415fc12ef23913f3d5fed9835bff60b5a903f44df75428c5",
	"decor-properties":     "06102815f81b3dc3894c417e6a750c2b080c8cc3d61681fb215589b0bb9aa1d8",
	"decor-sets":           "6ff8418c7122bdc200ce50ab0d5b57e2cd489e295b2ef434c29bb905d86fccb9",
	"decorations":          "453c2a2852572f5486630259c33933dae5cda3351d61a44edb1268b7f408127a",
	"equipment-boundaries": "5d50366d4f75efccb43c41d9baf63fac3815d31ed4d4c370faadd5810a1764ef",
	"equipment-sets":       "fa554b14c1f098317a2ee0596e2f33aa2b4397b8833613b9674e666ce24581aa",
	"exits-prefabs":        "353600a19665db12d56236107b88028d8929fde16fd0ff36f28b7dd74e3a6b31",
	"filtered-tokens":      "9e3acea2cfcacc07bdb958e45ce4c504f4313558d902db176378a9de16eb5a35",
	"foreach":              "cdbb3489f43dd761a26953b381dec8f089a3915c9d37f71f0b11758d6c73a4df",
	"frequency-validation": "83ed56440a122c1076d8baab0965803b2f0ba1761d6462184d920338275a7fd0",
	"full-files":           "7646c7f549b5f9ca275e8499912565321e13da2c83ac398661dff43593d1a36a",
	"operators":            "9a661b17463131316dd6095933f073286a6710ada7161d4f1d6efb36b8f72f74",
	"raw-tokens":           "6c27914fd71312bc40e28b194ddb44b036508005cbbef495f555d70f1a4ff3cd",
	"skips":                "f82c246c5ed7d37638310626eab2560e79def86dd1cf67d3873a3aae0eebb8bc",
	"spell-sets":           "e7bb8ddc2f9bbb0e8fee64a75dd59c78f732b82ef360de30bc0aba0e20a2f4ff",
	"template-removals":    "34b1ac9e2920987c8196c0f5ea808957bf0d1778da79912f74ac0789526fcf69",
}
