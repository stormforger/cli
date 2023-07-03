import "./modules/options.js"
import scenario from "./modules/scenario.js"

// NOTE: `--define` works on global identifiers only! To make `defines` a global identifier, it MUST NOT be defined via `var`/`let`/`const`.
// To replace fields of `defines`, use '--define defines.env="prod"'
defines = {};
var config = {
  env: defines.env || "staging",
  target: defines.target || "https://testapp.loadtest.party",
}

definition.addTarget(config.target)
definition.session("hello", scenario(config.env))
