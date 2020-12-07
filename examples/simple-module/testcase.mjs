/*
A simple example showing how to import functions and configuration from a module.
*/

import {helloworld, config} from "./modules/helper.mjs";

definition.setTarget(config.url);

definition.setArrivalPhases([
  {duration: 60, rate: 10},
]);

definition.session("helloworld", helloworld(config));
