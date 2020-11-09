import "./modules/options.js"
import scenario from "./modules/scenario.js"

// TODO do we want a "convention" here? An object
// that can be used to replace value into, would be a start…
const config = {
  env: ENV || "staging",
}

definition.session("hello", scenario(config))
