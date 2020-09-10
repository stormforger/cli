import "./modules/options.js"
import scenario from "./modules/scenario.js"

const config = {
  env: replacements.ENV,
}

definition.session("hello", scenario(config))
