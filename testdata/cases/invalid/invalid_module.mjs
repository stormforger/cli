import "./doesnotexist/options.js"
import scenario from "./doesnotexist/scenario.js"

// TODO do we want a "convention" here? An object
// that can be used to replace value into, would be a start…
var replacements = {}

const config = {
  env: replacements.ENV,
}

definition.session("hello", scenario(config))
