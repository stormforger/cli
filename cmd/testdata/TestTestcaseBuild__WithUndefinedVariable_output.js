// testdata/TestTestcaseBuild__WithUndefinedVariable_main.mjs
defines = {};
var config = {
  target: defines.target || "http://testapp.loadtest.party"
};
definition.setArrivalPhases([{ duration: 60, rate: 0 }]);
definition.addTarget(config.target);
