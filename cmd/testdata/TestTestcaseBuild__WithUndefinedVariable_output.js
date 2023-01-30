// testdata/TestTestcaseBuild__WithUndefinedVariable_main.mjs
var config = {
  env: ENV || "staging"
  // Intentionally undefined in the go test
};
definition.addTarget(env);
