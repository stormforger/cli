import "./modules/options.js"
import scenario from "./modules/scenario.js"

// NOTE: No 'var' or 'const' for defines!!!
// Inject defines here with '--define defines.env="prod"'
defines = {};
var config = {
  env: defines.env || "staging",
  target: defines.target || "https://testapp.loadtest.party",
}

definition.addTarget(config.target)
definition.session("hello", scenario(config.env))
