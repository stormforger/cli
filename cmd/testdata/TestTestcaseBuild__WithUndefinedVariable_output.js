// testdata/TestTestcaseBuild__WithUndefinedVariable_main.mjs
defines = {};
var config = {
  target: defines.target || "http://testapp.loadtest.party"
};
definition.addTarget(config.target);
